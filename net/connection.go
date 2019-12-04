package net

import (
	"errors"
	"fmt"
	"github.com/azd1997/zinx/iface"
	"github.com/azd1997/zinx/utils"
	"io"
	"net"
)

// Connection TCP连接，iface.IConnection接口的实现
type Connection struct {

	// 当前链接隶属于哪个server
	TcpServer iface.IServer

	// Conn 当前连接的TCP套接字
	Conn *net.TCPConn

	// ID 连接ID
	ID uint32

	// isClosed 当前连接状态
	isClosed bool

	// ExitChan 告知当前连接退出的channel
	ExitChan chan bool

	// MsgHandler 服务端注册的连接对应的消息管理模块（多路由）
	MsgHandler iface.IMsgHandler

	// msgChan 无缓冲通道，用于读写两个goroutine之间的消息通信
	msgChan chan []byte

	//有缓冲管道，用于读、写两个goroutine之间的消息通信
	msgBuffChan chan []byte
}

// NewConnection 新建TCP连接对象
func NewConnection(server iface.IServer, conn *net.TCPConn, id uint32, msgHandler iface.IMsgHandler) *Connection {
	c := &Connection{
		TcpServer:server,
		Conn:     conn,
		ID:       id,
		isClosed: false,
		MsgHandler:msgHandler,
		ExitChan: make(chan bool, 1), // 有缓冲通道
		msgChan:make(chan []byte),		// 读写goroutine间的消息通道
		msgBuffChan:make(chan []byte, utils.GlobalObject.MaxWorkerTaskLen),

	}

	//将新创建的Conn添加到链接管理中
	c.TcpServer.GetConnMgr().Add(c)   //将当前新创建的连接添加到ConnManager中

	return c
}

// Start 启动连接
func (c *Connection) Start() {
	fmt.Printf("Conn start... conn ID = %d\n", c.ID)

	// 启动当前连接的读数据业务
	go c.startReader()
	// 启动当前连接的写数据goroutine
	go c.startWriter()

	//按照用户传递进来的创建连接时需要处理的业务，执行钩子方法
	c.TcpServer.CallOnConnStart(c)
}

// Stop 停止连接，结束当前连接的工作
func (c *Connection) Stop() {
	fmt.Printf("Conn #%d stop\n", c.ID)

	// 如果当前连接已经关闭
	if c.isClosed {
		return
	}
	c.isClosed = true


	//如果用户注册了该链接的关闭回调业务，那么在此刻应该显示调用
	c.TcpServer.CallOnConnStop(c)

	// 关闭socket
	c.Conn.Close()
	c.ExitChan <- true

	// 将连接从链接管理器中删除
	c.TcpServer.GetConnMgr().Remove(c)

	// 回收资源
	close(c.ExitChan)
	close(c.msgBuffChan)
	close(c.msgChan)
}

// GetTCPConn 获取当前连接绑定的socket
func (c *Connection) GetTCPConn() *net.TCPConn {
	return c.Conn
}

// GetConnID 获取当前连接的ID
func (c *Connection) GetConnID() uint32 {
	return c.ID
}

// RemoteAddr 获取远程客户端的TCP状态 IP Port
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// SendMsg 发送数据
func (c *Connection) SendMsg(msgId uint32, data []byte) error {

	if c.isClosed {
		return errors.New("connection closed when send msg")
	}

	// 将data封包
	dp := NewDataPack()
	msg, err := dp.Pack(NewMessage(msgId, data))
	if err != nil {
		fmt.Println("Pack error msg id = ", msgId)
		return  errors.New("Pack error msg ")
	}

	c.msgChan <- msg

	return nil
}

// SendBuffMsg 发送数据
func (c *Connection) SendBuffMsg(msgId uint32, data []byte) error {

	if c.isClosed {
		return errors.New("connection closed when send msg")
	}

	// 将data封包
	dp := NewDataPack()
	msg, err := dp.Pack(NewMessage(msgId, data))
	if err != nil {
		fmt.Println("Pack error msg id = ", msgId)
		return  errors.New("Pack error msg ")
	}

	c.msgBuffChan <- msg

	return nil
}

// startReader 启动连接的读数据，持续从客户端读取数据，并交给handleAPI处理
func (c *Connection) startReader() {
	fmt.Println("Reader goroutine is running...")
	defer fmt.Printf("connID = %d; Reader exits, remote addr is %s\n", c.ID, c.RemoteAddr())
	defer c.Stop()

	for {

		// 创建拆包对象
		dp := NewDataPack()

		// 读取客户端的Msg Head
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConn(), headData); err != nil {
			fmt.Println("read msg head error ", err)
			break
		}

		//拆包，得到msgid 和 datalen 放在msg中
		msg , err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack error ", err)
			break
		}

		//根据 dataLen 读取 data，放在msg.Data中
		var data []byte
		if msg.GetDataLen() > 0 {
			data = make([]byte, msg.GetDataLen())
			if _, err := io.ReadFull(c.GetTCPConn(), data); err != nil {
				fmt.Println("read msg data error ", err)
				break
			}
		}
		msg.SetData(data)

		req := NewRequest(c, msg)


		// 从路由中找到注册绑定的Conn对应的router调用
		// 执行注册的路由方法
		// go c.MsgHandler.DoMsgHandler(req)

		// 如果工作池还有空闲worker则传给worker，否则从绑定好的消息和对应处理方法执行相应的handle方法
		if utils.GlobalObject.WorkerPoolSize > 0 {
			//已经启动工作池机制，将消息交给Worker处理
			c.MsgHandler.SendMsgToTaskQueue(req)
		} else {
			// 没有配置工作池机制或者worker已耗尽时则仍然按照原来的而做法单开goroutine去处理
			//从绑定好的消息和对应的处理方法中执行对应的Handle方法
			go c.MsgHandler.DoMsgHandler(req)
		}
	}
}

// startWriter 写消息goroutine，用户将数据发送给客户端
func (c *Connection) startWriter() {
	defer fmt.Println(c.RemoteAddr().String(), " conn writer exit.")

	// 除非出错或者收到退出信号，否则一直循环
	for {
		select {
		// 无缓冲数据通道的数据接收处理
		case data := <-c.msgChan:
			//有数据要写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send Data error:, ", err, " Conn Writer exit")
				return
			}
		// 有缓冲数据通道的数据接收处理
		case data, ok := <- c.msgBuffChan:
			if ok {
				//有数据要写给客户端
				if _, err := c.Conn.Write(data); err != nil {
					fmt.Println("Send Buff Data error:, ", err, " Conn Writer exit")
					return
				}
			} else {
				fmt.Println("msgBuffChan is closed")
				break
			}
		case <- c.ExitChan:
			//conn已经关闭
			return
		}
	}
}
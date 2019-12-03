package net

import (
	"errors"
	"fmt"
	"github.com/azd1997/zinx/iface"
	"io"
	"net"
)

// Connection TCP连接，iface.IConnection接口的实现
type Connection struct {

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
}

// NewConnection 新建TCP连接对象
func NewConnection(conn *net.TCPConn, id uint32, msgHandler iface.IMsgHandler) *Connection {
	return &Connection{
		Conn:     conn,
		ID:       id,
		isClosed: false,
		MsgHandler:msgHandler,
		ExitChan: make(chan bool, 1), // 有缓冲通道
		msgChan:make(chan []byte),		// 读写goroutine间的消息通道
	}
}

// Start 启动连接
func (c *Connection) Start() {
	fmt.Printf("Conn start... conn ID = %d\n", c.ID)

	// 启动当前连接的读数据业务
	go c.startReader()
	// 启动当前连接的写数据goroutine
	go c.startWriter()

	// 等待退出信息，来停止连接
	//for {
	//	select {
	//	case <-c.ExitChan:
	//		// 收到退出信息，不再阻塞
	//		return
	//	}
	//}
}

// Stop 停止连接，结束当前连接的工作
func (c *Connection) Stop() {
	fmt.Printf("Conn #%d stop\n", c.ID)

	// 如果当前连接已经关闭
	if c.isClosed {
		return
	}
	c.isClosed = true

	// 关闭socket
	c.Conn.Close()

	// 回收资源
	close(c.ExitChan)
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

// Send 发送数据
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

	// 写回客户端
	//cnt, err := c.Conn.Write(msg)
	//if cnt != len(data) || err != nil {
	//	fmt.Println("Write msg id ", msgId, " error ")
	//	c.ExitChan <- true
	//	return errors.New("conn Write error")
	//}
	c.msgChan <- msg

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
		go c.MsgHandler.DoMsgHandler(req)
	}
}

// startWriter 写消息goroutine，用户将数据发送给客户端
func (c *Connection) startWriter() {
	defer fmt.Println(c.RemoteAddr().String(), " conn writer exit.")

	// 除非出错或者收到退出信号，否则一直循环
	for {
		select {
		case data := <-c.msgChan:
			//有数据要写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send Data error:, ", err, " Conn Writer exit")
				return
			}
		case <- c.ExitChan:
			//conn已经关闭
			return
		}
	}
}
package net

import (
	"fmt"
	"github.com/azd1997/zinx/iface"
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

	// handleAPI 当前连接绑定的处理业务方法API
	//handleAPI iface.HandleFunc

	// ExitChan 告知当前连接退出的channel
	ExitChan chan bool

	// Router 该连接处理的方法路由
	Router iface.IRouter
}

// NewConnection 新建TCP连接对象
//func NewConnection(conn *net.TCPConn, id uint32, callbackAPI iface.HandleFunc) *Connection {
//	return &Connection{
//		Conn:      conn,
//		ID:        id,
//		isClosed:  false,
//		handleAPI: callbackAPI,
//		ExitChan:  make(chan bool, 1),	// 有缓冲通道
//	}
//}

// NewConnection 新建TCP连接对象
func NewConnection(conn *net.TCPConn, id uint32, router iface.IRouter) *Connection {
	return &Connection{
		Conn:      conn,
		ID:        id,
		isClosed:  false,
		Router: router,
		ExitChan:  make(chan bool, 1),	// 有缓冲通道
	}
}

// Start 启动连接
func (c *Connection) Start() {
	fmt.Printf("Conn start... conn ID = %d\n", c.ID)

	// 启动当前连接的读数据业务
	go c.startReader()

	// 等待退出信息，来停止连接
	for {
		select {
		case <-c.ExitChan:
			// 收到退出信息，不再阻塞
			return
		}
	}
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
func (c *Connection) Send(data []byte) error {

	// TODO: 待实现，这里随便写的

	cnt, err := c.Conn.Write(data)
	if cnt != len(data) || err != nil {
		return err
	}

	return nil
}

// startReader 启动连接的读数据，持续从客户端读取数据，并交给handleAPI处理
func (c *Connection) startReader() {
	fmt.Println("Reader goroutine is running...")
	defer fmt.Printf("connID = %d; Reader exits, remote addr is %s\n", c.ID, c.RemoteAddr())
	defer c.Stop()

	for {
		// 读客户端数据到buffer区
		buf := make([]byte, 512)
		//cnt, err := c.Conn.Read(buf)
		_, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Printf("Receive buf error: %s\n", err)
			continue
		}

		// 调用当前连接所绑定的handleAPI去处理客户端传来的数据
		//if err = c.handleAPI(c.Conn, buf, cnt); err != nil {
		//	fmt.Printf("Conn #%d: handleAPI error: %s\n", c.ID, err)
		//	break
		//}

		// 得到当前conn数据的Request
		//req := NewRequest(c, buf[:cnt])
		// 按cnt来截取缓冲区会将字符串末尾添加的\n删除，这样最后客户端打印打印出的before ping; ping; after ping就黏在一块了
		req := NewRequest(c, buf)

		// 从路由中找到注册绑定的Conn对应的router调用
		// 执行注册的路由方法
		go func(req iface.IRequest) {
			c.Router.PreHandle(req)
			c.Router.Handle(req)
			c.Router.PostHandle(req)
		}(req)
	}
}
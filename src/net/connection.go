package net

import (
	"fmt"
	"github.com/azd1997/zinx/src/iface"
	"net"
)

// Connection TCP连接，iface.IConnection接口的实现
type Connection struct {

	// Conn 当前连接的TCP套接字
	Conn *net.TCPConn

	// ID 连接ID
	ID uint32

	// 当前连接状态
	isClosed bool

	// 当前连接绑定的处理业务方法API
	handleAPI iface.HandleFunc

	// 告知当前连接退出的channel
	ExitChan chan bool
}

// NewConnection 新建TCP连接对象
func NewConnection(conn *net.TCPConn, id uint32, callbackAPI iface.HandleFunc) *Connection {
	return &Connection{
		Conn:      conn,
		ID:        id,
		isClosed:  false,
		handleAPI: callbackAPI,
		ExitChan:  make(chan bool, 1),	// 有缓冲通道
	}
}

// Start 启动连接
func (c *Connection) Start() {
	fmt.Printf("Conn start... conn ID = %d\n", c.ID)

	// 启动当前连接的读数据业务
	go c.startReader()
}

// Stop 停止连接
func (c *Connection) Stop() {

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
		cnt, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Printf("Receive buf error: %s\n", err)
			continue
		}

		// 调用当前连接所绑定的handleAPI去处理客户端传来的数据
		if err = c.handleAPI(c.Conn, buf, cnt); err != nil {
			fmt.Printf("Conn #%d: handleAPI error: %s\n", c.ID, err)
			break
		}
	}
}
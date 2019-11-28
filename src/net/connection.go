package net

import (
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

}

// Stop 停止连接
func (c *Connection) Stop() {

}

// GetTCPConn 获取当前连接绑定的socket
func (c *Connection) GetTCPConn() *net.TCPConn {

}

// GetConnID 获取当前连接的ID
func (c *Connection) GetConnID() uint32 {

}

// RemoteAddr 获取远程客户端的TCP状态 IP Port
func (c *Connection) RemoteAddr() net.Addr {

}

// Send 发送数据
func (c *Connection) Send(data []byte) error {

}
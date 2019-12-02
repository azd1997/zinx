package iface

import "net"

// IConnection 连接的抽象接口
type IConnection interface {
	// Start 启动连接
	Start()

	// Stop 停止连接
	Stop()

	// GetTCPConn 获取当前连接绑定的socket
	GetTCPConn() *net.TCPConn

	// GetConnID 获取当前连接的ID
	GetConnID() uint32

	// RemoteAddr 获取远程客户端的TCP状态 IP Port
	RemoteAddr() net.Addr

	// Send 发送数据
	Send(data []byte) error
}

// HandleFunc 定义连接绑定的 处理业务的函数类型。 data处理业务的数据， l为处理数据的长度
type HandleFunc func(conn *net.TCPConn, data []byte, l int) error
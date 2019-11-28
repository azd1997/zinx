package net

import "github.com/azd1997/zinx/src/iface"

type Request struct {

	// 已经和客户端建立好的连接
	conn iface.IConnection

	// 客户端的请求数据
	data []byte
}

// GetConn 获取当前连接
func (r *Request) GetConn() iface.IConnection {
	return r.conn
}

// GetData 获取当前连接的请求数据
func (r *Request) GetData() []byte {
	return r.data
}

//
func NewRequest(conn iface.IConnection, data []byte) iface.IRequest {
	return &Request{
		conn: conn,
		data: data,
	}
}

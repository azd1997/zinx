package net

import "github.com/azd1997/zinx/iface"

type Request struct {

	// 已经和客户端建立好的连接
	conn iface.IConnection

	// 客户端的请求数据
	//data []byte
	msg iface.IMessage
}

// GetConn 获取当前连接
func (r *Request) GetConn() iface.IConnection {
	return r.conn
}

// GetData 获取当前连接的请求数据
func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

// GetMsgID 获取请求的消息的ID
func (r *Request) GetMsgId() uint32 {
	return r.msg.GetMsgId()
}

// NewRequest 新建请求
func NewRequest(conn iface.IConnection, msg iface.IMessage) iface.IRequest {
	return &Request{
		conn: conn,
		msg: msg,
	}
}

package iface

// IRequest 将客户端请求的连接信息和请求的数据 包装到一个Request中
type IRequest interface {

	// GetConn 获取当前连接
	GetConn() IConnection

	// GetData 获取当前连接的请求数据
	GetData() []byte

	// GetMsgId 获取消息Id
	GetMsgId() uint32
}

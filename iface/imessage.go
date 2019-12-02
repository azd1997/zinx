package iface

// IMessage 将请求的一个消息封装到message中，定义抽象层接口
type IMessage interface {

	// GetMsgId 获取消息ID
	GetMsgId() uint32
	// GetDataLen 获取消息数据段长度
	GetDataLen() uint32
	// GetData 获取消息内容
	GetData() []byte

	// SetMsgId 设置消息ID
	SetMsgId(uint32)
	// 设置消息数据段长度
	SetDataLen(uint32)
	// 设置消息内容
	SetData([]byte)
}
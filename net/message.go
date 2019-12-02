package net

// Message 实现IMessage接口，对数据的一个封装
type Message struct {
	// Id 消息Id
	Id      uint32
	// DataLen 消息数据段长度
	DataLen uint32
	// Data 消息的数据段
	Data    []byte
}


// NewMessage 创建一个Message消息包
func NewMessage(id uint32, data []byte) *Message {
	return &Message{
		Id:     id,
		DataLen: uint32(len(data)),
		Data:   data,
	}
}

// GetDataLen 获取消息数据段长度
func (msg *Message) GetDataLen() uint32 {
	return msg.DataLen
}

// GetMsgId 获取消息ID
func (msg *Message) GetMsgId() uint32 {
	return msg.Id
}

// GetData 获取消息内容
func (msg *Message) GetData() []byte {
	return msg.Data
}

// SetDataLen 设置消息数据段长度
func (msg *Message) SetDataLen(len uint32) {
	msg.DataLen = len
}

// SetMsgId 设置消息ID
func (msg *Message) SetMsgId(msgId uint32) {
	msg.Id = msgId
}

// SetData 设置消息内容
func (msg *Message) SetData(data []byte) {
	msg.Data = data
}

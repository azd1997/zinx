package net

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/azd1997/zinx/iface"
	"github.com/azd1997/zinx/utils"
)

// DataPack 封包拆包类实例，暂时不需要成员
// TODO: 原作者此处使用了一个这样的面向对象方式，但Pack与Unpack其实完全可以在Message实现，所以我后面要修改此处，暂时跟着作者一样
type DataPack struct {}

// NewDataPack 封包拆包实例初始化方法
func NewDataPack() *DataPack {
	return &DataPack{}
}

// GetHeadLen 获取包头长度方法
func(dp *DataPack) GetHeadLen() uint32 {
	//Id uint32(4字节) +  DataLen uint32(4字节)
	return 8
}

// Pack 封包方法(压缩数据)
func(dp *DataPack) Pack(msg iface.IMessage) ([]byte, error) {
	//创建一个存放bytes字节的缓冲
	dataBuff := bytes.NewBuffer([]byte{})

	//写dataLen
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetDataLen()); err != nil {
		return nil, err
	}

	//写msgID
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}

	//写data数据
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil ,err
	}

	return dataBuff.Bytes(), nil
}

// Unpack 拆包方法(解压数据)
func(dp *DataPack) Unpack(binaryData []byte) (iface.IMessage, error) {
	//创建一个从输入二进制数据的ioReader
	dataBuff := bytes.NewReader(binaryData)

	//只解压head的信息，得到dataLen和msgID
	msg := &Message{}

	//读dataLen
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}

	//读msgID
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}

	//判断dataLen的长度是否超出我们允许的最大包长度
	if utils.GlobalObject.MaxPacketSize > 0 && msg.DataLen > utils.GlobalObject.MaxPacketSize {
		return nil, errors.New("too large msg data recieved")
	}

	//这里只需要把head的数据拆包出来就可以了，然后再通过head的长度，再从conn读取一次数据
	return msg, nil
}


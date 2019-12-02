package net

import (
	"fmt"
	"github.com/azd1997/zinx/iface"
	"strconv"
)

// MsgHandle 消息管理（多路由）
type MsgHandle struct{
	// Apis 存放每个MsgId 所对应的处理方法的map属性
	Apis map[uint32] iface.IRouter
}

// NewMsgHandle 新建消息管理
func NewMsgHandle() *MsgHandle {
	return &MsgHandle {
		Apis:make(map[uint32]iface.IRouter),
	}
}

// DoMsgHandler 马上以非阻塞方式处理消息
func (mh *MsgHandle) DoMsgHandler(request iface.IRequest)  {
	handler, ok := mh.Apis[request.GetMsgId()]
	if !ok {
		fmt.Println("api msgId = ", request.GetMsgId(), " is not FOUND!")
		return
	}

	//执行对应处理方法
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

// AddRouter 为消息添加具体的处理逻辑
func (mh *MsgHandle) AddRouter(msgId uint32, router iface.IRouter) {
	//1 判断当前msg绑定的API处理方法是否已经存在
	if _, ok := mh.Apis[msgId]; ok {
		panic("repeated api , msgId = " + strconv.Itoa(int(msgId)))
	}
	//2 添加msg与api的绑定关系
	mh.Apis[msgId] = router
	fmt.Println("Add api msgId = ", msgId)
}

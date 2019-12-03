package net

import (
	"fmt"
	"github.com/azd1997/zinx/iface"
	"github.com/azd1997/zinx/utils"
	"strconv"
)

// MsgHandle 消息管理（多路由）
type MsgHandle struct{
	// Apis 存放每个MsgId 所对应的处理方法的map属性
	Apis map[uint32] iface.IRouter

	// WorkerPoolSize 业务工作Worker池中worker的数量
	// 这里没有单独实现一个worker pool（实现了也得在这里使用）而是直接在MsgHandler中实现
	WorkerPoolSize uint32

	// TaskQueue Worker负责取任务的消息队列。一个worker对应一个chan
	TaskQueue      []chan iface.IRequest
}

// NewMsgHandle 新建消息管理
func NewMsgHandle() *MsgHandle {
	return &MsgHandle {
		Apis:make(map[uint32]iface.IRouter),
		WorkerPoolSize:utils.GlobalObject.WorkerPoolSize,
		TaskQueue:make([]chan iface.IRequest, utils.GlobalObject.WorkerPoolSize),
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

//启动一个Worker工作流程
func (mh *MsgHandle) StartOneWorker(workerID int, taskQueue chan iface.IRequest) {
	fmt.Println("Worker ID = ", workerID, " is started.")
	//不断的等待队列中的消息
	for {
		select {
		//有消息则取出队列的Request，并执行绑定的业务方法
		case request := <-taskQueue:
			mh.DoMsgHandler(request)
		}
	}
}

//启动worker工作池
func (mh *MsgHandle) StartWorkerPool() {
	//遍历需要启动worker的数量，依此启动
	for i:= 0; i < int(mh.WorkerPoolSize); i++ {
		//一个worker被启动
		//给当前worker对应的任务队列开辟空间
		mh.TaskQueue[i] = make(chan iface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		//启动当前Worker，阻塞的等待对应的任务队列是否有消息传递进来
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}

//将消息交给TaskQueue,由worker进行处理
func (mh *MsgHandle)SendMsgToTaskQueue(request iface.IRequest) {
	//根据ConnID来分配当前的连接应该由哪个worker负责处理
	//轮询的平均分配法则

	//得到需要处理此条连接的workerID
	workerID := request.GetConn().GetConnID() % mh.WorkerPoolSize
	fmt.Println("Add ConnID=", request.GetConn().GetConnID()," request msgID=", request.GetMsgId(), "to workerID=", workerID)
	//将请求消息发送给任务队列
	mh.TaskQueue[workerID] <- request
}
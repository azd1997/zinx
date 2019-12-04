package iface

type IServer interface {

	// Start 启动服务器
	Start()

	// Stop 停止服务器
	Stop()

	// Serve 运行服务器
	Serve()

	// AddRouter 路由功能，给当前的服务器注册一个路由方法，供客户端的连接处理使用
	AddRouter(msgId uint32, router IRouter)

	//得到链接管理
	GetConnMgr() IConnManager

	//设置该Server的连接创建时Hook函数
	SetOnConnStart(func(IConnection))
	//设置该Server的连接断开时的Hook函数
	SetOnConnStop(func(IConnection))
	//调用连接OnConnStart Hook函数
	CallOnConnStart(conn IConnection)
	//调用连接OnConnStop Hook函数
	CallOnConnStop(conn IConnection)
}

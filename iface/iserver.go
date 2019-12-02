package iface

type IServer interface {

	// Start 启动服务器
	Start()

	// Stop 停止服务器
	Stop()

	// Serve 运行服务器
	Serve()

	// AddRouter 路由功能，给当前的服务器注册一个路由方法，供客户端的连接处理使用
	AddRouter(router IRouter)
}
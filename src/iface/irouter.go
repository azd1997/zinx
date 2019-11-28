package iface

// 路由的抽象接口，路由里的数据都是IRequest
type IRouter interface {

	// PreHandle 处理conn业务之前的钩子方法(Hook)
	PreHandle(r IRequest)

	// Handle 处理conn业务的主方法(Hook)
	Handle(r IRequest)

	// PostHandle 处理conn业务之后的钩子方法(Hook)
	PostHandle(r IRequest)
}

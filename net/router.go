package net

import "github.com/azd1997/zinx/iface"

// 实现用户自己的IRouter，先嵌入这个BaseRouter，再根据需要对“基类”方法进行“重写”
type BaseRouter struct {}

// 这里的方法都为空的原因是，用户不一定都需要三个阶段的方法，所以用户自定义路由“继承”自BaseRouter可以不用再重写不需要的PreHandle或PostHandle方法

// PreHandle 处理conn业务之前的钩子方法(Hook)
func (br *BaseRouter) PreHandle(r iface.IRequest) {}

// Handle 处理conn业务的主方法(Hook)
func (br *BaseRouter) Handle(r iface.IRequest) {}

// PostHandle 处理conn业务之后的钩子方法(Hook)
func (br *BaseRouter) PostHandle(r iface.IRequest) {}
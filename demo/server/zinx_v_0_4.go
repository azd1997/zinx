// 基于Zinx框架开发的服务端应用程序
package main

import (
	"fmt"
	"github.com/azd1997/zinx/iface"
	"github.com/azd1997/zinx/net"
)

func main() {
	// 1. 创建Server
	// 输入配置文件时，如果是编译成了可执行文件，则应填"./conf/zinx.json"，
	// 但由于测试时是直接 go run zinx_v_0_x.go所以根目录是项目根目录，配置文件路径应填"./demo/server/conf/zinx.json"
	s := net.NewServer("./demo/server/conf/zinx.json")

	// 2. 添加自定义路由
	s.AddRouter(&PingRouter{})

	// 3. 启动Server
	s.Serve()
}

// 自定义路由
type PingRouter struct {
	net.BaseRouter // 先"继承"BaseRouter
}

// 覆盖（屏蔽内层方法）
func (r *PingRouter) PreHandle(req iface.IRequest) {
	fmt.Println("Call Router PreHandle")
	_, err := req.GetConn().GetTCPConn().Write([]byte("before ping\n"))
	if err != nil {
		fmt.Printf("prehandle error: %s\n", err)
	}
}

func (r *PingRouter) Handle(req iface.IRequest) {
	fmt.Println("Call Router Handle")
	_, err := req.GetConn().GetTCPConn().Write([]byte("ping ping ping\n"))
	if err != nil {
		fmt.Printf("handle error: %s\n", err)
	}
}

func (r *PingRouter) PostHandle(req iface.IRequest) {
	fmt.Println("Call Router PostHandle")
	_, err := req.GetConn().GetTCPConn().Write([]byte("after ping\n"))
	if err != nil {
		fmt.Printf("posthandle error: %s\n", err)
	}
}

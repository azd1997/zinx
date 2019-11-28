// 基于Zinx框架开发的服务端应用程序
package main

import "github.com/azd1997/zinx/src/net"

func main() {
	// 1. 创建Server
	s := net.NewServer("zinx_v_0_2")

	// 2. 启动Server
	s.Serve()
}

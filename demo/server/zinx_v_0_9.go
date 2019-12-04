// 基于Zinx框架开发的服务端应用程序
package main

import (
	"fmt"
	"github.com/azd1997/zinx/iface"
	"github.com/azd1997/zinx/net"
)


// 要测试消息管理模块，这里自定义了两个路由，
// 那么就开两个client，一个发Id为0的消息，一个Id为1
func main() {
	// 1. 创建Server
	// 输入配置文件时，如果是编译成了可执行文件，则应填"./conf/zinx.json"，
	// 但由于测试时是直接 go run zinx_v_0_x.go所以根目录是项目根目录，配置文件路径应填"./demo/server/conf/zinx.json"
	s := net.NewServer("./demo/server/conf/zinx.json")


	//注册链接hook回调函数
	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionLost)

	// 2. 添加自定义路由
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloRouter{})

	// 3. 启动Server
	s.Serve()
}

// 下面自定义了两个路由


type PingRouter struct {
	net.BaseRouter // 先"继承"BaseRouter
}

// 覆盖（屏蔽内层方法）
func (r *PingRouter) Handle(req iface.IRequest) {
	fmt.Println("Call Router Handle")

	// 读客户端消息
	fmt.Println("recv from client : msgId=", req.GetMsgId(), ", data=", string(req.GetData()))


	// 回写客户端消息
	err := req.GetConn().SendMsg(1, []byte("ping ping ping\n"))
	if err != nil {
		fmt.Printf("handle error: %s\n", err)
	}
}


type HelloRouter struct {
	net.BaseRouter // 先"继承"BaseRouter
}

// 覆盖（屏蔽内层方法）
func (r *HelloRouter) Handle(req iface.IRequest) {
	fmt.Println("Call Router Handle")

	// 读客户端消息
	fmt.Println("recv from client : msgId=", req.GetMsgId(), ", data=", string(req.GetData()))


	// 回写客户端消息
	err := req.GetConn().SendMsg(1, []byte("Hello Zinx Router v0.6\n"))
	if err != nil {
		fmt.Printf("handle error: %s\n", err)
	}
}


//创建连接的时候执行
func DoConnectionBegin(conn iface.IConnection) {
	fmt.Println("DoConnecionBegin is Called ... ")
	err := conn.SendMsg(2, []byte("DoConnection BEGIN..."))
	if err != nil {
		fmt.Println(err)
	}
}

//连接断开的时候执行
func DoConnectionLost(conn iface.IConnection) {
	fmt.Println("DoConneciotnLost is Called ... ")
}
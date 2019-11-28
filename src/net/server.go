package net

import (
	"fmt"
	"github.com/azd1997/zinx/src/iface"
	"net"
)

// Server iface.IServer接口的一个实现
type Server struct {

	// Name 服务器名称
	Name string

	// IPVersion 服务器绑定的IP版本
	IPVersion string

	// IP 服务器监听的IP地址
	IP string

	// Port 服务器监听端口
	Port int
}

// Start 启动
func (s *Server) Start() {
	fmt.Printf("[Start] Server starts listening at IP: %s, Port: %d\n", s.IP, s.Port)

	// 异步启动，避免阻塞主线程
	go func() {
		// 1. 获取一个TCP Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Printf("Start: %s\n", err)
			return
		}

		// 2. 监听服务器地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Printf("Start: %s\n", err)
			return
		}

		fmt.Printf("[Start] Server starts listening at IP: %s, Port: %d successful\n", s.IP, s.Port)

		var connID uint32 = 0

		// 3. 阻塞等待客户端连接，处理客户端连接业务
		for {
			// 如果有客户端连接，则阻塞会返回，往下执行
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Printf("Start: Listener AcceptTCP: %s\n", err)
				continue
			}

			// 已经与客户端进行连接，做一个基本的业务： 最大512字节长度的回显业务
			//go func() {
			//	for {
			//		buf := make([]byte, 512)
			//		cnt, err := conn.Read(buf)
			//		if err != nil {
			//			fmt.Printf("Start: Conn Read%s\n", err)
			//			continue
			//		}
			//		fmt.Printf("Receive from client: %s, cnt = %d\n", string(buf), cnt)
			//
			//		// 回显
			//		if _, err := conn.Write(buf[:cnt]); err != nil {
			//			fmt.Printf("Start: Conn Write back: %s\n", err)
			//			continue
			//		}
			//	}
			//}()

			// 原本回显的部分可以由connection实现，如下

			// 将处理当前连接的业务方法和conn进行绑定，得到我们的连接模块
			dealConn := NewConnection(conn, connID, CallBackToClient)
			connID++
			// 尝试启动连接模块
			go dealConn.Start()
		}
	}()
}

// Stop 停止
func (s *Server) Stop() {
	// TODO: 将服务器的资源、状态或者已经建立的连接等等进行停止或回收
}

// Serve 运行
func (s *Server) Serve() {
	// 启动server服务功能
	s.Start()

	// TODO: 可以做一些服务器启动之后的额外业务

	// 阻塞状态
	select {}
}

// NewServer 新建一个Server
func NewServer(name string) iface.IServer {
	return &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",	// TODO: 暂时写死
		Port:      8000,
	}
}

// CallBackToClient 回显函数，HandleFunc的一个实现（定义当前客户端连接所绑定的handleAPI），用于服务器的功能测试
// TODO: 以后要将这个handleAPI在demo里自定义
func CallBackToClient(conn *net.TCPConn, data []byte, cnt int) error {
	// 回显业务
	fmt.Println("[Conn Handle] CallBackToClient...")

	if _, err := conn.Write(data[:cnt]); err != nil {
		fmt.Printf("CallBackToClient: %s\n", err)
		return err
	} else {
		fmt.Printf("Received from client: %s, cnt = %d\n", string(data), cnt)
	}

	return nil
}
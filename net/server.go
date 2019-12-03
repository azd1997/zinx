package net

import (
	"fmt"
	"github.com/azd1997/zinx/iface"
	"github.com/azd1997/zinx/utils"
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

	// MsgHandler 服务端注册的连接对应的消息管理模块（多路由）
	MsgHandler iface.IMsgHandler
}

// Start 启动
func (s *Server) Start() {
	fmt.Printf("[Start] Server starts listening at IP: %s, Port: %d\n", s.IP, s.Port)
	fmt.Printf("[Zinx] Version: %s, MaxConn: %d, MaxPacketSize: %d\n",
		utils.GlobalObject.Version,
		utils.GlobalObject.MaxConn,
		utils.GlobalObject.MaxPacketSize)

	// 异步启动，避免阻塞主线程
	go func() {

		// 0. 启动 worker 工作池机制
		s.MsgHandler.StartWorkerPool()

		// 1. 获取一个TCP Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Printf("Start: resolving tcp addr err: %s\n", err)
			return
		}

		// 2. 监听服务器地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Printf("Start listen %s err: %s\n", s.IPVersion, err)
			return
		}

		fmt.Printf("[Start] Server starts listening at IP: %s, Port: %d successful\n", s.IP, s.Port)

		// TODO: 添加一个自动生成ID的方法
		var connID uint32 = 0

		// 3. 阻塞等待客户端连接，处理客户端连接业务
		for {
			// 如果有客户端连接，则阻塞会返回，往下执行
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Printf("Start: Listener AcceptTCP: %s\n", err)
				continue
			}

			// TODO: Server.Start()设置服务器最大连接数，炒锅最大连接数，则关闭该新连接

			// 将处理当前连接的业务方法和conn进行绑定，得到我们的连接模块
			dealConn := NewConnection(conn, connID, s.MsgHandler)
			connID++
			// 尝试启动连接模块
			go dealConn.Start()
		}
	}()
}

// Stop 停止
func (s *Server) Stop() {
	fmt.Printf("[STOP] Zinx server %s\n", s.Name)

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

// AddRouter 添加路由
func (s *Server) AddRouter(msgId uint32, router iface.IRouter) {
	s.MsgHandler.AddRouter(msgId, router)
	fmt.Println("Add Router Successfully!")
}

// NewServer 新建一个Server
func NewServer(configFile string) iface.IServer {

	// 先初始化全局配置文件
	// 尽管utils.init中已经加载过一次，但这里再加载可以保证每次启动服务器都能得到最新的配置
	utils.GlobalObject.Reload(configFile)

	return &Server{
		Name:      utils.GlobalObject.Name,
		IPVersion: "tcp4",
		IP:        utils.GlobalObject.Host,
		Port:      utils.GlobalObject.TcpPort,
		MsgHandler: NewMsgHandle(),
	}
}

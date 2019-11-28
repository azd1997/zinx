package net

// Server iface.IServer接口的一个实现
type Server struct {
	Name string
	IPVersion string
	IP string
	Port int
}

// Start 启动
func (s *Server) Start() {

}

// Stop 停止
func (s *Server) Stop() {

}

// Serve 运行
func (s *Server) Serve() {

}
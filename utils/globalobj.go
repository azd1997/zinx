package utils

import (
	"encoding/json"
	"github.com/azd1997/zinx/iface"
	"io/ioutil"
)


// GlobalObject 全局对象
var GlobalObject *GlobalObj

// GlobalObj 存储一切有关zinx框架的全局参数，供其他模块使用
type GlobalObj struct {

	// TcpServer 当前Zinx的全局Server对象
	TcpServer iface.IServer

	// Host 当前服务器主机IP
	Host string

	// TcpPort 当前服务器主机监听端口号
	TcpPort int

	// Name 当前服务器名称
	Name string

	// Version 当前Zinx版本号
	Version string

	// MaxPacketSize 数据包的最大值
	MaxPacketSize uint32

	// MaxConn 当前服务器主机允许的最大链接个数
	MaxConn int
}

// Reload 读取配置文件
func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("conf/zinx.json")
	if err != nil {
		panic(err)
	}
	//将json数据解析到struct中
	//fmt.Printf("json :%s\n", data)
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

// init 初始化配置
func init() {
	//初始化GlobalObject变量，设置一些默认值
	GlobalObject = &GlobalObj{
		Name:    "ZinxServerApp",
		Version: "V0.4",
		TcpPort: 8000,
		Host:    "0.0.0.0",
		MaxConn: 12000,
		MaxPacketSize:4096,
	}

	//从配置文件中加载一些用户配置的参数
	GlobalObject.Reload()
}
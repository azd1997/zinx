package utils

import (
	"encoding/json"
	"fmt"
	"github.com/azd1997/zinx/iface"
	"io/ioutil"
)


// GlobalObject 全局对象
var GlobalObject *GlobalObj

// GlobalObj 存储一切有关zinx框架的全局参数，供其他模块使用
type GlobalObj struct {

	/*server配置*/

	// TcpServer 当前Zinx的全局Server对象
	TcpServer iface.IServer

	// Host 当前服务器主机IP
	Host string

	// TcpPort 当前服务器主机监听端口号
	TcpPort int

	// Name 当前服务器名称
	Name string

	/*zinx配置*/

	// Version 当前Zinx版本号
	Version string

	// MaxPacketSize 数据包的最大值
	MaxPacketSize uint32

	// MaxConn 当前服务器主机允许的最大链接个数
	MaxConn int

	// WorkerPoolSize 工作池大小
	WorkerPoolSize uint32

	// MaxWorkerTaskLen 业务工作Worker对应负责的任务队列最大任务存储数量
	MaxWorkerTaskLen uint32

	/*配置文件*/
	// ConfFilePath 配置文件路径
	ConfFilePath string
}

// Reload 读取配置文件。
// 读取成功就将配置文件属性赋予； 如失败则退出等下一次被调用时指定可用的配置路径
func (g *GlobalObj) Reload(confFile string) {
	data, err := ioutil.ReadFile(confFile)
	if err != nil {
		fmt.Println("global config read: ", err)
		return
	}
	//将json数据解析到struct中
	//fmt.Printf("json :%s\n", data)
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		fmt.Println("global config read: ", err)
		return
	}

	g.ConfFilePath = confFile
}

// init 初始化配置
func init() {
	//初始化GlobalObject变量，设置一些默认值
	GlobalObject = &GlobalObj{
		Name:    "ZinxServerApp",
		Version: "V0.7",
		TcpPort: 8000,
		Host:    "0.0.0.0",
		MaxConn: 12000,
		MaxPacketSize:4096,
		WorkerPoolSize:10,
		MaxWorkerTaskLen:1024,

		ConfFilePath: "conf/zinx.json",
	}

	//从配置文件中加载一些用户配置的参数
	GlobalObject.Reload(GlobalObject.ConfFilePath)
}
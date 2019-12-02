# Zinx复刻

*原项目地址： https://github.com/aceld/zinx. 本项目仅用于学习*

## 结构

Zinx - 轻量级TCP服务器框架

1. 基础Server
    1. 方法
        - 启动服务器 Start: 创建addr;创建listener;处理客户端连接并回显
        - 停止服务器 Stop
        - 运行服务器 Serve: 调用Start(); 并阻塞不让goroutine退出
        - 初始化 NewServer
    2. 属性
        - 名称 Name
        - 监听IP
        - 监听Port

2. 简单的连接封装和业务绑定
    1. 方法
        - 启动连接 Start
        - 停止连接 Stop
        - 获取当前连接对象（套接字）
        - 得到连接对象ID
        - 得到客户端连接对象的地址和端口
        - 发送数据的方法 Send
    2. 属性
        - TCP套接字
        - 连接对象ID
        - 当前连接状态（是否关闭）
        - 与当前连接绑定的处理业务方法
        - 等待连接被动退出的channel

3. 基础路由模块：
    1. 对Request进行封装，将连接和数据绑定在一起
        1. 方法
           - 获取当前连接
           - 获取当前数据
           - 新建一个Request请求
        2. 属性
           - 连接 IConnection
           - 请求数据
    2. 抽象的IRouter及基类实现
        - IRouter, 处理业务的主方法、处理业务之前的一些方法（钩子）、处理业务之后的方法
        - 具体的BaseRouter
        - router说明，框架中用BaseRouter实现IRouter接口，但用户往往需要自定义路由，所有用户可以自己实现IRouter接口，也可以组合BaseRouter并“重写”（屏蔽）原有方法实现。
        1. 方法
           - 处理业务之前的方法
           - 处理业务的主方法
           - 处理业务之后的方法
        2. 属性 无
    3. 集成路由功能：
        - IServer增加 路由添加 功能
        - Server类 增加 Router成员，把之前的handleAPI相关删除
        - Connection类绑定一个Router成员
        - Connection调用 已经注册的Router处理业务

4. 全局配置模块：
   1. 服务器应用/conf/zinx.json（由用户填写）
   2. 创建zinx全局配置模块/utils/globalobj.go
   3. 将zinx框架中所有硬编码替换成globalobj参数
   4. 使用zinx v0.4进行测试

5. 消息封装
    1. 定义一个消息的结构Message（属性有ID、长度、内容， 方法有set和get）
    2. 基于TLV（Type-Length-Value）解决消息封包与解包：数据包抽象类IDataPack与结构体DataPack实现。(方法：GetDataLen、Pack、Unpack)
    3. 将消息封装机制继承到zinx中，并作测试服务器开发：
        - 将Message添加到Request
        - 修改connection读取数据的机制，需要按照TLV格式拆包读取
        - 给connection发包增加pack

6. 多路由（消息管理）
   1.

## 实现

1. 抽象层定义IServer接口并在实体层以Server实现。
2. 抽象层定义IConnection接口并在实体层以Connection实现。
3. 抽象层定义IRequest接口并在实体层以Request实现; 抽象层定义IRouter并实现BaseRouter基类。
4. 全局变量GlobalObject，从json加载配置
5. Message/DataPack.
6.

# Zinx复刻

*原项目地址： https://github.com/aceld/zinx. 本项目仅用于学习*

## 结构

Zinx - 轻量级TCP服务器框架

1. 基础Server
    1.1 方法
        - 启动服务器 Start: 创建addr;创建listener;处理客户端连接并回显
        - 停止服务器 Stop
        - 运行服务器 Serve: 调用Start(); 并阻塞不让goroutine退出
        - 初始化 NewServer
    1.2 属性
        - 名称 Name
        - 监听IP
        - 监听Port


2. 简单的连接封装和业务绑定
    2.1 方法
        - 启动连接 Start
        - 停止连接 Stop
        - 获取当前连接对象（套接字）
        - 得到连接对象ID
        - 得到客户端连接对象的地址和端口
        - 发送数据的方法 Send
    2.2 属性
        - TCP套接字
        - 连接对象ID
        - 当前连接状态（是否关闭）
        - 与当前连接绑定的处理业务方法
        - 等待连接被动退出的channel

3. 基础路由模块：
    (1)对Request进行封装，将连接和数据绑定在一起
    3.1 方法
        - 获取当前连接
        - 获取当前数据
        - 新建一个Request请求
    3.2 属性
        - 连接 IConnection
        - 请求数据
    (2)抽象的IRouter，处理业务的主方法、处理业务之前的一些方法（钩子）、处理业务之后的方法，以及具体的BaseRouter
    3.3 方法
        - 处理业务之前的方法
        - 处理业务的主方法
        - 处理业务之后的方法
    3.4 属性 无
    (3) router说明，框架中用BaseRouter实现IRouter接口，但用户往往需要自定义路由，所有用户可以自己实现IRouter接口，
        也可以组合BaseRouter并“重写”（屏蔽）原有方法实现。


## 实现

1. 抽象层定义IServer接口并在实体层以Server实现。
2. 抽象层定义IConnection接口并在实体层以Connection实现。
3. 抽象层定义IRequest接口并在实体层以Request实现。
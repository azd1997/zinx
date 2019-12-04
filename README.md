# Zinx复刻

*原项目地址： <https://github.com/aceld/zinx>. 本项目仅用于学习*

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
   4. 使用zinx vchannel

5. 消息封装
    1. 定义一个消息的结构Message（属性有ID、长度、内容， 方法有set和get）
    2. 基于TLV（Type-Length-Value）解决消息封包与解包：数据包抽象类IDataPack与结构体DataPack实现。(方法：GetDataLen、Pack、Unpack)
    3. 将消息封装机制继承到zinx中，并作测试服务器开发：
        - 将Message添加到Request
        - 修改connection读取数据的机制，需要按照TLV格式拆包读取
        - 给connection发包增加pack

6. 多路由（消息管理）
   1. 抽象类IMsgHandler， 包含方法：DoMsgHandler、AddRouter
   2. 具体实现MsgHandler, 属性字段Apis (map[uint32]IRouter)
   3. 将connection属性中Router替换为MsgHandler并修改其NewConnection方法
   4. connection之前使用router的方法都相应修改为使用MsgHandler

7. 读写分离
   1. connection增加msgChan字段用来将读goroutine数据消息传给写goroutine
   2. startReader中之前读完之后写回客户端的操作交给startWriter去做
   3. 在connection.Start中增加 go startWriter，使读/写goroutine一同启动。

8. 消息队列与多任务处理机制
   1. 前面读写分离的设计是客户端发数据->服务端读数据->服务端写回数据的连续过程。每一个conn对应了一个reader和一个writer goroutine。当并发量高起来之后，不能无限制开很多个goroutine，切换开销会变得明显. 因此需要实现工作池机制。worker pool中有固定数量的worker goroutine（处理真正客户端请求的核心业务），worker pool 通过消息队列接收外界请求执行的的任务。每个worker配一个消息队列（也就是channel）
   2. 步骤：
      1. 创建一个消息队列
      2. 创建多任务worker的工作池并启动
      3. 将之前的发送消息全部改成将消息发给消息队列和worker工作池来处理

9. 链接管理
   1. 创建链接管理模块
       1. 当链接数较多时，为了保证服务端响应速度，应设置链接限制，拒绝过多的链接请求
       2. 创建链接管理模块IConnManager/ConnManager
   2. 将连接管理模块集成到zinx
       1. 将connmanager集成到server中，server中增加相应方法和属性，
       2. connection中增加server字段并在New方法中增加添加链接方法（将自己添加到server.connmanager中）
       3. server.start中增加对连接数量判断，满了则关闭新建立的链接（拒绝新连接）并continue
       4. connection.stop中要增加将自己从server.connmanager删除的操作
       5. server.stop中也应该将所有链接清空
   3. 给链接增加带缓冲的发包方法
       1. 以前无缓冲的msgChan发包如果客户端链接比较多而对方处理不及时可能产生阻塞，影响其他连接的处理，因此提供一个有缓冲的消息通道和非阻塞的发包方法是有必要的
       2. IConnection增加sendBuffMsg方法，Connection增加msgBuffChan和该方法，New作相应修改
       3. connection.startwriter中增加对msgBuffChan消息的处理
   4. 为链接启动和停止增加Hook方法
       1. IServer中增加两个Hook的set和call方法，共四个方法
       2. Server中增加两个Hook函数字段，并实现四个方法
       3. 在connection.start和stop中调用两个钩子函数
   5. 测试
       1. 在zinx_v_0_9.go中自定义这两个钩子方法，并注册到server中
       2. client代码不变

10. 链接属性设置
    1. 使用连接时，希望绑定一些和用户有关的参数数据，所以这一版本增加相应传递参数的接口方法
    2. IConnection增加SetProperty、GetProperty，RemoveProperty
    3. Connection增加一个map字段用来存储属性以及一个用来保护其的读写锁，并实现三个方法
    4. 测试服务器在连接启动和连接停止两个钩子函数中设置并读取连接属性

11. TODO: 整理代码

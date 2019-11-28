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


## 实现

抽象层定义IServer接口并在实体层实现其。
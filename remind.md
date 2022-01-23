1、网关层 gate
    对外提供一切所需的http,tcp,ws服务，协议转化，接收json、protobuf协议
2、服务层 service
    提供具体的服务逻辑, tcp服务, 仅接受protobuf协议
3、集群层 cluster
    提供服务注册与发现
4、c --- g --- sss, 服务与服务之间, 没有直接连接。网关层维护一个sid, 全局唯一。消息类型有：Request,Response,Push,Notice,Rpc

5、gate
5.1、注册(本地注册、远程注册)、心跳、关闭

6、cluster
6.1 分配用户到网关服务器 http
6.2 记录全部服务路由，并通知到各个网关 tcp

7、server - 实际服务层
7.1 支持独立运行

8、运行流程
8.1 启动server -- tcp
8.2 启动gate -- http/tcp
8.3 启动cluster -- http/tcp


10、go标准库 https://studygolang.com/pkgdoc

工作安排
1、heart包封装一下算球



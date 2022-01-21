1、网关层 gate
    对外提供一切所需的http,tcp,ws服务，协议转化，接收json、protobuf协议
2、服务层 service
    提供具体的服务逻辑, tcp服务, 仅接受protobuf协议
3、集群层 cluster
    提供服务注册与发现
4、c --- g --- sss, 服务与服务之间, 没有直接连接。网关层维护一个sid, 全局唯一。消息类型有：Request,Response,Push,Notice,Rpc

5、gate
5.1、注册(本地注册、远程注册)、心跳、关闭

工作安排
1、gate和服务层分别获取连接
2、rpc连接池  路由-->client
3、日志

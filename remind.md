1、网关层 gate
    对外提供一切所需的http,tcp,ws服务，协议转化，接收json、protobuf协议
2、服务层 service
    提供具体的服务逻辑, tcp服务, 仅接受protobuf协议
3、集群层 cluster
    提供服务注册与发现


工作安排
1、服务注册、发现  思路，节点获取自己的公网ip, 然后上报服务地址和对应服务
2、rpc连接池

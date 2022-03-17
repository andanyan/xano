# xano
xano是一款简易的tcp分布式微服务框架，支持服务最小粒度分散部署，你可以在短短几分钟的时间内就可以看懂和使用它。核心架构分为服务注册于发现层(master)、网关层(member)、服务层(server)三层，membe和server向master提交自己提供的服务和相关信息，而master则负责向两者同步数据节点信息以及其他数据。

## 协议类型
本框架协议类型支持protobuf和json, 默认protobuf, 你可以通过配置文件进行修改，配置文件请参考example下的cluster.toml文件

```
#基础配置
[base]
#版本号
version = "1.0.0"
#协议类型 1-protobuf 2-json
TcpDeal = 1

#网关节点
[Member]
# 服务启动地址
tcpAddr = ":10000"
# 内部交互地址
innerAddr = ":10001"
# 主节点访问地址
masterAddr = ":20000"

#主节点
[Master]
# 主节点服务启动地址
tcpAddr = ":20000"
# http服务，对外提供节点信息访问
httpAddr = ":20001"

#服务节点
[server]
# 本节点服务启动地址
tcpAddr = ":30000"
# 主节点访问地址
masterAddr = ":20000"

```

## 数据结构
### packet
packet是本框架消息传递的最小单位，包含两个字段，Length是Data的长度，tcp将根据length读取每个包的长度和数据，解决粘包问题。

```
type Packet struct {
	// 数据源长度
	Length uint16
	// 源数据
	Data []byte
}
```
### msg
msg是packet Data字段的数据结构

```
type Msg struct {
	// 通信SessionId-全局唯一
	Sid uint64
	// 路由-大驼峰式
	Route string
	// 消息id-从1开始递增
	Mid uint64
	// 消息类型-request/response/notice/push
	MsgType uint32
	// data数据协议
	Deal uint32-protobuf or json
	// 数据 - 数据
	Data []byte
	// 版本号 - 版本号
	Version string
}
```
## 消息类型
本框架消息类型分为4个，即request/response/notice/push，其中request和response成对出现，如果request之后没有回复response，程序将会卡死，请切记遵守该规则。notice用于客户端单方面给服务端发送消息，无需回包。push是服务端向客户端推送消息，无需回报。


## session
本框架的session仅提供单节点的数据缓存和消息处理，每个session从member层进入后，都会发放唯一的SID，客户端在连接上member时，会收到SessionInit的push推送。 你可以利用SID做用户数据的外部缓存，如登录后绑定用户的信息，便于后续的请求访问。
### session封装了Rpc/Notice/Response/Push/RpcResponse/PushTo等消息发送方法，同时提供给服务端GetSid/Get/Set等方法以供使用。 你可以在任一服务端，根据sid给任一客户端连接进行回包

## 路由
本框架路由采用注册的方式进行实现，你可以如下方式进行路由注册。路由名称是RouterServer的Name字段加上方法名，采用大驼峰式结构进行访问，如SessionInit。路由方法函数结构也是有严格要求的，必须是如下方式进行数据传递，不然就会报错，请严格遵守规则。完成的代码，在example下都可以看到，十分简单易懂。
```
package main

import (
	"xano"
	"xano/example/pb"
	"xano/logger"
	"xano/router"
	"xano/session"
)

type B struct{}

func (b *B) Div(s session.Session, input *pb.DivRequest) error {
	addRes := new(pb.AddResponse)
	err := s.Rpc("Add", &pb.AddRequest{
		Args: []int64{input.A, input.B},
	}, addRes)
	if err != nil {
		logger.Error(err)
		return err
	}

	res := addRes.Result * (input.B - input.A)

	return s.Response("Div", &pb.DivResponse{
		Result: res,
	})
}

func main() {
	xano.WithConfig("./config/server.toml")

	xano.WithRoute(&router.RouterServer{
		Name:   "",
		Server: new(B),
	})

	xano.Run()
}
```


## 版本控制
本框架严格要求version进行服务分发，客户端和服务端version不一致的话，请求将无法送达。如实在不一致，请先通过请求主节点的网关信息，同步version，客户端选择对应的服务

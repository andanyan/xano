// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.19.2
// source: msg.proto

package deal

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// 包消息
type Msg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// 路由 模块_方法
	// @inject_tag: json:"route"
	Route string `protobuf:"bytes,1,opt,name=route,proto3" json:"route,omitempty"`
	// 消息id
	// @inject_tag: json:"mid"
	Mid uint64 `protobuf:"varint,2,opt,name=mid,proto3" json:"mid,omitempty"`
	// 消息类型
	MsgType uint32 `protobuf:"varint,3,opt,name=msg_type,json=msgType,proto3" json:"msg_type,omitempty"`
	// data数据协议
	// @inject_tag: json:"deal"
	Deal uint32 `protobuf:"varint,4,opt,name=deal,proto3" json:"deal,omitempty"`
	// 数据
	// @inject_tag: json:"data"
	Data []byte `protobuf:"bytes,5,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *Msg) Reset() {
	*x = Msg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_msg_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Msg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Msg) ProtoMessage() {}

func (x *Msg) ProtoReflect() protoreflect.Message {
	mi := &file_msg_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Msg.ProtoReflect.Descriptor instead.
func (*Msg) Descriptor() ([]byte, []int) {
	return file_msg_proto_rawDescGZIP(), []int{0}
}

func (x *Msg) GetRoute() string {
	if x != nil {
		return x.Route
	}
	return ""
}

func (x *Msg) GetMid() uint64 {
	if x != nil {
		return x.Mid
	}
	return 0
}

func (x *Msg) GetMsgType() uint32 {
	if x != nil {
		return x.MsgType
	}
	return 0
}

func (x *Msg) GetDeal() uint32 {
	if x != nil {
		return x.Deal
	}
	return 0
}

func (x *Msg) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

// 路由注册
type ServerRunNotice struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// 版本号
	Version string `protobuf:"bytes,1,opt,name=version,proto3" json:"version,omitempty"`
	// 端口
	Port string `protobuf:"bytes,2,opt,name=port,proto3" json:"port,omitempty"`
	// 路由列表
	Routes []string `protobuf:"bytes,3,rep,name=routes,proto3" json:"routes,omitempty"`
}

func (x *ServerRunNotice) Reset() {
	*x = ServerRunNotice{}
	if protoimpl.UnsafeEnabled {
		mi := &file_msg_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ServerRunNotice) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ServerRunNotice) ProtoMessage() {}

func (x *ServerRunNotice) ProtoReflect() protoreflect.Message {
	mi := &file_msg_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ServerRunNotice.ProtoReflect.Descriptor instead.
func (*ServerRunNotice) Descriptor() ([]byte, []int) {
	return file_msg_proto_rawDescGZIP(), []int{1}
}

func (x *ServerRunNotice) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

func (x *ServerRunNotice) GetPort() string {
	if x != nil {
		return x.Port
	}
	return ""
}

func (x *ServerRunNotice) GetRoutes() []string {
	if x != nil {
		return x.Routes
	}
	return nil
}

// 心跳
type ServerHeartRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ServerHeartRequest) Reset() {
	*x = ServerHeartRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_msg_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ServerHeartRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ServerHeartRequest) ProtoMessage() {}

func (x *ServerHeartRequest) ProtoReflect() protoreflect.Message {
	mi := &file_msg_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ServerHeartRequest.ProtoReflect.Descriptor instead.
func (*ServerHeartRequest) Descriptor() ([]byte, []int) {
	return file_msg_proto_rawDescGZIP(), []int{2}
}

type ServerHeartResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ServerHeartResponse) Reset() {
	*x = ServerHeartResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_msg_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ServerHeartResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ServerHeartResponse) ProtoMessage() {}

func (x *ServerHeartResponse) ProtoReflect() protoreflect.Message {
	mi := &file_msg_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ServerHeartResponse.ProtoReflect.Descriptor instead.
func (*ServerHeartResponse) Descriptor() ([]byte, []int) {
	return file_msg_proto_rawDescGZIP(), []int{3}
}

// 路由关闭
type ServerCloseNotice struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Version string `protobuf:"bytes,1,opt,name=version,proto3" json:"version,omitempty"`
	Port    string `protobuf:"bytes,2,opt,name=port,proto3" json:"port,omitempty"`
}

func (x *ServerCloseNotice) Reset() {
	*x = ServerCloseNotice{}
	if protoimpl.UnsafeEnabled {
		mi := &file_msg_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ServerCloseNotice) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ServerCloseNotice) ProtoMessage() {}

func (x *ServerCloseNotice) ProtoReflect() protoreflect.Message {
	mi := &file_msg_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ServerCloseNotice.ProtoReflect.Descriptor instead.
func (*ServerCloseNotice) Descriptor() ([]byte, []int) {
	return file_msg_proto_rawDescGZIP(), []int{4}
}

func (x *ServerCloseNotice) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

func (x *ServerCloseNotice) GetPort() string {
	if x != nil {
		return x.Port
	}
	return ""
}

var File_msg_proto protoreflect.FileDescriptor

var file_msg_proto_rawDesc = []byte{
	0x0a, 0x09, 0x6d, 0x73, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x64, 0x65, 0x61,
	0x6c, 0x22, 0x70, 0x0a, 0x03, 0x4d, 0x73, 0x67, 0x12, 0x14, 0x0a, 0x05, 0x72, 0x6f, 0x75, 0x74,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x72, 0x6f, 0x75, 0x74, 0x65, 0x12, 0x10,
	0x0a, 0x03, 0x6d, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x03, 0x6d, 0x69, 0x64,
	0x12, 0x19, 0x0a, 0x08, 0x6d, 0x73, 0x67, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x0d, 0x52, 0x07, 0x6d, 0x73, 0x67, 0x54, 0x79, 0x70, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x64,
	0x65, 0x61, 0x6c, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x04, 0x64, 0x65, 0x61, 0x6c, 0x12,
	0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64,
	0x61, 0x74, 0x61, 0x22, 0x57, 0x0a, 0x0f, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x52, 0x75, 0x6e,
	0x4e, 0x6f, 0x74, 0x69, 0x63, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f,
	0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e,
	0x12, 0x12, 0x0a, 0x04, 0x70, 0x6f, 0x72, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x70, 0x6f, 0x72, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x6f, 0x75, 0x74, 0x65, 0x73, 0x18, 0x03,
	0x20, 0x03, 0x28, 0x09, 0x52, 0x06, 0x72, 0x6f, 0x75, 0x74, 0x65, 0x73, 0x22, 0x14, 0x0a, 0x12,
	0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x48, 0x65, 0x61, 0x72, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x22, 0x15, 0x0a, 0x13, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x48, 0x65, 0x61, 0x72,
	0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x41, 0x0a, 0x11, 0x53, 0x65, 0x72,
	0x76, 0x65, 0x72, 0x43, 0x6c, 0x6f, 0x73, 0x65, 0x4e, 0x6f, 0x74, 0x69, 0x63, 0x65, 0x12, 0x18,
	0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x6f, 0x72, 0x74,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x70, 0x6f, 0x72, 0x74, 0x42, 0x09, 0x5a, 0x07,
	0x2e, 0x2e, 0x2f, 0x64, 0x65, 0x61, 0x6c, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_msg_proto_rawDescOnce sync.Once
	file_msg_proto_rawDescData = file_msg_proto_rawDesc
)

func file_msg_proto_rawDescGZIP() []byte {
	file_msg_proto_rawDescOnce.Do(func() {
		file_msg_proto_rawDescData = protoimpl.X.CompressGZIP(file_msg_proto_rawDescData)
	})
	return file_msg_proto_rawDescData
}

var file_msg_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_msg_proto_goTypes = []interface{}{
	(*Msg)(nil),                 // 0: deal.Msg
	(*ServerRunNotice)(nil),     // 1: deal.ServerRunNotice
	(*ServerHeartRequest)(nil),  // 2: deal.ServerHeartRequest
	(*ServerHeartResponse)(nil), // 3: deal.ServerHeartResponse
	(*ServerCloseNotice)(nil),   // 4: deal.ServerCloseNotice
}
var file_msg_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_msg_proto_init() }
func file_msg_proto_init() {
	if File_msg_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_msg_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Msg); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_msg_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ServerRunNotice); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_msg_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ServerHeartRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_msg_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ServerHeartResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_msg_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ServerCloseNotice); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_msg_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_msg_proto_goTypes,
		DependencyIndexes: file_msg_proto_depIdxs,
		MessageInfos:      file_msg_proto_msgTypes,
	}.Build()
	File_msg_proto = out.File
	file_msg_proto_rawDesc = nil
	file_msg_proto_goTypes = nil
	file_msg_proto_depIdxs = nil
}
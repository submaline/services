// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        (unknown)
// source: talk/v1/talk.proto

package talkv1

import (
	v1 "github.com/submaline/services/gen/types/v1"
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

type SendMessageRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message *v1.Message `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *SendMessageRequest) Reset() {
	*x = SendMessageRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_talk_v1_talk_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SendMessageRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendMessageRequest) ProtoMessage() {}

func (x *SendMessageRequest) ProtoReflect() protoreflect.Message {
	mi := &file_talk_v1_talk_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendMessageRequest.ProtoReflect.Descriptor instead.
func (*SendMessageRequest) Descriptor() ([]byte, []int) {
	return file_talk_v1_talk_proto_rawDescGZIP(), []int{0}
}

func (x *SendMessageRequest) GetMessage() *v1.Message {
	if x != nil {
		return x.Message
	}
	return nil
}

type SendMessageResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message *v1.Message `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *SendMessageResponse) Reset() {
	*x = SendMessageResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_talk_v1_talk_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SendMessageResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendMessageResponse) ProtoMessage() {}

func (x *SendMessageResponse) ProtoReflect() protoreflect.Message {
	mi := &file_talk_v1_talk_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendMessageResponse.ProtoReflect.Descriptor instead.
func (*SendMessageResponse) Descriptor() ([]byte, []int) {
	return file_talk_v1_talk_proto_rawDescGZIP(), []int{1}
}

func (x *SendMessageResponse) GetMessage() *v1.Message {
	if x != nil {
		return x.Message
	}
	return nil
}

type SendReadReceiptRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MessageId string `protobuf:"bytes,1,opt,name=message_id,json=messageId,proto3" json:"message_id,omitempty"`
}

func (x *SendReadReceiptRequest) Reset() {
	*x = SendReadReceiptRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_talk_v1_talk_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SendReadReceiptRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendReadReceiptRequest) ProtoMessage() {}

func (x *SendReadReceiptRequest) ProtoReflect() protoreflect.Message {
	mi := &file_talk_v1_talk_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendReadReceiptRequest.ProtoReflect.Descriptor instead.
func (*SendReadReceiptRequest) Descriptor() ([]byte, []int) {
	return file_talk_v1_talk_proto_rawDescGZIP(), []int{2}
}

func (x *SendReadReceiptRequest) GetMessageId() string {
	if x != nil {
		return x.MessageId
	}
	return ""
}

type SendReadReceiptResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *SendReadReceiptResponse) Reset() {
	*x = SendReadReceiptResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_talk_v1_talk_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SendReadReceiptResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendReadReceiptResponse) ProtoMessage() {}

func (x *SendReadReceiptResponse) ProtoReflect() protoreflect.Message {
	mi := &file_talk_v1_talk_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendReadReceiptResponse.ProtoReflect.Descriptor instead.
func (*SendReadReceiptResponse) Descriptor() ([]byte, []int) {
	return file_talk_v1_talk_proto_rawDescGZIP(), []int{3}
}

var File_talk_v1_talk_proto protoreflect.FileDescriptor

var file_talk_v1_talk_proto_rawDesc = []byte{
	0x0a, 0x12, 0x74, 0x61, 0x6c, 0x6b, 0x2f, 0x76, 0x31, 0x2f, 0x74, 0x61, 0x6c, 0x6b, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x74, 0x61, 0x6c, 0x6b, 0x2e, 0x76, 0x31, 0x1a, 0x14, 0x74,
	0x79, 0x70, 0x65, 0x73, 0x2f, 0x76, 0x31, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0x41, 0x0a, 0x12, 0x53, 0x65, 0x6e, 0x64, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2b, 0x0a, 0x07, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x74, 0x79, 0x70,
	0x65, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x07, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x42, 0x0a, 0x13, 0x53, 0x65, 0x6e, 0x64, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2b, 0x0a,
	0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11,
	0x2e, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x37, 0x0a, 0x16, 0x53, 0x65,
	0x6e, 0x64, 0x52, 0x65, 0x61, 0x64, 0x52, 0x65, 0x63, 0x65, 0x69, 0x70, 0x74, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x5f,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x49, 0x64, 0x22, 0x19, 0x0a, 0x17, 0x53, 0x65, 0x6e, 0x64, 0x52, 0x65, 0x61, 0x64, 0x52,
	0x65, 0x63, 0x65, 0x69, 0x70, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x32, 0xb1,
	0x01, 0x0a, 0x0b, 0x54, 0x61, 0x6c, 0x6b, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x4a,
	0x0a, 0x0b, 0x53, 0x65, 0x6e, 0x64, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x1b, 0x2e,
	0x74, 0x61, 0x6c, 0x6b, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x65, 0x6e, 0x64, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e, 0x74, 0x61, 0x6c,
	0x6b, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x65, 0x6e, 0x64, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x56, 0x0a, 0x0f, 0x53, 0x65,
	0x6e, 0x64, 0x52, 0x65, 0x61, 0x64, 0x52, 0x65, 0x63, 0x65, 0x69, 0x70, 0x74, 0x12, 0x1f, 0x2e,
	0x74, 0x61, 0x6c, 0x6b, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x65, 0x6e, 0x64, 0x52, 0x65, 0x61, 0x64,
	0x52, 0x65, 0x63, 0x65, 0x69, 0x70, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x20,
	0x2e, 0x74, 0x61, 0x6c, 0x6b, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x65, 0x6e, 0x64, 0x52, 0x65, 0x61,
	0x64, 0x52, 0x65, 0x63, 0x65, 0x69, 0x70, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x00, 0x42, 0x87, 0x01, 0x0a, 0x0b, 0x63, 0x6f, 0x6d, 0x2e, 0x74, 0x61, 0x6c, 0x6b, 0x2e,
	0x76, 0x31, 0x42, 0x09, 0x54, 0x61, 0x6c, 0x6b, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a,
	0x30, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x75, 0x62, 0x6d,
	0x61, 0x6c, 0x69, 0x6e, 0x65, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2f, 0x67,
	0x65, 0x6e, 0x2f, 0x74, 0x61, 0x6c, 0x6b, 0x2f, 0x76, 0x31, 0x3b, 0x74, 0x61, 0x6c, 0x6b, 0x76,
	0x31, 0xa2, 0x02, 0x03, 0x54, 0x58, 0x58, 0xaa, 0x02, 0x07, 0x54, 0x61, 0x6c, 0x6b, 0x2e, 0x56,
	0x31, 0xca, 0x02, 0x07, 0x54, 0x61, 0x6c, 0x6b, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x13, 0x54, 0x61,
	0x6c, 0x6b, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74,
	0x61, 0xea, 0x02, 0x08, 0x54, 0x61, 0x6c, 0x6b, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_talk_v1_talk_proto_rawDescOnce sync.Once
	file_talk_v1_talk_proto_rawDescData = file_talk_v1_talk_proto_rawDesc
)

func file_talk_v1_talk_proto_rawDescGZIP() []byte {
	file_talk_v1_talk_proto_rawDescOnce.Do(func() {
		file_talk_v1_talk_proto_rawDescData = protoimpl.X.CompressGZIP(file_talk_v1_talk_proto_rawDescData)
	})
	return file_talk_v1_talk_proto_rawDescData
}

var file_talk_v1_talk_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_talk_v1_talk_proto_goTypes = []interface{}{
	(*SendMessageRequest)(nil),      // 0: talk.v1.SendMessageRequest
	(*SendMessageResponse)(nil),     // 1: talk.v1.SendMessageResponse
	(*SendReadReceiptRequest)(nil),  // 2: talk.v1.SendReadReceiptRequest
	(*SendReadReceiptResponse)(nil), // 3: talk.v1.SendReadReceiptResponse
	(*v1.Message)(nil),              // 4: types.v1.Message
}
var file_talk_v1_talk_proto_depIdxs = []int32{
	4, // 0: talk.v1.SendMessageRequest.message:type_name -> types.v1.Message
	4, // 1: talk.v1.SendMessageResponse.message:type_name -> types.v1.Message
	0, // 2: talk.v1.TalkService.SendMessage:input_type -> talk.v1.SendMessageRequest
	2, // 3: talk.v1.TalkService.SendReadReceipt:input_type -> talk.v1.SendReadReceiptRequest
	1, // 4: talk.v1.TalkService.SendMessage:output_type -> talk.v1.SendMessageResponse
	3, // 5: talk.v1.TalkService.SendReadReceipt:output_type -> talk.v1.SendReadReceiptResponse
	4, // [4:6] is the sub-list for method output_type
	2, // [2:4] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_talk_v1_talk_proto_init() }
func file_talk_v1_talk_proto_init() {
	if File_talk_v1_talk_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_talk_v1_talk_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SendMessageRequest); i {
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
		file_talk_v1_talk_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SendMessageResponse); i {
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
		file_talk_v1_talk_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SendReadReceiptRequest); i {
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
		file_talk_v1_talk_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SendReadReceiptResponse); i {
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
			RawDescriptor: file_talk_v1_talk_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_talk_v1_talk_proto_goTypes,
		DependencyIndexes: file_talk_v1_talk_proto_depIdxs,
		MessageInfos:      file_talk_v1_talk_proto_msgTypes,
	}.Build()
	File_talk_v1_talk_proto = out.File
	file_talk_v1_talk_proto_rawDesc = nil
	file_talk_v1_talk_proto_goTypes = nil
	file_talk_v1_talk_proto_depIdxs = nil
}

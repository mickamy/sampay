// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.1
// 	protoc        (unknown)
// source: auth/v1/common.proto

package authv1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Token struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Value         string                 `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
	ExpiresAt     *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=expires_at,json=expiresAt,proto3" json:"expires_at,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Token) Reset() {
	*x = Token{}
	mi := &file_auth_v1_common_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Token) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Token) ProtoMessage() {}

func (x *Token) ProtoReflect() protoreflect.Message {
	mi := &file_auth_v1_common_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Token.ProtoReflect.Descriptor instead.
func (*Token) Descriptor() ([]byte, []int) {
	return file_auth_v1_common_proto_rawDescGZIP(), []int{0}
}

func (x *Token) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

func (x *Token) GetExpiresAt() *timestamppb.Timestamp {
	if x != nil {
		return x.ExpiresAt
	}
	return nil
}

type Tokens struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Access        *Token                 `protobuf:"bytes,1,opt,name=access,proto3" json:"access,omitempty"`
	Refresh       *Token                 `protobuf:"bytes,2,opt,name=refresh,proto3" json:"refresh,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Tokens) Reset() {
	*x = Tokens{}
	mi := &file_auth_v1_common_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Tokens) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Tokens) ProtoMessage() {}

func (x *Tokens) ProtoReflect() protoreflect.Message {
	mi := &file_auth_v1_common_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Tokens.ProtoReflect.Descriptor instead.
func (*Tokens) Descriptor() ([]byte, []int) {
	return file_auth_v1_common_proto_rawDescGZIP(), []int{1}
}

func (x *Tokens) GetAccess() *Token {
	if x != nil {
		return x.Access
	}
	return nil
}

func (x *Tokens) GetRefresh() *Token {
	if x != nil {
		return x.Refresh
	}
	return nil
}

var File_auth_v1_common_proto protoreflect.FileDescriptor

var file_auth_v1_common_proto_rawDesc = []byte{
	0x0a, 0x14, 0x61, 0x75, 0x74, 0x68, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x76, 0x31, 0x1a,
	0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0x58, 0x0a, 0x05, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x12,
	0x39, 0x0a, 0x0a, 0x65, 0x78, 0x70, 0x69, 0x72, 0x65, 0x73, 0x5f, 0x61, 0x74, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52,
	0x09, 0x65, 0x78, 0x70, 0x69, 0x72, 0x65, 0x73, 0x41, 0x74, 0x22, 0x5a, 0x0a, 0x06, 0x54, 0x6f,
	0x6b, 0x65, 0x6e, 0x73, 0x12, 0x26, 0x0a, 0x06, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x76, 0x31, 0x2e, 0x54,
	0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x06, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12, 0x28, 0x0a, 0x07,
	0x72, 0x65, 0x66, 0x72, 0x65, 0x73, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e,
	0x61, 0x75, 0x74, 0x68, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x07, 0x72,
	0x65, 0x66, 0x72, 0x65, 0x73, 0x68, 0x42, 0x7e, 0x0a, 0x0b, 0x63, 0x6f, 0x6d, 0x2e, 0x61, 0x75,
	0x74, 0x68, 0x2e, 0x76, 0x31, 0x42, 0x0b, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x50, 0x72, 0x6f,
	0x74, 0x6f, 0x50, 0x01, 0x5a, 0x25, 0x6d, 0x69, 0x63, 0x6b, 0x61, 0x6d, 0x79, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x73, 0x61, 0x6d, 0x70, 0x61, 0x79, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x61, 0x75, 0x74,
	0x68, 0x2f, 0x76, 0x31, 0x3b, 0x61, 0x75, 0x74, 0x68, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x41, 0x58,
	0x58, 0xaa, 0x02, 0x07, 0x41, 0x75, 0x74, 0x68, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x07, 0x41, 0x75,
	0x74, 0x68, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x13, 0x41, 0x75, 0x74, 0x68, 0x5c, 0x56, 0x31, 0x5c,
	0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x08, 0x41, 0x75,
	0x74, 0x68, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_auth_v1_common_proto_rawDescOnce sync.Once
	file_auth_v1_common_proto_rawDescData = file_auth_v1_common_proto_rawDesc
)

func file_auth_v1_common_proto_rawDescGZIP() []byte {
	file_auth_v1_common_proto_rawDescOnce.Do(func() {
		file_auth_v1_common_proto_rawDescData = protoimpl.X.CompressGZIP(file_auth_v1_common_proto_rawDescData)
	})
	return file_auth_v1_common_proto_rawDescData
}

var file_auth_v1_common_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_auth_v1_common_proto_goTypes = []any{
	(*Token)(nil),                 // 0: auth.v1.Token
	(*Tokens)(nil),                // 1: auth.v1.Tokens
	(*timestamppb.Timestamp)(nil), // 2: google.protobuf.Timestamp
}
var file_auth_v1_common_proto_depIdxs = []int32{
	2, // 0: auth.v1.Token.expires_at:type_name -> google.protobuf.Timestamp
	0, // 1: auth.v1.Tokens.access:type_name -> auth.v1.Token
	0, // 2: auth.v1.Tokens.refresh:type_name -> auth.v1.Token
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_auth_v1_common_proto_init() }
func file_auth_v1_common_proto_init() {
	if File_auth_v1_common_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_auth_v1_common_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_auth_v1_common_proto_goTypes,
		DependencyIndexes: file_auth_v1_common_proto_depIdxs,
		MessageInfos:      file_auth_v1_common_proto_msgTypes,
	}.Build()
	File_auth_v1_common_proto = out.File
	file_auth_v1_common_proto_rawDesc = nil
	file_auth_v1_common_proto_goTypes = nil
	file_auth_v1_common_proto_depIdxs = nil
}

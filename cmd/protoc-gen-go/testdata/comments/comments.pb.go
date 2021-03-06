// Code generated by protoc-gen-go. DO NOT EDIT.
// source: comments/comments.proto

// COMMENT: package goproto.protoc.comments;

package comments

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoregistry "google.golang.org/protobuf/reflect/protoregistry"
	protoiface "google.golang.org/protobuf/runtime/protoiface"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	sync "sync"
)

const _ = protoimpl.EnforceVersion(protoimpl.Version - 0)

// COMMENT: Message1
type Message1 struct {
	// COMMENT: Field1A
	Field1A *string `protobuf:"bytes,1,opt,name=Field1A" json:"Field1A,omitempty"`
	// COMMENT: Oneof1A
	//
	// Types that are valid to be assigned to Oneof1A:
	// COMMENT: Oneof1AField1
	//	*Message1_Oneof1AField1
	Oneof1A              isMessage1_Oneof1A      `protobuf_oneof:"Oneof1a"`
	XXX_NoUnkeyedLiteral struct{}                `json:"-"`
	XXX_unrecognized     protoimpl.UnknownFields `json:"-"`
	XXX_sizecache        protoimpl.SizeCache     `json:"-"`
}

func (x *Message1) Reset() {
	*x = Message1{}
}

func (x *Message1) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message1) ProtoMessage() {}

func (x *Message1) ProtoReflect() protoreflect.Message {
	return file_comments_comments_proto_msgTypes[0].MessageOf(x)
}

func (m *Message1) XXX_Methods() *protoiface.Methods {
	return file_comments_comments_proto_msgTypes[0].Methods()
}

// Deprecated: Use Message1.ProtoReflect.Type instead.
func (*Message1) Descriptor() ([]byte, []int) {
	return file_comments_comments_proto_rawDescGZIP(), []int{0}
}

func (x *Message1) GetField1A() string {
	if x != nil && x.Field1A != nil {
		return *x.Field1A
	}
	return ""
}

func (m *Message1) GetOneof1A() isMessage1_Oneof1A {
	if m != nil {
		return m.Oneof1A
	}
	return nil
}

func (x *Message1) GetOneof1AField1() string {
	if x, ok := x.GetOneof1A().(*Message1_Oneof1AField1); ok {
		return x.Oneof1AField1
	}
	return ""
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*Message1) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*Message1_Oneof1AField1)(nil),
	}
}

type isMessage1_Oneof1A interface {
	isMessage1_Oneof1A()
}

type Message1_Oneof1AField1 struct {
	Oneof1AField1 string `protobuf:"bytes,2,opt,name=Oneof1AField1,oneof"`
}

func (*Message1_Oneof1AField1) isMessage1_Oneof1A() {}

// COMMENT: Message2
type Message2 struct {
	XXX_NoUnkeyedLiteral struct{}                `json:"-"`
	XXX_unrecognized     protoimpl.UnknownFields `json:"-"`
	XXX_sizecache        protoimpl.SizeCache     `json:"-"`
}

func (x *Message2) Reset() {
	*x = Message2{}
}

func (x *Message2) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message2) ProtoMessage() {}

func (x *Message2) ProtoReflect() protoreflect.Message {
	return file_comments_comments_proto_msgTypes[1].MessageOf(x)
}

func (m *Message2) XXX_Methods() *protoiface.Methods {
	return file_comments_comments_proto_msgTypes[1].Methods()
}

// Deprecated: Use Message2.ProtoReflect.Type instead.
func (*Message2) Descriptor() ([]byte, []int) {
	return file_comments_comments_proto_rawDescGZIP(), []int{1}
}

// COMMENT: Message1A
type Message1_Message1A struct {
	XXX_NoUnkeyedLiteral struct{}                `json:"-"`
	XXX_unrecognized     protoimpl.UnknownFields `json:"-"`
	XXX_sizecache        protoimpl.SizeCache     `json:"-"`
}

func (x *Message1_Message1A) Reset() {
	*x = Message1_Message1A{}
}

func (x *Message1_Message1A) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message1_Message1A) ProtoMessage() {}

func (x *Message1_Message1A) ProtoReflect() protoreflect.Message {
	return file_comments_comments_proto_msgTypes[2].MessageOf(x)
}

func (m *Message1_Message1A) XXX_Methods() *protoiface.Methods {
	return file_comments_comments_proto_msgTypes[2].Methods()
}

// Deprecated: Use Message1_Message1A.ProtoReflect.Type instead.
func (*Message1_Message1A) Descriptor() ([]byte, []int) {
	return file_comments_comments_proto_rawDescGZIP(), []int{0, 0}
}

// COMMENT: Message1B
type Message1_Message1B struct {
	XXX_NoUnkeyedLiteral struct{}                `json:"-"`
	XXX_unrecognized     protoimpl.UnknownFields `json:"-"`
	XXX_sizecache        protoimpl.SizeCache     `json:"-"`
}

func (x *Message1_Message1B) Reset() {
	*x = Message1_Message1B{}
}

func (x *Message1_Message1B) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message1_Message1B) ProtoMessage() {}

func (x *Message1_Message1B) ProtoReflect() protoreflect.Message {
	return file_comments_comments_proto_msgTypes[3].MessageOf(x)
}

func (m *Message1_Message1B) XXX_Methods() *protoiface.Methods {
	return file_comments_comments_proto_msgTypes[3].Methods()
}

// Deprecated: Use Message1_Message1B.ProtoReflect.Type instead.
func (*Message1_Message1B) Descriptor() ([]byte, []int) {
	return file_comments_comments_proto_rawDescGZIP(), []int{0, 1}
}

// COMMENT: Message2A
type Message2_Message2A struct {
	XXX_NoUnkeyedLiteral struct{}                `json:"-"`
	XXX_unrecognized     protoimpl.UnknownFields `json:"-"`
	XXX_sizecache        protoimpl.SizeCache     `json:"-"`
}

func (x *Message2_Message2A) Reset() {
	*x = Message2_Message2A{}
}

func (x *Message2_Message2A) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message2_Message2A) ProtoMessage() {}

func (x *Message2_Message2A) ProtoReflect() protoreflect.Message {
	return file_comments_comments_proto_msgTypes[4].MessageOf(x)
}

func (m *Message2_Message2A) XXX_Methods() *protoiface.Methods {
	return file_comments_comments_proto_msgTypes[4].Methods()
}

// Deprecated: Use Message2_Message2A.ProtoReflect.Type instead.
func (*Message2_Message2A) Descriptor() ([]byte, []int) {
	return file_comments_comments_proto_rawDescGZIP(), []int{1, 0}
}

// COMMENT: Message2B
type Message2_Message2B struct {
	XXX_NoUnkeyedLiteral struct{}                `json:"-"`
	XXX_unrecognized     protoimpl.UnknownFields `json:"-"`
	XXX_sizecache        protoimpl.SizeCache     `json:"-"`
}

func (x *Message2_Message2B) Reset() {
	*x = Message2_Message2B{}
}

func (x *Message2_Message2B) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message2_Message2B) ProtoMessage() {}

func (x *Message2_Message2B) ProtoReflect() protoreflect.Message {
	return file_comments_comments_proto_msgTypes[5].MessageOf(x)
}

func (m *Message2_Message2B) XXX_Methods() *protoiface.Methods {
	return file_comments_comments_proto_msgTypes[5].Methods()
}

// Deprecated: Use Message2_Message2B.ProtoReflect.Type instead.
func (*Message2_Message2B) Descriptor() ([]byte, []int) {
	return file_comments_comments_proto_rawDescGZIP(), []int{1, 1}
}

var File_comments_comments_proto protoreflect.FileDescriptor

var file_comments_comments_proto_rawDesc = []byte{
	0x0a, 0x17, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x65,
	0x6e, 0x74, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x17, 0x67, 0x6f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e,
	0x74, 0x73, 0x22, 0x71, 0x0a, 0x08, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x31, 0x12, 0x18,
	0x0a, 0x07, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x31, 0x41, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x07, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x31, 0x41, 0x12, 0x26, 0x0a, 0x0d, 0x4f, 0x6e, 0x65, 0x6f,
	0x66, 0x31, 0x41, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x31, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x48,
	0x00, 0x52, 0x0d, 0x4f, 0x6e, 0x65, 0x6f, 0x66, 0x31, 0x41, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x31,
	0x1a, 0x0b, 0x0a, 0x09, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x31, 0x41, 0x1a, 0x0b, 0x0a,
	0x09, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x31, 0x42, 0x42, 0x09, 0x0a, 0x07, 0x4f, 0x6e,
	0x65, 0x6f, 0x66, 0x31, 0x61, 0x22, 0x24, 0x0a, 0x08, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x32, 0x1a, 0x0b, 0x0a, 0x09, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x32, 0x41, 0x1a, 0x0b,
	0x0a, 0x09, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x32, 0x42, 0x42, 0x40, 0x5a, 0x3e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x67, 0x6f, 0x6c, 0x61, 0x6e, 0x67, 0x2e, 0x6f, 0x72, 0x67,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x63, 0x6d, 0x64, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x63, 0x2d, 0x67, 0x65, 0x6e, 0x2d, 0x67, 0x6f, 0x2f, 0x74, 0x65, 0x73, 0x74,
	0x64, 0x61, 0x74, 0x61, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x73,
}

var (
	file_comments_comments_proto_rawDescOnce sync.Once
	file_comments_comments_proto_rawDescData = file_comments_comments_proto_rawDesc
)

func file_comments_comments_proto_rawDescGZIP() []byte {
	file_comments_comments_proto_rawDescOnce.Do(func() {
		file_comments_comments_proto_rawDescData = protoimpl.X.CompressGZIP(file_comments_comments_proto_rawDescData)
	})
	return file_comments_comments_proto_rawDescData
}

var file_comments_comments_proto_msgTypes = make([]protoimpl.MessageType, 6)
var file_comments_comments_proto_goTypes = []interface{}{
	(*Message1)(nil),           // 0: goproto.protoc.comments.Message1
	(*Message2)(nil),           // 1: goproto.protoc.comments.Message2
	(*Message1_Message1A)(nil), // 2: goproto.protoc.comments.Message1.Message1A
	(*Message1_Message1B)(nil), // 3: goproto.protoc.comments.Message1.Message1B
	(*Message2_Message2A)(nil), // 4: goproto.protoc.comments.Message2.Message2A
	(*Message2_Message2B)(nil), // 5: goproto.protoc.comments.Message2.Message2B
}
var file_comments_comments_proto_depIdxs = []int32{}

func init() { file_comments_comments_proto_init() }
func file_comments_comments_proto_init() {
	if File_comments_comments_proto != nil {
		return
	}
	File_comments_comments_proto = protoimpl.FileBuilder{
		RawDescriptor:      file_comments_comments_proto_rawDesc,
		GoTypes:            file_comments_comments_proto_goTypes,
		DependencyIndexes:  file_comments_comments_proto_depIdxs,
		MessageOutputTypes: file_comments_comments_proto_msgTypes,
		FilesRegistry:      protoregistry.GlobalFiles,
		TypesRegistry:      protoregistry.GlobalTypes,
	}.Init()
	file_comments_comments_proto_rawDesc = nil
	file_comments_comments_proto_goTypes = nil
	file_comments_comments_proto_depIdxs = nil
}

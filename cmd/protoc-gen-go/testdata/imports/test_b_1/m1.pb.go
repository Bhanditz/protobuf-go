// Code generated by protoc-gen-go. DO NOT EDIT.
// source: imports/test_b_1/m1.proto

package beta

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoregistry "google.golang.org/protobuf/reflect/protoregistry"
	protoiface "google.golang.org/protobuf/runtime/protoiface"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	sync "sync"
)

const _ = protoimpl.EnforceVersion(protoimpl.Version - 0)

type M1 struct {
	XXX_NoUnkeyedLiteral struct{}                `json:"-"`
	XXX_unrecognized     protoimpl.UnknownFields `json:"-"`
	XXX_sizecache        protoimpl.SizeCache     `json:"-"`
}

func (x *M1) Reset() {
	*x = M1{}
}

func (x *M1) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*M1) ProtoMessage() {}

func (x *M1) ProtoReflect() protoreflect.Message {
	return file_imports_test_b_1_m1_proto_msgTypes[0].MessageOf(x)
}

func (m *M1) XXX_Methods() *protoiface.Methods {
	return file_imports_test_b_1_m1_proto_msgTypes[0].Methods()
}

// Deprecated: Use M1.ProtoReflect.Type instead.
func (*M1) Descriptor() ([]byte, []int) {
	return file_imports_test_b_1_m1_proto_rawDescGZIP(), []int{0}
}

var File_imports_test_b_1_m1_proto protoreflect.FileDescriptor

var file_imports_test_b_1_m1_proto_rawDesc = []byte{
	0x0a, 0x19, 0x69, 0x6d, 0x70, 0x6f, 0x72, 0x74, 0x73, 0x2f, 0x74, 0x65, 0x73, 0x74, 0x5f, 0x62,
	0x5f, 0x31, 0x2f, 0x6d, 0x31, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0c, 0x74, 0x65, 0x73,
	0x74, 0x2e, 0x62, 0x2e, 0x70, 0x61, 0x72, 0x74, 0x31, 0x22, 0x04, 0x0a, 0x02, 0x4d, 0x31, 0x42,
	0x4d, 0x5a, 0x4b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x67, 0x6f, 0x6c, 0x61, 0x6e, 0x67,
	0x2e, 0x6f, 0x72, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x63, 0x6d,
	0x64, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x2d, 0x67, 0x65, 0x6e, 0x2d, 0x67, 0x6f, 0x2f,
	0x74, 0x65, 0x73, 0x74, 0x64, 0x61, 0x74, 0x61, 0x2f, 0x69, 0x6d, 0x70, 0x6f, 0x72, 0x74, 0x73,
	0x2f, 0x74, 0x65, 0x73, 0x74, 0x5f, 0x62, 0x5f, 0x31, 0x3b, 0x62, 0x65, 0x74, 0x61, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_imports_test_b_1_m1_proto_rawDescOnce sync.Once
	file_imports_test_b_1_m1_proto_rawDescData = file_imports_test_b_1_m1_proto_rawDesc
)

func file_imports_test_b_1_m1_proto_rawDescGZIP() []byte {
	file_imports_test_b_1_m1_proto_rawDescOnce.Do(func() {
		file_imports_test_b_1_m1_proto_rawDescData = protoimpl.X.CompressGZIP(file_imports_test_b_1_m1_proto_rawDescData)
	})
	return file_imports_test_b_1_m1_proto_rawDescData
}

var file_imports_test_b_1_m1_proto_msgTypes = make([]protoimpl.MessageType, 1)
var file_imports_test_b_1_m1_proto_goTypes = []interface{}{
	(*M1)(nil), // 0: test.b.part1.M1
}
var file_imports_test_b_1_m1_proto_depIdxs = []int32{}

func init() { file_imports_test_b_1_m1_proto_init() }
func file_imports_test_b_1_m1_proto_init() {
	if File_imports_test_b_1_m1_proto != nil {
		return
	}
	File_imports_test_b_1_m1_proto = protoimpl.FileBuilder{
		RawDescriptor:      file_imports_test_b_1_m1_proto_rawDesc,
		GoTypes:            file_imports_test_b_1_m1_proto_goTypes,
		DependencyIndexes:  file_imports_test_b_1_m1_proto_depIdxs,
		MessageOutputTypes: file_imports_test_b_1_m1_proto_msgTypes,
		FilesRegistry:      protoregistry.GlobalFiles,
		TypesRegistry:      protoregistry.GlobalTypes,
	}.Init()
	file_imports_test_b_1_m1_proto_rawDesc = nil
	file_imports_test_b_1_m1_proto_goTypes = nil
	file_imports_test_b_1_m1_proto_depIdxs = nil
}

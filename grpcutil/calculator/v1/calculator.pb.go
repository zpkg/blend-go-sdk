/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/



// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.12.3
// source: calculator.proto

package v1

import (
	reflect "reflect"
	sync "sync"

	proto "github.com/golang/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type Numbers struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Values []float64 `protobuf:"fixed64,1,rep,packed,name=Values,proto3" json:"Values,omitempty"`
}

func (x *Numbers) Reset() {
	*x = Numbers{}
	if protoimpl.UnsafeEnabled {
		mi := &file_calculator_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Numbers) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Numbers) ProtoMessage() {}

func (x *Numbers) ProtoReflect() protoreflect.Message {
	mi := &file_calculator_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Numbers.ProtoReflect.Descriptor instead.
func (*Numbers) Descriptor() ([]byte, []int) {
	return file_calculator_proto_rawDescGZIP(), []int{0}
}

func (x *Numbers) GetValues() []float64 {
	if x != nil {
		return x.Values
	}
	return nil
}

type Number struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value float64 `protobuf:"fixed64,1,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *Number) Reset() {
	*x = Number{}
	if protoimpl.UnsafeEnabled {
		mi := &file_calculator_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Number) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Number) ProtoMessage() {}

func (x *Number) ProtoReflect() protoreflect.Message {
	mi := &file_calculator_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Number.ProtoReflect.Descriptor instead.
func (*Number) Descriptor() ([]byte, []int) {
	return file_calculator_proto_rawDescGZIP(), []int{1}
}

func (x *Number) GetValue() float64 {
	if x != nil {
		return x.Value
	}
	return 0
}

var File_calculator_proto protoreflect.FileDescriptor

var file_calculator_proto_rawDesc = []byte{
	0x0a, 0x10, 0x63, 0x61, 0x6c, 0x63, 0x75, 0x6c, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x02, 0x76, 0x31, 0x22, 0x21, 0x0a, 0x07, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72,
	0x73, 0x12, 0x16, 0x0a, 0x06, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28,
	0x01, 0x52, 0x06, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x22, 0x1e, 0x0a, 0x06, 0x4e, 0x75, 0x6d,
	0x62, 0x65, 0x72, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x01, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x32, 0xd2, 0x02, 0x0a, 0x0a, 0x43, 0x61,
	0x6c, 0x63, 0x75, 0x6c, 0x61, 0x74, 0x6f, 0x72, 0x12, 0x20, 0x0a, 0x03, 0x41, 0x64, 0x64, 0x12,
	0x0b, 0x2e, 0x76, 0x31, 0x2e, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x1a, 0x0a, 0x2e, 0x76,
	0x31, 0x2e, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x22, 0x00, 0x12, 0x27, 0x0a, 0x09, 0x41, 0x64,
	0x64, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x12, 0x0a, 0x2e, 0x76, 0x31, 0x2e, 0x4e, 0x75, 0x6d,
	0x62, 0x65, 0x72, 0x1a, 0x0a, 0x2e, 0x76, 0x31, 0x2e, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x22,
	0x00, 0x28, 0x01, 0x12, 0x25, 0x0a, 0x08, 0x53, 0x75, 0x62, 0x74, 0x72, 0x61, 0x63, 0x74, 0x12,
	0x0b, 0x2e, 0x76, 0x31, 0x2e, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x1a, 0x0a, 0x2e, 0x76,
	0x31, 0x2e, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x22, 0x00, 0x12, 0x2c, 0x0a, 0x0e, 0x53, 0x75,
	0x62, 0x74, 0x72, 0x61, 0x63, 0x74, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x12, 0x0a, 0x2e, 0x76,
	0x31, 0x2e, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x1a, 0x0a, 0x2e, 0x76, 0x31, 0x2e, 0x4e, 0x75,
	0x6d, 0x62, 0x65, 0x72, 0x22, 0x00, 0x28, 0x01, 0x12, 0x25, 0x0a, 0x08, 0x4d, 0x75, 0x6c, 0x74,
	0x69, 0x70, 0x6c, 0x79, 0x12, 0x0b, 0x2e, 0x76, 0x31, 0x2e, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72,
	0x73, 0x1a, 0x0a, 0x2e, 0x76, 0x31, 0x2e, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x22, 0x00, 0x12,
	0x2c, 0x0a, 0x0e, 0x4d, 0x75, 0x6c, 0x74, 0x69, 0x70, 0x6c, 0x79, 0x53, 0x74, 0x72, 0x65, 0x61,
	0x6d, 0x12, 0x0a, 0x2e, 0x76, 0x31, 0x2e, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x1a, 0x0a, 0x2e,
	0x76, 0x31, 0x2e, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x22, 0x00, 0x28, 0x01, 0x12, 0x23, 0x0a,
	0x06, 0x44, 0x69, 0x76, 0x69, 0x64, 0x65, 0x12, 0x0b, 0x2e, 0x76, 0x31, 0x2e, 0x4e, 0x75, 0x6d,
	0x62, 0x65, 0x72, 0x73, 0x1a, 0x0a, 0x2e, 0x76, 0x31, 0x2e, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72,
	0x22, 0x00, 0x12, 0x2a, 0x0a, 0x0c, 0x44, 0x69, 0x76, 0x69, 0x64, 0x65, 0x53, 0x74, 0x72, 0x65,
	0x61, 0x6d, 0x12, 0x0a, 0x2e, 0x76, 0x31, 0x2e, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x1a, 0x0a,
	0x2e, 0x76, 0x31, 0x2e, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x22, 0x00, 0x28, 0x01, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_calculator_proto_rawDescOnce sync.Once
	file_calculator_proto_rawDescData = file_calculator_proto_rawDesc
)

func file_calculator_proto_rawDescGZIP() []byte {
	file_calculator_proto_rawDescOnce.Do(func() {
		file_calculator_proto_rawDescData = protoimpl.X.CompressGZIP(file_calculator_proto_rawDescData)
	})
	return file_calculator_proto_rawDescData
}

var file_calculator_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_calculator_proto_goTypes = []interface{}{
	(*Numbers)(nil), // 0: v1.Numbers
	(*Number)(nil),  // 1: v1.Number
}
var file_calculator_proto_depIdxs = []int32{
	0, // 0: v1.Calculator.Add:input_type -> v1.Numbers
	1, // 1: v1.Calculator.AddStream:input_type -> v1.Number
	0, // 2: v1.Calculator.Subtract:input_type -> v1.Numbers
	1, // 3: v1.Calculator.SubtractStream:input_type -> v1.Number
	0, // 4: v1.Calculator.Multiply:input_type -> v1.Numbers
	1, // 5: v1.Calculator.MultiplyStream:input_type -> v1.Number
	0, // 6: v1.Calculator.Divide:input_type -> v1.Numbers
	1, // 7: v1.Calculator.DivideStream:input_type -> v1.Number
	1, // 8: v1.Calculator.Add:output_type -> v1.Number
	1, // 9: v1.Calculator.AddStream:output_type -> v1.Number
	1, // 10: v1.Calculator.Subtract:output_type -> v1.Number
	1, // 11: v1.Calculator.SubtractStream:output_type -> v1.Number
	1, // 12: v1.Calculator.Multiply:output_type -> v1.Number
	1, // 13: v1.Calculator.MultiplyStream:output_type -> v1.Number
	1, // 14: v1.Calculator.Divide:output_type -> v1.Number
	1, // 15: v1.Calculator.DivideStream:output_type -> v1.Number
	8, // [8:16] is the sub-list for method output_type
	0, // [0:8] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_calculator_proto_init() }
func file_calculator_proto_init() {
	if File_calculator_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_calculator_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Numbers); i {
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
		file_calculator_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Number); i {
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
			RawDescriptor: file_calculator_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_calculator_proto_goTypes,
		DependencyIndexes: file_calculator_proto_depIdxs,
		MessageInfos:      file_calculator_proto_msgTypes,
	}.Build()
	File_calculator_proto = out.File
	file_calculator_proto_rawDesc = nil
	file_calculator_proto_goTypes = nil
	file_calculator_proto_depIdxs = nil
}

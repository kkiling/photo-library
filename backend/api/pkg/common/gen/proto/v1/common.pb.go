// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        (unknown)
// source: proto/v1/common.proto

package pbv1

import (
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2/options"
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

type Paginator struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Offset uint32 `protobuf:"varint,1,opt,name=offset,proto3" json:"offset,omitempty"`
	Limit  uint32 `protobuf:"varint,2,opt,name=limit,proto3" json:"limit,omitempty"`
}

func (x *Paginator) Reset() {
	*x = Paginator{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_v1_common_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Paginator) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Paginator) ProtoMessage() {}

func (x *Paginator) ProtoReflect() protoreflect.Message {
	mi := &file_proto_v1_common_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Paginator.ProtoReflect.Descriptor instead.
func (*Paginator) Descriptor() ([]byte, []int) {
	return file_proto_v1_common_proto_rawDescGZIP(), []int{0}
}

func (x *Paginator) GetOffset() uint32 {
	if x != nil {
		return x.Offset
	}
	return 0
}

func (x *Paginator) GetLimit() uint32 {
	if x != nil {
		return x.Limit
	}
	return 0
}

type ServerError struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code    uint32       `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Message string       `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	Details []*ErrorInfo `protobuf:"bytes,3,rep,name=details,proto3" json:"details,omitempty"`
}

func (x *ServerError) Reset() {
	*x = ServerError{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_v1_common_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ServerError) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ServerError) ProtoMessage() {}

func (x *ServerError) ProtoReflect() protoreflect.Message {
	mi := &file_proto_v1_common_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ServerError.ProtoReflect.Descriptor instead.
func (*ServerError) Descriptor() ([]byte, []int) {
	return file_proto_v1_common_proto_rawDescGZIP(), []int{1}
}

func (x *ServerError) GetCode() uint32 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *ServerError) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *ServerError) GetDetails() []*ErrorInfo {
	if x != nil {
		return x.Details
	}
	return nil
}

type ErrorInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Description     string            `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	FieldViolations []*FieldViolation `protobuf:"bytes,4,rep,name=field_violations,json=fieldViolations,proto3" json:"field_violations,omitempty"`
}

func (x *ErrorInfo) Reset() {
	*x = ErrorInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_v1_common_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ErrorInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ErrorInfo) ProtoMessage() {}

func (x *ErrorInfo) ProtoReflect() protoreflect.Message {
	mi := &file_proto_v1_common_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ErrorInfo.ProtoReflect.Descriptor instead.
func (*ErrorInfo) Descriptor() ([]byte, []int) {
	return file_proto_v1_common_proto_rawDescGZIP(), []int{2}
}

func (x *ErrorInfo) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *ErrorInfo) GetFieldViolations() []*FieldViolation {
	if x != nil {
		return x.FieldViolations
	}
	return nil
}

type FieldViolation struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Field string `protobuf:"bytes,1,opt,name=field,proto3" json:"field,omitempty"`
	Error string `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty"`
}

func (x *FieldViolation) Reset() {
	*x = FieldViolation{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_v1_common_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FieldViolation) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FieldViolation) ProtoMessage() {}

func (x *FieldViolation) ProtoReflect() protoreflect.Message {
	mi := &file_proto_v1_common_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FieldViolation.ProtoReflect.Descriptor instead.
func (*FieldViolation) Descriptor() ([]byte, []int) {
	return file_proto_v1_common_proto_rawDescGZIP(), []int{3}
}

func (x *FieldViolation) GetField() string {
	if x != nil {
		return x.Field
	}
	return ""
}

func (x *FieldViolation) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

var File_proto_v1_common_proto protoreflect.FileDescriptor

var file_proto_v1_common_proto_rawDesc = []byte{
	0x0a, 0x15, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f,
	0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70, 0x62, 0x2e, 0x76, 0x31, 0x1a, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x2d, 0x67, 0x65, 0x6e, 0x2d, 0x6f, 0x70, 0x65, 0x6e, 0x61,
	0x70, 0x69, 0x76, 0x32, 0x2f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x61, 0x6e, 0x6e,
	0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x39,
	0x0a, 0x09, 0x50, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x6f, 0x72, 0x12, 0x16, 0x0a, 0x06, 0x6f,
	0x66, 0x66, 0x73, 0x65, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x6f, 0x66, 0x66,
	0x73, 0x65, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0d, 0x52, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x22, 0x92, 0x01, 0x0a, 0x0b, 0x53, 0x65,
	0x72, 0x76, 0x65, 0x72, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x18, 0x0a,
	0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x2a, 0x0a, 0x07, 0x64, 0x65, 0x74, 0x61, 0x69,
	0x6c, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x70, 0x62, 0x2e, 0x76, 0x31,
	0x2e, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x07, 0x64, 0x65, 0x74, 0x61,
	0x69, 0x6c, 0x73, 0x3a, 0x29, 0x92, 0x41, 0x26, 0x0a, 0x24, 0x2a, 0x0b, 0x53, 0x65, 0x72, 0x76,
	0x65, 0x72, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x32, 0x15, 0x43, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x20,
	0x65, 0x72, 0x72, 0x6f, 0x72, 0x20, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x6f,
	0x0a, 0x09, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x20, 0x0a, 0x0b, 0x64,
	0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x40, 0x0a,
	0x10, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x5f, 0x76, 0x69, 0x6f, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x70, 0x62, 0x2e, 0x76, 0x31, 0x2e,
	0x46, 0x69, 0x65, 0x6c, 0x64, 0x56, 0x69, 0x6f, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0f,
	0x66, 0x69, 0x65, 0x6c, 0x64, 0x56, 0x69, 0x6f, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x22,
	0x3c, 0x0a, 0x0e, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x56, 0x69, 0x6f, 0x6c, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x12, 0x14, 0x0a, 0x05, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x42, 0x51, 0x5a,
	0x4f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6b, 0x6b, 0x69, 0x6c,
	0x69, 0x6e, 0x67, 0x2f, 0x70, 0x68, 0x6f, 0x74, 0x6f, 0x2d, 0x6c, 0x69, 0x62, 0x72, 0x61, 0x72,
	0x79, 0x2f, 0x62, 0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64, 0x2f, 0x73, 0x79, 0x6e, 0x63, 0x5f, 0x70,
	0x68, 0x6f, 0x74, 0x6f, 0x73, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2f, 0x70, 0x6b, 0x67,
	0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x70, 0x62, 0x76, 0x31,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_v1_common_proto_rawDescOnce sync.Once
	file_proto_v1_common_proto_rawDescData = file_proto_v1_common_proto_rawDesc
)

func file_proto_v1_common_proto_rawDescGZIP() []byte {
	file_proto_v1_common_proto_rawDescOnce.Do(func() {
		file_proto_v1_common_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_v1_common_proto_rawDescData)
	})
	return file_proto_v1_common_proto_rawDescData
}

var file_proto_v1_common_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_proto_v1_common_proto_goTypes = []interface{}{
	(*Paginator)(nil),      // 0: pb.v1.Paginator
	(*ServerError)(nil),    // 1: pb.v1.ServerError
	(*ErrorInfo)(nil),      // 2: pb.v1.ErrorInfo
	(*FieldViolation)(nil), // 3: pb.v1.FieldViolation
}
var file_proto_v1_common_proto_depIdxs = []int32{
	2, // 0: pb.v1.ServerError.details:type_name -> pb.v1.ErrorInfo
	3, // 1: pb.v1.ErrorInfo.field_violations:type_name -> pb.v1.FieldViolation
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_proto_v1_common_proto_init() }
func file_proto_v1_common_proto_init() {
	if File_proto_v1_common_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_v1_common_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Paginator); i {
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
		file_proto_v1_common_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ServerError); i {
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
		file_proto_v1_common_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ErrorInfo); i {
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
		file_proto_v1_common_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FieldViolation); i {
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
			RawDescriptor: file_proto_v1_common_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_proto_v1_common_proto_goTypes,
		DependencyIndexes: file_proto_v1_common_proto_depIdxs,
		MessageInfos:      file_proto_v1_common_proto_msgTypes,
	}.Build()
	File_proto_v1_common_proto = out.File
	file_proto_v1_common_proto_rawDesc = nil
	file_proto_v1_common_proto_goTypes = nil
	file_proto_v1_common_proto_depIdxs = nil
}

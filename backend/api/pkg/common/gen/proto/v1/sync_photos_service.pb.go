// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.32.0
// 	protoc        (unknown)
// source: proto/v1/sync_photos_service.proto

package pbv1

import (
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2/options"
	_ "google.golang.org/genproto/googleapis/api/annotations"
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

type UploadPhotoRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Пути фотографий которые загружаем (может быть несколько если фото одинаковые)
	Paths []string `protobuf:"bytes,1,rep,name=paths,proto3" json:"paths,omitempty"`
	// Рассчитанный на клиенте хеш фотографии
	Hash string `protobuf:"bytes,2,opt,name=hash,proto3" json:"hash,omitempty"`
	// Данные фото
	Body []byte `protobuf:"bytes,3,opt,name=body,proto3" json:"body,omitempty"`
	// Информация о последнем изменении фото
	UpdateAt *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=update_at,json=updateAt,proto3" json:"update_at,omitempty"`
	// Идентификатор клиента
	ClientId string `protobuf:"bytes,5,opt,name=client_id,json=clientId,proto3" json:"client_id,omitempty"`
	// Ключ доступа
	AccessKey string `protobuf:"bytes,6,opt,name=access_key,json=accessKey,proto3" json:"access_key,omitempty"`
}

func (x *UploadPhotoRequest) Reset() {
	*x = UploadPhotoRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_v1_sync_photos_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UploadPhotoRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UploadPhotoRequest) ProtoMessage() {}

func (x *UploadPhotoRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_v1_sync_photos_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UploadPhotoRequest.ProtoReflect.Descriptor instead.
func (*UploadPhotoRequest) Descriptor() ([]byte, []int) {
	return file_proto_v1_sync_photos_service_proto_rawDescGZIP(), []int{0}
}

func (x *UploadPhotoRequest) GetPaths() []string {
	if x != nil {
		return x.Paths
	}
	return nil
}

func (x *UploadPhotoRequest) GetHash() string {
	if x != nil {
		return x.Hash
	}
	return ""
}

func (x *UploadPhotoRequest) GetBody() []byte {
	if x != nil {
		return x.Body
	}
	return nil
}

func (x *UploadPhotoRequest) GetUpdateAt() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdateAt
	}
	return nil
}

func (x *UploadPhotoRequest) GetClientId() string {
	if x != nil {
		return x.ClientId
	}
	return ""
}

func (x *UploadPhotoRequest) GetAccessKey() string {
	if x != nil {
		return x.AccessKey
	}
	return ""
}

type UploadPhotoResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Фото было загружено ранее
	HasBeenUploadedBefore bool `protobuf:"varint,1,opt,name=has_been_uploaded_before,json=hasBeenUploadedBefore,proto3" json:"has_been_uploaded_before,omitempty"`
	// Хеш фотографии
	Hash string `protobuf:"bytes,3,opt,name=hash,proto3" json:"hash,omitempty"`
}

func (x *UploadPhotoResponse) Reset() {
	*x = UploadPhotoResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_v1_sync_photos_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UploadPhotoResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UploadPhotoResponse) ProtoMessage() {}

func (x *UploadPhotoResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_v1_sync_photos_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UploadPhotoResponse.ProtoReflect.Descriptor instead.
func (*UploadPhotoResponse) Descriptor() ([]byte, []int) {
	return file_proto_v1_sync_photos_service_proto_rawDescGZIP(), []int{1}
}

func (x *UploadPhotoResponse) GetHasBeenUploadedBefore() bool {
	if x != nil {
		return x.HasBeenUploadedBefore
	}
	return false
}

func (x *UploadPhotoResponse) GetHash() string {
	if x != nil {
		return x.Hash
	}
	return ""
}

var File_proto_v1_sync_photos_service_proto protoreflect.FileDescriptor

var file_proto_v1_sync_photos_service_proto_rawDesc = []byte{
	0x0a, 0x22, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x76, 0x31, 0x2f, 0x73, 0x79, 0x6e, 0x63, 0x5f,
	0x70, 0x68, 0x6f, 0x74, 0x6f, 0x73, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70, 0x62, 0x2e, 0x76, 0x31, 0x1a, 0x1c, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x63, 0x2d, 0x67, 0x65, 0x6e, 0x2d, 0x6f, 0x70, 0x65, 0x6e, 0x61, 0x70, 0x69, 0x76, 0x32, 0x2f,
	0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xc7, 0x01, 0x0a, 0x12, 0x55,
	0x70, 0x6c, 0x6f, 0x61, 0x64, 0x50, 0x68, 0x6f, 0x74, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x14, 0x0a, 0x05, 0x70, 0x61, 0x74, 0x68, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09,
	0x52, 0x05, 0x70, 0x61, 0x74, 0x68, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x68, 0x61, 0x73, 0x68, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x68, 0x61, 0x73, 0x68, 0x12, 0x12, 0x0a, 0x04, 0x62,
	0x6f, 0x64, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x12,
	0x37, 0x0a, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x5f, 0x61, 0x74, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x08,
	0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x41, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x63, 0x6c, 0x69, 0x65,
	0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x63, 0x6c, 0x69,
	0x65, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x5f,
	0x6b, 0x65, 0x79, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x61, 0x63, 0x63, 0x65, 0x73,
	0x73, 0x4b, 0x65, 0x79, 0x22, 0x62, 0x0a, 0x13, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x50, 0x68,
	0x6f, 0x74, 0x6f, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x37, 0x0a, 0x18, 0x68,
	0x61, 0x73, 0x5f, 0x62, 0x65, 0x65, 0x6e, 0x5f, 0x75, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x65, 0x64,
	0x5f, 0x62, 0x65, 0x66, 0x6f, 0x72, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x15, 0x68,
	0x61, 0x73, 0x42, 0x65, 0x65, 0x6e, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x65, 0x64, 0x42, 0x65,
	0x66, 0x6f, 0x72, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x68, 0x61, 0x73, 0x68, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x68, 0x61, 0x73, 0x68, 0x32, 0xbe, 0x01, 0x0a, 0x11, 0x53, 0x79, 0x6e,
	0x63, 0x50, 0x68, 0x6f, 0x74, 0x6f, 0x73, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0xa8,
	0x01, 0x0a, 0x0b, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x50, 0x68, 0x6f, 0x74, 0x6f, 0x12, 0x19,
	0x2e, 0x70, 0x62, 0x2e, 0x76, 0x31, 0x2e, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x50, 0x68, 0x6f,
	0x74, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x70, 0x62, 0x2e, 0x76,
	0x31, 0x2e, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x50, 0x68, 0x6f, 0x74, 0x6f, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x62, 0x92, 0x41, 0x3e, 0x0a, 0x0a, 0x73, 0x79, 0x6e, 0x63,
	0x2d, 0x70, 0x68, 0x6f, 0x74, 0x6f, 0x12, 0x30, 0xd0, 0x97, 0xd0, 0xb0, 0xd0, 0xb3, 0xd1, 0x80,
	0xd1, 0x83, 0xd0, 0xb7, 0xd0, 0xba, 0xd0, 0xb0, 0x20, 0xd0, 0xbd, 0xd0, 0xbe, 0xd0, 0xb2, 0xd0,
	0xbe, 0xd0, 0xb9, 0x20, 0xd1, 0x84, 0xd0, 0xbe, 0xd1, 0x82, 0xd0, 0xbe, 0xd0, 0xb3, 0xd1, 0x80,
	0xd0, 0xb0, 0xd1, 0x84, 0xd0, 0xb8, 0xd0, 0xb8, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x1b, 0x3a, 0x01,
	0x2a, 0x22, 0x16, 0x2f, 0x76, 0x31, 0x2f, 0x73, 0x79, 0x6e, 0x63, 0x2d, 0x70, 0x68, 0x6f, 0x74,
	0x6f, 0x73, 0x2f, 0x75, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x42, 0xa3, 0x01, 0x92, 0x41, 0x4f, 0x12,
	0x16, 0x0a, 0x0f, 0x53, 0x79, 0x6e, 0x63, 0x20, 0x70, 0x68, 0x6f, 0x74, 0x6f, 0x73, 0x20, 0x41,
	0x50, 0x49, 0x32, 0x03, 0x30, 0x2e, 0x31, 0x1a, 0x0e, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x68, 0x6f,
	0x73, 0x74, 0x3a, 0x38, 0x30, 0x38, 0x30, 0x2a, 0x01, 0x01, 0x32, 0x10, 0x61, 0x70, 0x70, 0x6c,
	0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x6a, 0x73, 0x6f, 0x6e, 0x3a, 0x10, 0x61, 0x70,
	0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x6a, 0x73, 0x6f, 0x6e, 0x5a, 0x4f,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6b, 0x6b, 0x69, 0x6c, 0x69,
	0x6e, 0x67, 0x2f, 0x70, 0x68, 0x6f, 0x74, 0x6f, 0x2d, 0x6c, 0x69, 0x62, 0x72, 0x61, 0x72, 0x79,
	0x2f, 0x62, 0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64, 0x2f, 0x73, 0x79, 0x6e, 0x63, 0x5f, 0x70, 0x68,
	0x6f, 0x74, 0x6f, 0x73, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2f, 0x70, 0x6b, 0x67, 0x2f,
	0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x70, 0x62, 0x76, 0x31, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_v1_sync_photos_service_proto_rawDescOnce sync.Once
	file_proto_v1_sync_photos_service_proto_rawDescData = file_proto_v1_sync_photos_service_proto_rawDesc
)

func file_proto_v1_sync_photos_service_proto_rawDescGZIP() []byte {
	file_proto_v1_sync_photos_service_proto_rawDescOnce.Do(func() {
		file_proto_v1_sync_photos_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_v1_sync_photos_service_proto_rawDescData)
	})
	return file_proto_v1_sync_photos_service_proto_rawDescData
}

var file_proto_v1_sync_photos_service_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_proto_v1_sync_photos_service_proto_goTypes = []interface{}{
	(*UploadPhotoRequest)(nil),    // 0: pb.v1.UploadPhotoRequest
	(*UploadPhotoResponse)(nil),   // 1: pb.v1.UploadPhotoResponse
	(*timestamppb.Timestamp)(nil), // 2: google.protobuf.Timestamp
}
var file_proto_v1_sync_photos_service_proto_depIdxs = []int32{
	2, // 0: pb.v1.UploadPhotoRequest.update_at:type_name -> google.protobuf.Timestamp
	0, // 1: pb.v1.SyncPhotosService.UploadPhoto:input_type -> pb.v1.UploadPhotoRequest
	1, // 2: pb.v1.SyncPhotosService.UploadPhoto:output_type -> pb.v1.UploadPhotoResponse
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_proto_v1_sync_photos_service_proto_init() }
func file_proto_v1_sync_photos_service_proto_init() {
	if File_proto_v1_sync_photos_service_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_v1_sync_photos_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UploadPhotoRequest); i {
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
		file_proto_v1_sync_photos_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UploadPhotoResponse); i {
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
			RawDescriptor: file_proto_v1_sync_photos_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_v1_sync_photos_service_proto_goTypes,
		DependencyIndexes: file_proto_v1_sync_photos_service_proto_depIdxs,
		MessageInfos:      file_proto_v1_sync_photos_service_proto_msgTypes,
	}.Build()
	File_proto_v1_sync_photos_service_proto = out.File
	file_proto_v1_sync_photos_service_proto_rawDesc = nil
	file_proto_v1_sync_photos_service_proto_goTypes = nil
	file_proto_v1_sync_photos_service_proto_depIdxs = nil
}
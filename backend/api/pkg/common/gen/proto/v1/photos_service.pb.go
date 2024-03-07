// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.32.0
// 	protoc        (unknown)
// source: proto/v1/photos_service.proto

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

type GetPhotoGroupsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Page    int32 `protobuf:"varint,1,opt,name=page,proto3" json:"page,omitempty"`
	PerPage int32 `protobuf:"varint,2,opt,name=per_page,json=perPage,proto3" json:"per_page,omitempty"`
}

func (x *GetPhotoGroupsRequest) Reset() {
	*x = GetPhotoGroupsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_v1_photos_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetPhotoGroupsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetPhotoGroupsRequest) ProtoMessage() {}

func (x *GetPhotoGroupsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_v1_photos_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetPhotoGroupsRequest.ProtoReflect.Descriptor instead.
func (*GetPhotoGroupsRequest) Descriptor() ([]byte, []int) {
	return file_proto_v1_photos_service_proto_rawDescGZIP(), []int{0}
}

func (x *GetPhotoGroupsRequest) GetPage() int32 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *GetPhotoGroupsRequest) GetPerPage() int32 {
	if x != nil {
		return x.PerPage
	}
	return 0
}

type GetPhotoGroupsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Items      []*PhotoGroup `protobuf:"bytes,1,rep,name=items,proto3" json:"items,omitempty"`
	TotalItems int32         `protobuf:"varint,2,opt,name=total_items,json=totalItems,proto3" json:"total_items,omitempty"`
}

func (x *GetPhotoGroupsResponse) Reset() {
	*x = GetPhotoGroupsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_v1_photos_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetPhotoGroupsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetPhotoGroupsResponse) ProtoMessage() {}

func (x *GetPhotoGroupsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_v1_photos_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetPhotoGroupsResponse.ProtoReflect.Descriptor instead.
func (*GetPhotoGroupsResponse) Descriptor() ([]byte, []int) {
	return file_proto_v1_photos_service_proto_rawDescGZIP(), []int{1}
}

func (x *GetPhotoGroupsResponse) GetItems() []*PhotoGroup {
	if x != nil {
		return x.Items
	}
	return nil
}

func (x *GetPhotoGroupsResponse) GetTotalItems() int32 {
	if x != nil {
		return x.TotalItems
	}
	return 0
}

type GetPhotoGroupRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	GroupId string `protobuf:"bytes,1,opt,name=group_id,json=groupId,proto3" json:"group_id,omitempty"`
}

func (x *GetPhotoGroupRequest) Reset() {
	*x = GetPhotoGroupRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_v1_photos_service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetPhotoGroupRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetPhotoGroupRequest) ProtoMessage() {}

func (x *GetPhotoGroupRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_v1_photos_service_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetPhotoGroupRequest.ProtoReflect.Descriptor instead.
func (*GetPhotoGroupRequest) Descriptor() ([]byte, []int) {
	return file_proto_v1_photos_service_proto_rawDescGZIP(), []int{2}
}

func (x *GetPhotoGroupRequest) GetGroupId() string {
	if x != nil {
		return x.GroupId
	}
	return ""
}

type GetPhotoGroupResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        string           `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	MainPhoto *PhotoWithData   `protobuf:"bytes,2,opt,name=main_photo,json=mainPhoto,proto3" json:"main_photo,omitempty"`
	Photos    []*PhotoWithData `protobuf:"bytes,3,rep,name=photos,proto3" json:"photos,omitempty"`
}

func (x *GetPhotoGroupResponse) Reset() {
	*x = GetPhotoGroupResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_v1_photos_service_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetPhotoGroupResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetPhotoGroupResponse) ProtoMessage() {}

func (x *GetPhotoGroupResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_v1_photos_service_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetPhotoGroupResponse.ProtoReflect.Descriptor instead.
func (*GetPhotoGroupResponse) Descriptor() ([]byte, []int) {
	return file_proto_v1_photos_service_proto_rawDescGZIP(), []int{3}
}

func (x *GetPhotoGroupResponse) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *GetPhotoGroupResponse) GetMainPhoto() *PhotoWithData {
	if x != nil {
		return x.MainPhoto
	}
	return nil
}

func (x *GetPhotoGroupResponse) GetPhotos() []*PhotoWithData {
	if x != nil {
		return x.Photos
	}
	return nil
}

type PhotoWithData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id       string    `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Url      string    `protobuf:"bytes,2,opt,name=url,proto3" json:"url,omitempty"`
	MetaData *Metadata `protobuf:"bytes,3,opt,name=meta_data,json=metaData,proto3" json:"meta_data,omitempty"`
	Tags     []*Tag    `protobuf:"bytes,4,rep,name=tags,proto3" json:"tags,omitempty"`
}

func (x *PhotoWithData) Reset() {
	*x = PhotoWithData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_v1_photos_service_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PhotoWithData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PhotoWithData) ProtoMessage() {}

func (x *PhotoWithData) ProtoReflect() protoreflect.Message {
	mi := &file_proto_v1_photos_service_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PhotoWithData.ProtoReflect.Descriptor instead.
func (*PhotoWithData) Descriptor() ([]byte, []int) {
	return file_proto_v1_photos_service_proto_rawDescGZIP(), []int{4}
}

func (x *PhotoWithData) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *PhotoWithData) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *PhotoWithData) GetMetaData() *Metadata {
	if x != nil {
		return x.MetaData
	}
	return nil
}

func (x *PhotoWithData) GetTags() []*Tag {
	if x != nil {
		return x.Tags
	}
	return nil
}

type Geo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Latitude  float64 `protobuf:"fixed64,1,opt,name=Latitude,proto3" json:"Latitude,omitempty"`
	Longitude float64 `protobuf:"fixed64,2,opt,name=Longitude,proto3" json:"Longitude,omitempty"`
}

func (x *Geo) Reset() {
	*x = Geo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_v1_photos_service_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Geo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Geo) ProtoMessage() {}

func (x *Geo) ProtoReflect() protoreflect.Message {
	mi := &file_proto_v1_photos_service_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Geo.ProtoReflect.Descriptor instead.
func (*Geo) Descriptor() ([]byte, []int) {
	return file_proto_v1_photos_service_proto_rawDescGZIP(), []int{5}
}

func (x *Geo) GetLatitude() float64 {
	if x != nil {
		return x.Latitude
	}
	return 0
}

func (x *Geo) GetLongitude() float64 {
	if x != nil {
		return x.Longitude
	}
	return 0
}

type Metadata struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ModelInfo   *string                `protobuf:"bytes,1,opt,name=model_info,json=modelInfo,proto3,oneof" json:"model_info,omitempty"`
	SizeBytes   int32                  `protobuf:"varint,2,opt,name=size_bytes,json=sizeBytes,proto3" json:"size_bytes,omitempty"`
	WidthPixel  int32                  `protobuf:"varint,3,opt,name=width_pixel,json=widthPixel,proto3" json:"width_pixel,omitempty"`
	HeightPixel int32                  `protobuf:"varint,4,opt,name=height_pixel,json=heightPixel,proto3" json:"height_pixel,omitempty"`
	DataTime    *timestamppb.Timestamp `protobuf:"bytes,5,opt,name=data_time,json=dataTime,proto3" json:"data_time,omitempty"`
	UpdateAt    *timestamppb.Timestamp `protobuf:"bytes,6,opt,name=update_at,json=updateAt,proto3" json:"update_at,omitempty"`
	Geo         *Geo                   `protobuf:"bytes,7,opt,name=geo,proto3" json:"geo,omitempty"`
}

func (x *Metadata) Reset() {
	*x = Metadata{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_v1_photos_service_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Metadata) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Metadata) ProtoMessage() {}

func (x *Metadata) ProtoReflect() protoreflect.Message {
	mi := &file_proto_v1_photos_service_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Metadata.ProtoReflect.Descriptor instead.
func (*Metadata) Descriptor() ([]byte, []int) {
	return file_proto_v1_photos_service_proto_rawDescGZIP(), []int{6}
}

func (x *Metadata) GetModelInfo() string {
	if x != nil && x.ModelInfo != nil {
		return *x.ModelInfo
	}
	return ""
}

func (x *Metadata) GetSizeBytes() int32 {
	if x != nil {
		return x.SizeBytes
	}
	return 0
}

func (x *Metadata) GetWidthPixel() int32 {
	if x != nil {
		return x.WidthPixel
	}
	return 0
}

func (x *Metadata) GetHeightPixel() int32 {
	if x != nil {
		return x.HeightPixel
	}
	return 0
}

func (x *Metadata) GetDataTime() *timestamppb.Timestamp {
	if x != nil {
		return x.DataTime
	}
	return nil
}

func (x *Metadata) GetUpdateAt() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdateAt
	}
	return nil
}

func (x *Metadata) GetGeo() *Geo {
	if x != nil {
		return x.Geo
	}
	return nil
}

type Tag struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id    string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name  string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Type  string `protobuf:"bytes,3,opt,name=type,proto3" json:"type,omitempty"`
	Color string `protobuf:"bytes,4,opt,name=color,proto3" json:"color,omitempty"`
}

func (x *Tag) Reset() {
	*x = Tag{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_v1_photos_service_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Tag) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Tag) ProtoMessage() {}

func (x *Tag) ProtoReflect() protoreflect.Message {
	mi := &file_proto_v1_photos_service_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Tag.ProtoReflect.Descriptor instead.
func (*Tag) Descriptor() ([]byte, []int) {
	return file_proto_v1_photos_service_proto_rawDescGZIP(), []int{7}
}

func (x *Tag) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Tag) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Tag) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *Tag) GetColor() string {
	if x != nil {
		return x.Color
	}
	return ""
}

type Photo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id  string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Url string `protobuf:"bytes,2,opt,name=url,proto3" json:"url,omitempty"`
}

func (x *Photo) Reset() {
	*x = Photo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_v1_photos_service_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Photo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Photo) ProtoMessage() {}

func (x *Photo) ProtoReflect() protoreflect.Message {
	mi := &file_proto_v1_photos_service_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Photo.ProtoReflect.Descriptor instead.
func (*Photo) Descriptor() ([]byte, []int) {
	return file_proto_v1_photos_service_proto_rawDescGZIP(), []int{8}
}

func (x *Photo) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Photo) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

type PhotoGroup struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id          string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	MainPhoto   *Photo `protobuf:"bytes,2,opt,name=main_photo,json=mainPhoto,proto3" json:"main_photo,omitempty"`
	PhotosCount int32  `protobuf:"varint,3,opt,name=photos_count,json=photosCount,proto3" json:"photos_count,omitempty"`
}

func (x *PhotoGroup) Reset() {
	*x = PhotoGroup{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_v1_photos_service_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PhotoGroup) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PhotoGroup) ProtoMessage() {}

func (x *PhotoGroup) ProtoReflect() protoreflect.Message {
	mi := &file_proto_v1_photos_service_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PhotoGroup.ProtoReflect.Descriptor instead.
func (*PhotoGroup) Descriptor() ([]byte, []int) {
	return file_proto_v1_photos_service_proto_rawDescGZIP(), []int{9}
}

func (x *PhotoGroup) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *PhotoGroup) GetMainPhoto() *Photo {
	if x != nil {
		return x.MainPhoto
	}
	return nil
}

func (x *PhotoGroup) GetPhotosCount() int32 {
	if x != nil {
		return x.PhotosCount
	}
	return 0
}

var File_proto_v1_photos_service_proto protoreflect.FileDescriptor

var file_proto_v1_photos_service_proto_rawDesc = []byte{
	0x0a, 0x1d, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x76, 0x31, 0x2f, 0x70, 0x68, 0x6f, 0x74, 0x6f,
	0x73, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x05, 0x70, 0x62, 0x2e, 0x76, 0x31, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61,
	0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x2d, 0x67, 0x65, 0x6e,
	0x2d, 0x6f, 0x70, 0x65, 0x6e, 0x61, 0x70, 0x69, 0x76, 0x32, 0x2f, 0x6f, 0x70, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x46, 0x0a, 0x15, 0x47, 0x65, 0x74, 0x50, 0x68, 0x6f, 0x74,
	0x6f, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12,
	0x0a, 0x04, 0x70, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x70, 0x61,
	0x67, 0x65, 0x12, 0x19, 0x0a, 0x08, 0x70, 0x65, 0x72, 0x5f, 0x70, 0x61, 0x67, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x70, 0x65, 0x72, 0x50, 0x61, 0x67, 0x65, 0x22, 0x62, 0x0a,
	0x16, 0x47, 0x65, 0x74, 0x50, 0x68, 0x6f, 0x74, 0x6f, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x73, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x27, 0x0a, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x70, 0x62, 0x2e, 0x76, 0x31, 0x2e, 0x50,
	0x68, 0x6f, 0x74, 0x6f, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x52, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73,
	0x12, 0x1f, 0x0a, 0x0b, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x5f, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x49, 0x74, 0x65, 0x6d,
	0x73, 0x22, 0x31, 0x0a, 0x14, 0x47, 0x65, 0x74, 0x50, 0x68, 0x6f, 0x74, 0x6f, 0x47, 0x72, 0x6f,
	0x75, 0x70, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x19, 0x0a, 0x08, 0x67, 0x72, 0x6f,
	0x75, 0x70, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x67, 0x72, 0x6f,
	0x75, 0x70, 0x49, 0x64, 0x22, 0x8a, 0x01, 0x0a, 0x15, 0x47, 0x65, 0x74, 0x50, 0x68, 0x6f, 0x74,
	0x6f, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x0e,
	0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x33,
	0x0a, 0x0a, 0x6d, 0x61, 0x69, 0x6e, 0x5f, 0x70, 0x68, 0x6f, 0x74, 0x6f, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x14, 0x2e, 0x70, 0x62, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x68, 0x6f, 0x74, 0x6f,
	0x57, 0x69, 0x74, 0x68, 0x44, 0x61, 0x74, 0x61, 0x52, 0x09, 0x6d, 0x61, 0x69, 0x6e, 0x50, 0x68,
	0x6f, 0x74, 0x6f, 0x12, 0x2c, 0x0a, 0x06, 0x70, 0x68, 0x6f, 0x74, 0x6f, 0x73, 0x18, 0x03, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x70, 0x62, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x68, 0x6f, 0x74,
	0x6f, 0x57, 0x69, 0x74, 0x68, 0x44, 0x61, 0x74, 0x61, 0x52, 0x06, 0x70, 0x68, 0x6f, 0x74, 0x6f,
	0x73, 0x22, 0x7f, 0x0a, 0x0d, 0x50, 0x68, 0x6f, 0x74, 0x6f, 0x57, 0x69, 0x74, 0x68, 0x44, 0x61,
	0x74, 0x61, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02,
	0x69, 0x64, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x03, 0x75, 0x72, 0x6c, 0x12, 0x2c, 0x0a, 0x09, 0x6d, 0x65, 0x74, 0x61, 0x5f, 0x64, 0x61, 0x74,
	0x61, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x70, 0x62, 0x2e, 0x76, 0x31, 0x2e,
	0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x52, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x44, 0x61,
	0x74, 0x61, 0x12, 0x1e, 0x0a, 0x04, 0x74, 0x61, 0x67, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x0a, 0x2e, 0x70, 0x62, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x61, 0x67, 0x52, 0x04, 0x74, 0x61,
	0x67, 0x73, 0x22, 0x3f, 0x0a, 0x03, 0x47, 0x65, 0x6f, 0x12, 0x1a, 0x0a, 0x08, 0x4c, 0x61, 0x74,
	0x69, 0x74, 0x75, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x01, 0x52, 0x08, 0x4c, 0x61, 0x74,
	0x69, 0x74, 0x75, 0x64, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x4c, 0x6f, 0x6e, 0x67, 0x69, 0x74, 0x75,
	0x64, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01, 0x52, 0x09, 0x4c, 0x6f, 0x6e, 0x67, 0x69, 0x74,
	0x75, 0x64, 0x65, 0x22, 0xb0, 0x02, 0x0a, 0x08, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61,
	0x12, 0x22, 0x0a, 0x0a, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x09, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x49, 0x6e, 0x66,
	0x6f, 0x88, 0x01, 0x01, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x69, 0x7a, 0x65, 0x5f, 0x62, 0x79, 0x74,
	0x65, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x73, 0x69, 0x7a, 0x65, 0x42, 0x79,
	0x74, 0x65, 0x73, 0x12, 0x1f, 0x0a, 0x0b, 0x77, 0x69, 0x64, 0x74, 0x68, 0x5f, 0x70, 0x69, 0x78,
	0x65, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x77, 0x69, 0x64, 0x74, 0x68, 0x50,
	0x69, 0x78, 0x65, 0x6c, 0x12, 0x21, 0x0a, 0x0c, 0x68, 0x65, 0x69, 0x67, 0x68, 0x74, 0x5f, 0x70,
	0x69, 0x78, 0x65, 0x6c, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0b, 0x68, 0x65, 0x69, 0x67,
	0x68, 0x74, 0x50, 0x69, 0x78, 0x65, 0x6c, 0x12, 0x37, 0x0a, 0x09, 0x64, 0x61, 0x74, 0x61, 0x5f,
	0x74, 0x69, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x08, 0x64, 0x61, 0x74, 0x61, 0x54, 0x69, 0x6d, 0x65,
	0x12, 0x37, 0x0a, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x5f, 0x61, 0x74, 0x18, 0x06, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52,
	0x08, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x41, 0x74, 0x12, 0x1c, 0x0a, 0x03, 0x67, 0x65, 0x6f,
	0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0a, 0x2e, 0x70, 0x62, 0x2e, 0x76, 0x31, 0x2e, 0x47,
	0x65, 0x6f, 0x52, 0x03, 0x67, 0x65, 0x6f, 0x42, 0x0d, 0x0a, 0x0b, 0x5f, 0x6d, 0x6f, 0x64, 0x65,
	0x6c, 0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x22, 0x53, 0x0a, 0x03, 0x54, 0x61, 0x67, 0x12, 0x0e, 0x0a,
	0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x6f, 0x6c, 0x6f, 0x72, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x63, 0x6f, 0x6c, 0x6f, 0x72, 0x22, 0x29, 0x0a, 0x05, 0x50,
	0x68, 0x6f, 0x74, 0x6f, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x02, 0x69, 0x64, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x22, 0x6c, 0x0a, 0x0a, 0x50, 0x68, 0x6f, 0x74, 0x6f, 0x47,
	0x72, 0x6f, 0x75, 0x70, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x02, 0x69, 0x64, 0x12, 0x2b, 0x0a, 0x0a, 0x6d, 0x61, 0x69, 0x6e, 0x5f, 0x70, 0x68, 0x6f,
	0x74, 0x6f, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x70, 0x62, 0x2e, 0x76, 0x31,
	0x2e, 0x50, 0x68, 0x6f, 0x74, 0x6f, 0x52, 0x09, 0x6d, 0x61, 0x69, 0x6e, 0x50, 0x68, 0x6f, 0x74,
	0x6f, 0x12, 0x21, 0x0a, 0x0c, 0x70, 0x68, 0x6f, 0x74, 0x6f, 0x73, 0x5f, 0x63, 0x6f, 0x75, 0x6e,
	0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0b, 0x70, 0x68, 0x6f, 0x74, 0x6f, 0x73, 0x43,
	0x6f, 0x75, 0x6e, 0x74, 0x32, 0xf8, 0x02, 0x0a, 0x0d, 0x50, 0x68, 0x6f, 0x74, 0x6f, 0x73, 0x53,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0xb3, 0x01, 0x0a, 0x0e, 0x47, 0x65, 0x74, 0x50, 0x68,
	0x6f, 0x74, 0x6f, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x73, 0x12, 0x1c, 0x2e, 0x70, 0x62, 0x2e, 0x76,
	0x31, 0x2e, 0x47, 0x65, 0x74, 0x50, 0x68, 0x6f, 0x74, 0x6f, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x73,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1d, 0x2e, 0x70, 0x62, 0x2e, 0x76, 0x31, 0x2e,
	0x47, 0x65, 0x74, 0x50, 0x68, 0x6f, 0x74, 0x6f, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x73, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x64, 0x92, 0x41, 0x48, 0x0a, 0x05, 0x70, 0x68, 0x6f,
	0x74, 0x6f, 0x12, 0x3f, 0xd0, 0x9f, 0xd0, 0xbe, 0xd0, 0xbb, 0xd1, 0x83, 0xd1, 0x87, 0xd0, 0xb5,
	0xd0, 0xbd, 0xd0, 0xb8, 0xd0, 0xb5, 0x20, 0xd1, 0x81, 0xd0, 0xbf, 0xd0, 0xb8, 0xd1, 0x81, 0xd0,
	0xba, 0xd0, 0xb0, 0x20, 0xd0, 0xb3, 0xd1, 0x80, 0xd1, 0x83, 0xd0, 0xbf, 0xd0, 0xbf, 0x20, 0xd1,
	0x84, 0xd0, 0xbe, 0xd1, 0x82, 0xd0, 0xbe, 0xd0, 0xb3, 0xd1, 0x80, 0xd0, 0xb0, 0xd1, 0x84, 0xd0,
	0xb8, 0xd0, 0xb9, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x13, 0x12, 0x11, 0x2f, 0x76, 0x31, 0x2f, 0x70,
	0x68, 0x6f, 0x74, 0x6f, 0x73, 0x2f, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x73, 0x12, 0xb0, 0x01, 0x0a,
	0x0d, 0x47, 0x65, 0x74, 0x50, 0x68, 0x6f, 0x74, 0x6f, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x12, 0x1b,
	0x2e, 0x70, 0x62, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x50, 0x68, 0x6f, 0x74, 0x6f, 0x47,
	0x72, 0x6f, 0x75, 0x70, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e, 0x70, 0x62,
	0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x50, 0x68, 0x6f, 0x74, 0x6f, 0x47, 0x72, 0x6f, 0x75,
	0x70, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x64, 0x92, 0x41, 0x3d, 0x0a, 0x05,
	0x70, 0x68, 0x6f, 0x74, 0x6f, 0x12, 0x34, 0xd0, 0x9f, 0xd0, 0xbe, 0xd0, 0xbb, 0xd1, 0x83, 0xd1,
	0x87, 0xd0, 0xb5, 0xd0, 0xbd, 0xd0, 0xb8, 0xd0, 0xb5, 0x20, 0xd0, 0xb3, 0xd1, 0x80, 0xd1, 0x83,
	0xd0, 0xbf, 0xd0, 0xbf, 0xd1, 0x8b, 0x20, 0xd1, 0x84, 0xd0, 0xbe, 0xd1, 0x82, 0xd0, 0xbe, 0xd0,
	0xb3, 0xd1, 0x80, 0xd0, 0xb0, 0xd1, 0x84, 0xd0, 0xb8, 0xd0, 0xb9, 0x82, 0xd3, 0xe4, 0x93, 0x02,
	0x1e, 0x12, 0x1c, 0x2f, 0x76, 0x31, 0x2f, 0x70, 0x68, 0x6f, 0x74, 0x6f, 0x73, 0x2f, 0x67, 0x72,
	0x6f, 0x75, 0x70, 0x73, 0x2f, 0x7b, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x5f, 0x69, 0x64, 0x7d, 0x42,
	0xc6, 0x01, 0x92, 0x41, 0x72, 0x12, 0x18, 0x0a, 0x11, 0x50, 0x68, 0x6f, 0x74, 0x6f, 0x20, 0x6c,
	0x69, 0x62, 0x72, 0x61, 0x72, 0x79, 0x20, 0x41, 0x50, 0x49, 0x32, 0x03, 0x30, 0x2e, 0x31, 0x1a,
	0x0e, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x68, 0x6f, 0x73, 0x74, 0x3a, 0x38, 0x30, 0x38, 0x30, 0x2a,
	0x01, 0x01, 0x32, 0x10, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f,
	0x6a, 0x73, 0x6f, 0x6e, 0x3a, 0x10, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x2f, 0x6a, 0x73, 0x6f, 0x6e, 0x5a, 0x1f, 0x0a, 0x1d, 0x0a, 0x06, 0x42, 0x65, 0x61, 0x72,
	0x65, 0x72, 0x12, 0x13, 0x08, 0x02, 0x1a, 0x0d, 0x41, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x20, 0x02, 0x5a, 0x4f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x6b, 0x6b, 0x69, 0x6c, 0x69, 0x6e, 0x67, 0x2f, 0x70, 0x68, 0x6f, 0x74,
	0x6f, 0x2d, 0x6c, 0x69, 0x62, 0x72, 0x61, 0x72, 0x79, 0x2f, 0x62, 0x61, 0x63, 0x6b, 0x65, 0x6e,
	0x64, 0x2f, 0x73, 0x79, 0x6e, 0x63, 0x5f, 0x70, 0x68, 0x6f, 0x74, 0x6f, 0x73, 0x5f, 0x73, 0x65,
	0x72, 0x76, 0x65, 0x72, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f,
	0x67, 0x65, 0x6e, 0x2f, 0x70, 0x62, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_v1_photos_service_proto_rawDescOnce sync.Once
	file_proto_v1_photos_service_proto_rawDescData = file_proto_v1_photos_service_proto_rawDesc
)

func file_proto_v1_photos_service_proto_rawDescGZIP() []byte {
	file_proto_v1_photos_service_proto_rawDescOnce.Do(func() {
		file_proto_v1_photos_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_v1_photos_service_proto_rawDescData)
	})
	return file_proto_v1_photos_service_proto_rawDescData
}

var file_proto_v1_photos_service_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_proto_v1_photos_service_proto_goTypes = []interface{}{
	(*GetPhotoGroupsRequest)(nil),  // 0: pb.v1.GetPhotoGroupsRequest
	(*GetPhotoGroupsResponse)(nil), // 1: pb.v1.GetPhotoGroupsResponse
	(*GetPhotoGroupRequest)(nil),   // 2: pb.v1.GetPhotoGroupRequest
	(*GetPhotoGroupResponse)(nil),  // 3: pb.v1.GetPhotoGroupResponse
	(*PhotoWithData)(nil),          // 4: pb.v1.PhotoWithData
	(*Geo)(nil),                    // 5: pb.v1.Geo
	(*Metadata)(nil),               // 6: pb.v1.Metadata
	(*Tag)(nil),                    // 7: pb.v1.Tag
	(*Photo)(nil),                  // 8: pb.v1.Photo
	(*PhotoGroup)(nil),             // 9: pb.v1.PhotoGroup
	(*timestamppb.Timestamp)(nil),  // 10: google.protobuf.Timestamp
}
var file_proto_v1_photos_service_proto_depIdxs = []int32{
	9,  // 0: pb.v1.GetPhotoGroupsResponse.items:type_name -> pb.v1.PhotoGroup
	4,  // 1: pb.v1.GetPhotoGroupResponse.main_photo:type_name -> pb.v1.PhotoWithData
	4,  // 2: pb.v1.GetPhotoGroupResponse.photos:type_name -> pb.v1.PhotoWithData
	6,  // 3: pb.v1.PhotoWithData.meta_data:type_name -> pb.v1.Metadata
	7,  // 4: pb.v1.PhotoWithData.tags:type_name -> pb.v1.Tag
	10, // 5: pb.v1.Metadata.data_time:type_name -> google.protobuf.Timestamp
	10, // 6: pb.v1.Metadata.update_at:type_name -> google.protobuf.Timestamp
	5,  // 7: pb.v1.Metadata.geo:type_name -> pb.v1.Geo
	8,  // 8: pb.v1.PhotoGroup.main_photo:type_name -> pb.v1.Photo
	0,  // 9: pb.v1.PhotosService.GetPhotoGroups:input_type -> pb.v1.GetPhotoGroupsRequest
	2,  // 10: pb.v1.PhotosService.GetPhotoGroup:input_type -> pb.v1.GetPhotoGroupRequest
	1,  // 11: pb.v1.PhotosService.GetPhotoGroups:output_type -> pb.v1.GetPhotoGroupsResponse
	3,  // 12: pb.v1.PhotosService.GetPhotoGroup:output_type -> pb.v1.GetPhotoGroupResponse
	11, // [11:13] is the sub-list for method output_type
	9,  // [9:11] is the sub-list for method input_type
	9,  // [9:9] is the sub-list for extension type_name
	9,  // [9:9] is the sub-list for extension extendee
	0,  // [0:9] is the sub-list for field type_name
}

func init() { file_proto_v1_photos_service_proto_init() }
func file_proto_v1_photos_service_proto_init() {
	if File_proto_v1_photos_service_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_v1_photos_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetPhotoGroupsRequest); i {
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
		file_proto_v1_photos_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetPhotoGroupsResponse); i {
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
		file_proto_v1_photos_service_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetPhotoGroupRequest); i {
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
		file_proto_v1_photos_service_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetPhotoGroupResponse); i {
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
		file_proto_v1_photos_service_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PhotoWithData); i {
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
		file_proto_v1_photos_service_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Geo); i {
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
		file_proto_v1_photos_service_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Metadata); i {
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
		file_proto_v1_photos_service_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Tag); i {
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
		file_proto_v1_photos_service_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Photo); i {
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
		file_proto_v1_photos_service_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PhotoGroup); i {
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
	file_proto_v1_photos_service_proto_msgTypes[6].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_v1_photos_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_v1_photos_service_proto_goTypes,
		DependencyIndexes: file_proto_v1_photos_service_proto_depIdxs,
		MessageInfos:      file_proto_v1_photos_service_proto_msgTypes,
	}.Build()
	File_proto_v1_photos_service_proto = out.File
	file_proto_v1_photos_service_proto_rawDesc = nil
	file_proto_v1_photos_service_proto_goTypes = nil
	file_proto_v1_photos_service_proto_depIdxs = nil
}

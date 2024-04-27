// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: proto/v1/photo_tags.proto

package pbv1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	PhotoTagsService_AddPhotoTag_FullMethodName      = "/pb.v1.PhotoTagsService/AddPhotoTag"
	PhotoTagsService_GetPhotoTags_FullMethodName     = "/pb.v1.PhotoTagsService/GetPhotoTags"
	PhotoTagsService_DeletePhotoTag_FullMethodName   = "/pb.v1.PhotoTagsService/DeletePhotoTag"
	PhotoTagsService_GetTagCategories_FullMethodName = "/pb.v1.PhotoTagsService/GetTagCategories"
)

// PhotoTagsServiceClient is the client API for PhotoTagsService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PhotoTagsServiceClient interface {
	// --- PhotoTags ---
	AddPhotoTag(ctx context.Context, in *AddPhotoTagRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	GetPhotoTags(ctx context.Context, in *GetPhotoTagsRequest, opts ...grpc.CallOption) (*GetPhotoTagsResponse, error)
	DeletePhotoTag(ctx context.Context, in *DeletePhotoTagRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// --- TAGS ---
	GetTagCategories(ctx context.Context, in *GetTagCategoriesRequest, opts ...grpc.CallOption) (*GetTagCategoriesResponse, error)
}

type photoTagsServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPhotoTagsServiceClient(cc grpc.ClientConnInterface) PhotoTagsServiceClient {
	return &photoTagsServiceClient{cc}
}

func (c *photoTagsServiceClient) AddPhotoTag(ctx context.Context, in *AddPhotoTagRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, PhotoTagsService_AddPhotoTag_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *photoTagsServiceClient) GetPhotoTags(ctx context.Context, in *GetPhotoTagsRequest, opts ...grpc.CallOption) (*GetPhotoTagsResponse, error) {
	out := new(GetPhotoTagsResponse)
	err := c.cc.Invoke(ctx, PhotoTagsService_GetPhotoTags_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *photoTagsServiceClient) DeletePhotoTag(ctx context.Context, in *DeletePhotoTagRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, PhotoTagsService_DeletePhotoTag_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *photoTagsServiceClient) GetTagCategories(ctx context.Context, in *GetTagCategoriesRequest, opts ...grpc.CallOption) (*GetTagCategoriesResponse, error) {
	out := new(GetTagCategoriesResponse)
	err := c.cc.Invoke(ctx, PhotoTagsService_GetTagCategories_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PhotoTagsServiceServer is the server API for PhotoTagsService service.
// All implementations should embed UnimplementedPhotoTagsServiceServer
// for forward compatibility
type PhotoTagsServiceServer interface {
	// --- PhotoTags ---
	AddPhotoTag(context.Context, *AddPhotoTagRequest) (*emptypb.Empty, error)
	GetPhotoTags(context.Context, *GetPhotoTagsRequest) (*GetPhotoTagsResponse, error)
	DeletePhotoTag(context.Context, *DeletePhotoTagRequest) (*emptypb.Empty, error)
	// --- TAGS ---
	GetTagCategories(context.Context, *GetTagCategoriesRequest) (*GetTagCategoriesResponse, error)
}

// UnimplementedPhotoTagsServiceServer should be embedded to have forward compatible implementations.
type UnimplementedPhotoTagsServiceServer struct {
}

func (UnimplementedPhotoTagsServiceServer) AddPhotoTag(context.Context, *AddPhotoTagRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddPhotoTag not implemented")
}
func (UnimplementedPhotoTagsServiceServer) GetPhotoTags(context.Context, *GetPhotoTagsRequest) (*GetPhotoTagsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPhotoTags not implemented")
}
func (UnimplementedPhotoTagsServiceServer) DeletePhotoTag(context.Context, *DeletePhotoTagRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeletePhotoTag not implemented")
}
func (UnimplementedPhotoTagsServiceServer) GetTagCategories(context.Context, *GetTagCategoriesRequest) (*GetTagCategoriesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTagCategories not implemented")
}

// UnsafePhotoTagsServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PhotoTagsServiceServer will
// result in compilation errors.
type UnsafePhotoTagsServiceServer interface {
	mustEmbedUnimplementedPhotoTagsServiceServer()
}

func RegisterPhotoTagsServiceServer(s grpc.ServiceRegistrar, srv PhotoTagsServiceServer) {
	s.RegisterService(&PhotoTagsService_ServiceDesc, srv)
}

func _PhotoTagsService_AddPhotoTag_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddPhotoTagRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PhotoTagsServiceServer).AddPhotoTag(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PhotoTagsService_AddPhotoTag_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PhotoTagsServiceServer).AddPhotoTag(ctx, req.(*AddPhotoTagRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PhotoTagsService_GetPhotoTags_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPhotoTagsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PhotoTagsServiceServer).GetPhotoTags(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PhotoTagsService_GetPhotoTags_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PhotoTagsServiceServer).GetPhotoTags(ctx, req.(*GetPhotoTagsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PhotoTagsService_DeletePhotoTag_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeletePhotoTagRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PhotoTagsServiceServer).DeletePhotoTag(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PhotoTagsService_DeletePhotoTag_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PhotoTagsServiceServer).DeletePhotoTag(ctx, req.(*DeletePhotoTagRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PhotoTagsService_GetTagCategories_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTagCategoriesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PhotoTagsServiceServer).GetTagCategories(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PhotoTagsService_GetTagCategories_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PhotoTagsServiceServer).GetTagCategories(ctx, req.(*GetTagCategoriesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// PhotoTagsService_ServiceDesc is the grpc.ServiceDesc for PhotoTagsService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PhotoTagsService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pb.v1.PhotoTagsService",
	HandlerType: (*PhotoTagsServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddPhotoTag",
			Handler:    _PhotoTagsService_AddPhotoTag_Handler,
		},
		{
			MethodName: "GetPhotoTags",
			Handler:    _PhotoTagsService_GetPhotoTags_Handler,
		},
		{
			MethodName: "DeletePhotoTag",
			Handler:    _PhotoTagsService_DeletePhotoTag_Handler,
		},
		{
			MethodName: "GetTagCategories",
			Handler:    _PhotoTagsService_GetTagCategories_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/v1/photo_tags.proto",
}

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: proto/v1/photos_service.proto

package pbv1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	PhotosService_GetPhotoGroups_FullMethodName = "/pb.v1.PhotosService/GetPhotoGroups"
)

// PhotosServiceClient is the client API for PhotosService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PhotosServiceClient interface {
	GetPhotoGroups(ctx context.Context, in *GetPhotoGroupsRequest, opts ...grpc.CallOption) (*GetPhotoGroupsResponse, error)
}

type photosServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPhotosServiceClient(cc grpc.ClientConnInterface) PhotosServiceClient {
	return &photosServiceClient{cc}
}

func (c *photosServiceClient) GetPhotoGroups(ctx context.Context, in *GetPhotoGroupsRequest, opts ...grpc.CallOption) (*GetPhotoGroupsResponse, error) {
	out := new(GetPhotoGroupsResponse)
	err := c.cc.Invoke(ctx, PhotosService_GetPhotoGroups_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PhotosServiceServer is the server API for PhotosService service.
// All implementations should embed UnimplementedPhotosServiceServer
// for forward compatibility
type PhotosServiceServer interface {
	GetPhotoGroups(context.Context, *GetPhotoGroupsRequest) (*GetPhotoGroupsResponse, error)
}

// UnimplementedPhotosServiceServer should be embedded to have forward compatible implementations.
type UnimplementedPhotosServiceServer struct {
}

func (UnimplementedPhotosServiceServer) GetPhotoGroups(context.Context, *GetPhotoGroupsRequest) (*GetPhotoGroupsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPhotoGroups not implemented")
}

// UnsafePhotosServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PhotosServiceServer will
// result in compilation errors.
type UnsafePhotosServiceServer interface {
	mustEmbedUnimplementedPhotosServiceServer()
}

func RegisterPhotosServiceServer(s grpc.ServiceRegistrar, srv PhotosServiceServer) {
	s.RegisterService(&PhotosService_ServiceDesc, srv)
}

func _PhotosService_GetPhotoGroups_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPhotoGroupsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PhotosServiceServer).GetPhotoGroups(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PhotosService_GetPhotoGroups_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PhotosServiceServer).GetPhotoGroups(ctx, req.(*GetPhotoGroupsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// PhotosService_ServiceDesc is the grpc.ServiceDesc for PhotosService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PhotosService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pb.v1.PhotosService",
	HandlerType: (*PhotosServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetPhotoGroups",
			Handler:    _PhotosService_GetPhotoGroups_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/v1/photos_service.proto",
}

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: proto/v1/photo_metadata.proto

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
	PhotoMetadataService_GetPhotoMetaData_FullMethodName = "/pb.v1.PhotoMetadataService/GetPhotoMetaData"
)

// PhotoMetadataServiceClient is the client API for PhotoMetadataService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PhotoMetadataServiceClient interface {
	// --- MetaData ---
	GetPhotoMetaData(ctx context.Context, in *GetPhotoMetaDataRequest, opts ...grpc.CallOption) (*Metadata, error)
}

type photoMetadataServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPhotoMetadataServiceClient(cc grpc.ClientConnInterface) PhotoMetadataServiceClient {
	return &photoMetadataServiceClient{cc}
}

func (c *photoMetadataServiceClient) GetPhotoMetaData(ctx context.Context, in *GetPhotoMetaDataRequest, opts ...grpc.CallOption) (*Metadata, error) {
	out := new(Metadata)
	err := c.cc.Invoke(ctx, PhotoMetadataService_GetPhotoMetaData_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PhotoMetadataServiceServer is the server API for PhotoMetadataService service.
// All implementations should embed UnimplementedPhotoMetadataServiceServer
// for forward compatibility
type PhotoMetadataServiceServer interface {
	// --- MetaData ---
	GetPhotoMetaData(context.Context, *GetPhotoMetaDataRequest) (*Metadata, error)
}

// UnimplementedPhotoMetadataServiceServer should be embedded to have forward compatible implementations.
type UnimplementedPhotoMetadataServiceServer struct {
}

func (UnimplementedPhotoMetadataServiceServer) GetPhotoMetaData(context.Context, *GetPhotoMetaDataRequest) (*Metadata, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPhotoMetaData not implemented")
}

// UnsafePhotoMetadataServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PhotoMetadataServiceServer will
// result in compilation errors.
type UnsafePhotoMetadataServiceServer interface {
	mustEmbedUnimplementedPhotoMetadataServiceServer()
}

func RegisterPhotoMetadataServiceServer(s grpc.ServiceRegistrar, srv PhotoMetadataServiceServer) {
	s.RegisterService(&PhotoMetadataService_ServiceDesc, srv)
}

func _PhotoMetadataService_GetPhotoMetaData_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPhotoMetaDataRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PhotoMetadataServiceServer).GetPhotoMetaData(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PhotoMetadataService_GetPhotoMetaData_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PhotoMetadataServiceServer).GetPhotoMetaData(ctx, req.(*GetPhotoMetaDataRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// PhotoMetadataService_ServiceDesc is the grpc.ServiceDesc for PhotoMetadataService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PhotoMetadataService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pb.v1.PhotoMetadataService",
	HandlerType: (*PhotoMetadataServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetPhotoMetaData",
			Handler:    _PhotoMetadataService_GetPhotoMetaData_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/v1/photo_metadata.proto",
}

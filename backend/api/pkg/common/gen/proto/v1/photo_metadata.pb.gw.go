// Code generated by protoc-gen-grpc-gateway. DO NOT EDIT.
// source: proto/v1/photo_metadata.proto

/*
Package pbv1 is a reverse proxy.

It translates gRPC into RESTful JSON APIs.
*/
package pbv1

import (
	"context"
	"io"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/grpc-ecosystem/grpc-gateway/v2/utilities"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

// Suppress "imported and not used" errors
var _ codes.Code
var _ io.Reader
var _ status.Status
var _ = runtime.String
var _ = utilities.NewDoubleArray
var _ = metadata.Join

func request_PhotoMetadataService_GetPhotoMetaData_0(ctx context.Context, marshaler runtime.Marshaler, client PhotoMetadataServiceClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq GetPhotoMetaDataRequest
	var metadata runtime.ServerMetadata

	var (
		val string
		ok  bool
		err error
		_   = err
	)

	val, ok = pathParams["photo_id"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "photo_id")
	}

	protoReq.PhotoId, err = runtime.String(val)
	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "photo_id", err)
	}

	msg, err := client.GetPhotoMetaData(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err

}

func local_request_PhotoMetadataService_GetPhotoMetaData_0(ctx context.Context, marshaler runtime.Marshaler, server PhotoMetadataServiceServer, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq GetPhotoMetaDataRequest
	var metadata runtime.ServerMetadata

	var (
		val string
		ok  bool
		err error
		_   = err
	)

	val, ok = pathParams["photo_id"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "photo_id")
	}

	protoReq.PhotoId, err = runtime.String(val)
	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "photo_id", err)
	}

	msg, err := server.GetPhotoMetaData(ctx, &protoReq)
	return msg, metadata, err

}

// RegisterPhotoMetadataServiceHandlerServer registers the http handlers for service PhotoMetadataService to "mux".
// UnaryRPC     :call PhotoMetadataServiceServer directly.
// StreamingRPC :currently unsupported pending https://github.com/grpc/grpc-go/issues/906.
// Note that using this registration option will cause many gRPC library features to stop working. Consider using RegisterPhotoMetadataServiceHandlerFromEndpoint instead.
func RegisterPhotoMetadataServiceHandlerServer(ctx context.Context, mux *runtime.ServeMux, server PhotoMetadataServiceServer) error {

	mux.Handle("GET", pattern_PhotoMetadataService_GetPhotoMetaData_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		var stream runtime.ServerTransportStream
		ctx = grpc.NewContextWithServerTransportStream(ctx, &stream)
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		var err error
		var annotatedContext context.Context
		annotatedContext, err = runtime.AnnotateIncomingContext(ctx, mux, req, "/pb.v1.PhotoMetadataService/GetPhotoMetaData", runtime.WithHTTPPathPattern("/v1/photos/{photo_id}/metadata"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := local_request_PhotoMetadataService_GetPhotoMetaData_0(annotatedContext, inboundMarshaler, server, req, pathParams)
		md.HeaderMD, md.TrailerMD = metadata.Join(md.HeaderMD, stream.Header()), metadata.Join(md.TrailerMD, stream.Trailer())
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_PhotoMetadataService_GetPhotoMetaData_0(annotatedContext, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	return nil
}

// RegisterPhotoMetadataServiceHandlerFromEndpoint is same as RegisterPhotoMetadataServiceHandler but
// automatically dials to "endpoint" and closes the connection when "ctx" gets done.
func RegisterPhotoMetadataServiceHandlerFromEndpoint(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error) {
	conn, err := grpc.DialContext(ctx, endpoint, opts...)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if cerr := conn.Close(); cerr != nil {
				grpclog.Infof("Failed to close conn to %s: %v", endpoint, cerr)
			}
			return
		}
		go func() {
			<-ctx.Done()
			if cerr := conn.Close(); cerr != nil {
				grpclog.Infof("Failed to close conn to %s: %v", endpoint, cerr)
			}
		}()
	}()

	return RegisterPhotoMetadataServiceHandler(ctx, mux, conn)
}

// RegisterPhotoMetadataServiceHandler registers the http handlers for service PhotoMetadataService to "mux".
// The handlers forward requests to the grpc endpoint over "conn".
func RegisterPhotoMetadataServiceHandler(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return RegisterPhotoMetadataServiceHandlerClient(ctx, mux, NewPhotoMetadataServiceClient(conn))
}

// RegisterPhotoMetadataServiceHandlerClient registers the http handlers for service PhotoMetadataService
// to "mux". The handlers forward requests to the grpc endpoint over the given implementation of "PhotoMetadataServiceClient".
// Note: the gRPC framework executes interceptors within the gRPC handler. If the passed in "PhotoMetadataServiceClient"
// doesn't go through the normal gRPC flow (creating a gRPC client etc.) then it will be up to the passed in
// "PhotoMetadataServiceClient" to call the correct interceptors.
func RegisterPhotoMetadataServiceHandlerClient(ctx context.Context, mux *runtime.ServeMux, client PhotoMetadataServiceClient) error {

	mux.Handle("GET", pattern_PhotoMetadataService_GetPhotoMetaData_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		var err error
		var annotatedContext context.Context
		annotatedContext, err = runtime.AnnotateContext(ctx, mux, req, "/pb.v1.PhotoMetadataService/GetPhotoMetaData", runtime.WithHTTPPathPattern("/v1/photos/{photo_id}/metadata"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_PhotoMetadataService_GetPhotoMetaData_0(annotatedContext, inboundMarshaler, client, req, pathParams)
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_PhotoMetadataService_GetPhotoMetaData_0(annotatedContext, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	return nil
}

var (
	pattern_PhotoMetadataService_GetPhotoMetaData_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1, 1, 0, 4, 1, 5, 2, 2, 3}, []string{"v1", "photos", "photo_id", "metadata"}, ""))
)

var (
	forward_PhotoMetadataService_GetPhotoMetaData_0 = runtime.ForwardResponseMessage
)

package service

import (
	"context"
	"fmt"
	pbv1 "github.com/kkiling/photo-library/backend/api/pkg/common/gen/proto/v1"
	"google.golang.org/grpc"
)

type PhotoService struct {
}

func (s *PhotoService) NewTestTestInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		fmt.Println(info.FullMethod)
		return handler(ctx, req)
	}
}

func NewPhotoService() *PhotoService {
	return &PhotoService{}
}

func (s *PhotoService) CheckHashPhoto(ctx context.Context, request *pbv1.CheckHashPhotoRequest) (*pbv1.CheckHashPhotoResponse, error) {
	return &pbv1.CheckHashPhotoResponse{AlreadyLoaded: true}, nil
}

func (s *PhotoService) UploadPhoto(ctx context.Context, request *pbv1.UploadPhotoRequest) (*pbv1.UploadPhotoResponse, error) {
	return &pbv1.UploadPhotoResponse{
		Success: false,
	}, nil
}

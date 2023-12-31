package server

import (
	pbv1 "github.com/kkiling/photo-library/backend/api/pkg/common/gen/proto/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrInternal = func(err error) error {
		return NewResponseError(codes.Internal, err)
	}
	ErrNotFound = func(err error) error {
		return NewResponseError(codes.NotFound, err)
	}
	ErrPermissionDenied = func(err error) error {
		return NewResponseError(codes.PermissionDenied, err)
	}
	ErrInvalidArgument = func(err error) error {
		return NewResponseError(codes.InvalidArgument, err)
	}
	ErrUnauthenticated = func(err error) error {
		return NewResponseError(codes.Unauthenticated, err)
	}
	ErrAlreadyExists = func(err error) error {
		return NewResponseError(codes.AlreadyExists, err)
	}
	ErrTooManyRequests = func(err error) error {
		return NewResponseError(codes.ResourceExhausted, err)
	}
)

func NewResponseError(code codes.Code, err error) error {
	s, err := status.New(code, err.Error()).WithDetails(&pbv1.Error{
		Code: uint32(code),
	})

	if err != nil {
		return status.Errorf(codes.Internal, err.Error())
	}

	return s.Err()
}

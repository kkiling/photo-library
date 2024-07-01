package interceptor

import (
	"context"
	"fmt"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/kkiling/photo-library/backend/api/internal/ctxutils"
	"github.com/kkiling/photo-library/backend/api/internal/service/model"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/kkiling/photo-library/backend/api/pkg/common/log"
	"github.com/kkiling/photo-library/backend/api/pkg/common/server"
	methoddescriptor "github.com/kkiling/photo-library/backend/api/pkg/common/server/method_descriptor"
)

type SessionManager interface {
	GetSessionByToken(token string) (model.Session, error)
}

type ApiTokenService interface {
	GetApiToken(ctx context.Context, token string) (model.ApiToken, error)
}

type CustomDescriptor struct {
	method interface{}
	roles  []model.AuthRole
}

func NewCustomDescriptor(method interface{}, roles ...model.AuthRole) *CustomDescriptor {
	return &CustomDescriptor{
		method: method,
		roles:  roles,
	}
}

func (c *CustomDescriptor) Method() interface{} {
	return c.method
}

func getCustomDescriptor(descriptors methoddescriptor.DescriptorsMap, fullName string) *CustomDescriptor {
	ds, ok := descriptors.GetByFullName(fullName)
	if !ok {
		return nil
	}

	if result, ok := ds.(*CustomDescriptor); !ok {
		panic("cannot convert method descriptor to customDescriptor")
	} else {
		return result
	}
}

func containsRole(descriptorRoles []model.AuthRole, sessionRoles model.AuthRole) bool {
	for _, dr := range descriptorRoles {
		if dr == sessionRoles {
			return true
		}
	}
	return false
}

func NewApiTokenInterceptor(apiToken ApiTokenService) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		token, err := grpc_auth.AuthFromMD(ctx, "bearer")
		if err != nil {
			return nil, server.ErrUnauthenticated(fmt.Errorf("grpc_auth.AuthFromMD: %w", err))
		}

		value, err := apiToken.GetApiToken(ctx, token)
		if err != nil {
			return nil, server.ErrUnauthenticated(fmt.Errorf("sessionManager.GetApiToken: %w", err))
		}

		return handler(ctxutils.Set(ctx, ctxutils.ApiToken, value), req)
	}
}

func NewAuthInterceptor(descriptors methoddescriptor.DescriptorsMap, sessionManager SessionManager) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		ds := getCustomDescriptor(descriptors, info.FullMethod)
		if ds == nil {
			return nil, server.ErrUnauthenticated(methoddescriptor.ErrMethodDescriptorNotFound)
		}

		if len(ds.roles) == 0 {
			return handler(ctx, req)
		}

		token, err := grpc_auth.AuthFromMD(ctx, "bearer")
		if err != nil {
			return nil, server.ErrUnauthenticated(fmt.Errorf("grpc_auth.AuthFromMD: %w", err))
		}

		session, err := sessionManager.GetSessionByToken(token)
		if err != nil {
			return nil, server.ErrUnauthenticated(fmt.Errorf("sessionManager.GetSessionByToken: %w", err))
		}

		if !containsRole(ds.roles, session.Role) {
			return nil, server.ErrPermissionDenied(fmt.Errorf("no access"))
		}

		return handler(ctxutils.Set(ctx, ctxutils.Session, session), req)
	}
}

func NewPanicRecoverInterceptor(logger log.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				logger.Errorf("panic recovered: %v", r)
				err = server.ErrInternal(err)
			}
		}()

		resp, err = handler(ctx, req)
		return resp, err
	}
}

func NewLoggerInterceptor(logger log.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		resp, err = handler(ctx, req)
		if err == nil {
			return resp, err
		}
		switch status.Code(err) {
		case codes.Internal:
			logger.Errorf(err.Error())
		default:
			logger.Warnf(err.Error())
		}
		return resp, err
	}
}

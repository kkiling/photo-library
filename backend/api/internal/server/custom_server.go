package server

import (
	"context"
	"fmt"
	"google.golang.org/grpc"

	"github.com/kkiling/photo-library/backend/api/internal/interceptor"
	"github.com/kkiling/photo-library/backend/api/pkg/common/log"
	"github.com/kkiling/photo-library/backend/api/pkg/common/server"
	methoddescriptor "github.com/kkiling/photo-library/backend/api/pkg/common/server/method_descriptor"
)

type CustomHandlerService interface {
	server.HandlerService
	GetMethodDescriptors() []methoddescriptor.Descriptor
}

type CustomServer struct {
	server         *server.Server
	logger         log.Logger
	sessionManager interceptor.SessionManager
	services       []CustomHandlerService
}

func NewCustomServer(
	logger log.Logger,
	serverConfig server.Config,
	sessionManager interceptor.SessionManager,
	services ...CustomHandlerService,
) *CustomServer {
	return &CustomServer{
		server:         server.NewServer(logger, serverConfig),
		logger:         logger,
		services:       services,
		sessionManager: sessionManager,
	}
}

func (p *CustomServer) unaryServerInterceptors() ([]grpc.UnaryServerInterceptor, error) {
	var descriptors []methoddescriptor.Descriptor
	for _, service := range p.services {
		descriptors = append(descriptors, service.GetMethodDescriptors()...)
	}

	descriptorMap, err := methoddescriptor.NewMethodDescriptorMap(descriptors)
	if err != nil {
		return nil, fmt.Errorf("method_descriptor.NewMethodDescriptorMap: %w", err)
	}

	return []grpc.UnaryServerInterceptor{
		interceptor.NewPanicRecoverInterceptor(p.logger),
		interceptor.NewLoggerInterceptor(p.logger),
		interceptor.NewAuthInterceptor(descriptorMap, p.sessionManager),
	}, nil
}

func (p *CustomServer) Start(ctx context.Context, swaggerName string) error {
	unaryServerInterceptors, err := p.unaryServerInterceptors()
	if err != nil {
		return fmt.Errorf("p.unaryServerInterceptors(): %w", err)
	}

	p.server.WitUnaryServerInterceptor(unaryServerInterceptors...)

	var impl []server.HandlerService
	for _, service := range p.services {
		impl = append(impl, service)
	}

	if err = p.server.Start(ctx, swaggerName, impl...); err != nil {
		return fmt.Errorf("server.Start: %w", err)
	}

	return nil
}

func (p *CustomServer) Stop() {
	p.server.Stop()
}

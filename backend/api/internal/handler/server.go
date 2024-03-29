package handler

import (
	"context"
	"fmt"
	"github.com/kkiling/photo-library/backend/api/pkg/common/log"
	"github.com/kkiling/photo-library/backend/api/pkg/common/server"
	"github.com/kkiling/photo-library/backend/api/pkg/common/server/method_descriptor"
	"google.golang.org/grpc"
)

type CustomHandlerService interface {
	server.HandlerService
	GetMethodDescriptors() []methoddescriptor.Descriptor
}

type CustomServer struct {
	server   *server.Server
	logger   log.Logger
	services []CustomHandlerService
}

func NewCustomServer(logger log.Logger, serverConfig server.Config, services ...CustomHandlerService) *CustomServer {
	return &CustomServer{
		server:   server.NewServer(logger, serverConfig),
		logger:   logger,
		services: services,
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
		NewPanicRecoverInterceptor(p.logger),
		NewLoggerInterceptor(p.logger),
		NewAuthInterceptor(p.logger, descriptorMap),
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

	if err := p.server.Start(ctx, swaggerName, impl...); err != nil {
		return fmt.Errorf("server.Start: %w", err)
	}

	return nil
}

func (p *CustomServer) Stop() {
	p.server.Stop()
}

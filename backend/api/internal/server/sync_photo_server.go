package server

import (
	"context"
	"fmt"
	"github.com/kkiling/photo-library/backend/api/internal/interceptor"
	"github.com/kkiling/photo-library/backend/api/pkg/common/log"
	"github.com/kkiling/photo-library/backend/api/pkg/common/server"
	"google.golang.org/grpc"
)

type ApiServer struct {
	server   *server.Server
	logger   log.Logger
	services []server.HandlerService
	apiToken interceptor.ApiTokenService
}

func NewApiServer(
	logger log.Logger,
	serverConfig server.Config,
	apiToken interceptor.ApiTokenService,
	services ...server.HandlerService,
) *ApiServer {
	return &ApiServer{
		server:   server.NewServer(logger, serverConfig),
		logger:   logger,
		services: services,
		apiToken: apiToken,
	}
}

func (p *ApiServer) unaryServerInterceptors() ([]grpc.UnaryServerInterceptor, error) {
	return []grpc.UnaryServerInterceptor{
		interceptor.NewPanicRecoverInterceptor(p.logger),
		interceptor.NewLoggerInterceptor(p.logger),
		interceptor.NewApiTokenInterceptor(p.apiToken),
	}, nil
}

func (p *ApiServer) Start(ctx context.Context, swaggerName string) error {
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

func (p *ApiServer) Stop() {
	p.server.Stop()
}

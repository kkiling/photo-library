package server

import (
	"context"
	"fmt"
	httpSwagger "github.com/swaggo/http-swagger"
	"net"
	"net/http"
	"time"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	rn "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/kkiling/photo-library/backend/api/pkg/common/log"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

type Descriptor struct {
	GatewayRegistrar     func(context.Context, *rn.ServeMux, string, []grpc.DialOption) error
	OnRegisterGrpcServer func(grpcServer *grpc.Server)
}

type Server struct {
	cfg           Config
	logger        log.Logger
	grpcServer    *grpc.Server
	gatewayServer *http.Server
	mux           *rn.ServeMux
	opts          []grpc.DialOption
	errors        chan error
}

func NewServer(logger log.Logger, cfg Config, interceptor ...grpc.UnaryServerInterceptor) *Server {
	muxOption := rn.WithMarshalerOption(rn.MIMEWildcard, &rn.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames:   true,
			EmitUnpopulated: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{},
	})

	//r := http.NewServeMux()
	// r.Handle("/app.swagger.json", http.FileServer(http.Dir(".")))

	return &Server{
		logger: logger,
		grpcServer: grpc.NewServer(
			grpc.MaxRecvMsgSize(cfg.MaxReceiveMessageLength),
			grpc.MaxSendMsgSize(cfg.MaxSendMessageLength),
			grpc.ChainUnaryInterceptor(interceptor...),
			grpc.StreamInterceptor(grpc_prometheus.StreamServerInterceptor),
			grpc.UnaryInterceptor(grpc_prometheus.UnaryServerInterceptor),
		),
		mux: rn.NewServeMux(muxOption),
		opts: []grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithDefaultCallOptions(
				grpc.MaxCallRecvMsgSize(cfg.MaxReceiveMessageLength),
				grpc.MaxCallSendMsgSize(cfg.MaxSendMessageLength),
			),
		},
		cfg:    cfg,
		errors: make(chan error, 1),
	}
}

func (s *Server) Register(ctx context.Context, descriptor Descriptor) error {
	// Регистрация grpc сервера
	if descriptor.OnRegisterGrpcServer == nil {
		return fmt.Errorf("OnRegisterGrpcServer is requered")
	}

	descriptor.OnRegisterGrpcServer(s.grpcServer)

	// Регистрация rest api gateway
	if descriptor.GatewayRegistrar != nil {
		host := fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.GrpcPort)
		if err := descriptor.GatewayRegistrar(ctx, s.mux, host, s.opts); err != nil {
			return fmt.Errorf("failed to register HTTP server: %v", err)
		}
	}

	// После инициализации сервера:
	grpc_prometheus.Register(s.grpcServer)

	// Нужно что бы сервер сам отдавал описание методов
	// например для postman
	reflection.Register(s.grpcServer)

	return nil
}

func (s *Server) Start(name string, customHandler func(mux *http.ServeMux)) error {
	// Запуск grpc сервера
	go func(logger log.Logger) {
		netAddress := fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.GrpcPort)

		logger.Infof("start server at %s", netAddress)
		socket, err := net.Listen("tcp", netAddress)
		if err != nil {
			s.errors <- err
			return
		}
		s.errors <- s.grpcServer.Serve(socket)
	}(s.logger.Named("grpc_server"))

	// Запуск rest api сервера с gateway
	go func(logger log.Logger) {
		netAddress := fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.HttpPort)

		httpMux := http.NewServeMux()
		customHandler(httpMux)

		// OpenApi спецификация апи
		swagger := fmt.Sprintf("/%s.swagger.json", name)

		httpMux.Handle(swagger, http.FileServer(http.Dir("./swagger")))
		// Swagger в браузере
		httpMux.Handle("/swagger/", httpSwagger.Handler(
			httpSwagger.URL(swagger),
		))
		// Метрики
		httpMux.Handle("/metrics", promhttp.Handler())
		// Обрабатываем остальные запросы через gRPC-Gateway
		httpMux.Handle("/", s.mux)

		s.gatewayServer = &http.Server{
			Addr:    netAddress,
			Handler: httpMux,
		}

		logger.Infof("start gateway at %s", netAddress)
		s.errors <- s.gatewayServer.ListenAndServe()
	}(s.logger.Named("http_server"))

	return <-s.errors
}

func (s *Server) Stop() {
	// Пробуем по хорошему
	go s.grpcServer.GracefulStop()
	// Ждем
	time.Sleep(time.Duration(s.cfg.ShutdownTimeout) * time.Second)
	// Уже по плохому
	s.grpcServer.Stop()
}

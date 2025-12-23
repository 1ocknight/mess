package grpc

import (
	"fmt"
	"net"

	"github.com/TATAROmangol/mess/shared/logger"
	"github.com/TATAROmangol/mess/tokenissuer/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Config struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type Server struct {
	cfg    Config
	server *grpc.Server
}

func NewServer(cfg Config, log logger.Logger, svc service.Verify) *Server {
	interceptor := NewInterceptorImpl(log)

	srv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptor.InitLogger(log),
			interceptor.SetMetadataWithRequestID,
			interceptor.LogResponse,
		),
	)

	handler := NewHandlerImpl(svc)

	Register(srv, handler)
	reflection.Register(srv)

	return &Server{
		cfg:    cfg,
		server: srv,
	}
}

func (s *Server) Run() error {
	lis, err := net.Listen("tcp", fmt.Sprintf("%v:%v", s.cfg.Host, s.cfg.Port))
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}

	return s.server.Serve(lis)
}

func (s *Server) Stop() {
	s.server.Stop()
}

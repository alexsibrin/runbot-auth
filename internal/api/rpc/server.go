package rpc

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"net"
)

type IServiceRegister interface {
	Register(*grpc.Server)
}

type Config struct {
	Port int
}

type Server struct {
	server *grpc.Server
	config *Config
}

func NewServer(c *Config) (*Server, error) {
	server := grpc.NewServer()
	return &Server{
		server: server,
		config: c,
	}, nil
}

func (s *Server) Add(service IServiceRegister) {
	service.Register(s.server)
}

func (s *Server) Run(ctx context.Context) {
	ctx, cause := context.WithCancelCause(ctx)
	list, err := net.Listen("tcp", fmt.Sprintf(":%d", s.config.Port))
	if err != nil {
		cause(err)
		return
	}

	err = s.server.Serve(list)
	if err != nil {
		cause(err)
		return
	}
	
	<-ctx.Done()
	s.server.GracefulStop()
}

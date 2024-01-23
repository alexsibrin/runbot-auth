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
	// TODO: Add checking
	server := grpc.NewServer()
	return &Server{
		server: server,
		config: c,
	}, nil
}

func (s *Server) Add(service IServiceRegister) {
	service.Register(s.server)
}

func (s *Server) Run(ctx context.Context) error {
	// TODO: return error and handle it outside
	addr := fmt.Sprintf(":%d", s.config.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	errchan := make(chan error)

	go func() {
		errchan <- s.server.Serve(listener)
	}()

	select {
	case <-ctx.Done():
		s.server.GracefulStop()
		return nil
	case err := <-errchan:
		return err
	}
}

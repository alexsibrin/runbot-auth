package rpc

import "google.golang.org/grpc"

type IServiceRegister interface {
	Register(*Server)
}

type Config struct {
}

type Server struct {
	server *grpc.Server
}

func NewServer(c *Config) (*Server, error) {
	server := grpc.NewServer()
	return &Server{
		server: server,
	}, nil
}

func (s *Server) Add(service IServiceRegister) {
	service.Register(s)
}

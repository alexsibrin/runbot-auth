package rpc

import "google.golang.org/grpc"

type Server struct {
	server *grpc.Server
}

func NewServer() (*Server, error) {
	server := grpc.NewServer()
	return &Server{
		server: server,
	}, nil
}

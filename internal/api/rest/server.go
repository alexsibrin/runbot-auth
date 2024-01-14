package rest

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrDepServerIsNil        = errors.New("DepServer is nil")
	ErrDepServerConfigIsNil  = errors.New("DepServer.Config is nil")
	ErrDepServerHandlerIsNil = errors.New("DepServer.Handler is nil")
)

type Config struct {
	Host string
	Port string
}

type DependenciesServer struct {
	Config  *Config
	Handler http.Handler
}

type Server struct {
	config  *Config
	handler http.Handler

	server *http.Server
}

func NewServer(d *DependenciesServer) (*Server, error) {
	if d == nil {
		return nil, ErrDepServerIsNil
	}

	if d.Config == nil {
		return nil, ErrDepServerConfigIsNil
	}

	if d.Handler == nil {
		return nil, ErrDepServerHandlerIsNil
	}

	addr := fmt.Sprintf("%s:%s", d.Config.Host, d.Config.Port)

	hs := &http.Server{
		Addr:    addr,
		Handler: d.Handler,
	}

	return &Server{
		config: d.Config,
		server: hs,
	}, nil
}

func (s *Server) Run(ctx context.Context) {
	ctx, cancel := context.WithCancelCause(ctx)

	if err := s.server.ListenAndServe(); err != nil {
		cancel(err)
		return
	}

	<-ctx.Done()
	s.ShutDown(ctx)
}

func (s *Server) ShutDown(ctx context.Context) {
	s.server.Shutdown(ctx)
}

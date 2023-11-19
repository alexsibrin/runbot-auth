package rest

import (
	"context"
	"errors"
	"net/http"
	"runbot-auth/internal/config"
)

var (
	ErrDepHttpServerIsNil        = errors.New("DepHttpServer is nil")
	ErrDepHttpServerConfigIsNil  = errors.New("DepHttpServer.Config is nil")
	ErrDepHttpServerHandlerIsNil = errors.New("DepHttpServer.Handler is nil")
)

// TODO: describe config in this package

type DependenciesHttpServer struct {
	Config  *config.HttpServer
	Handler http.Handler
}

type HttpServer struct {
	config  *config.HttpServer
	handler http.Handler

	server *http.Server
}

func NewHttpServer(d *DependenciesHttpServer) (*HttpServer, error) {
	if d == nil {
		return nil, ErrDepHttpServerIsNil
	}

	if d.Config == nil {
		return nil, ErrDepHttpServerConfigIsNil
	}

	if d.Handler == nil {
		return nil, ErrDepHttpServerHandlerIsNil
	}

	hs := &http.Server{
		Addr:    d.Config.Addr,
		Handler: d.Handler,
	}

	return &HttpServer{
		config: d.Config,
		server: hs,
	}, nil
}

func (s *HttpServer) Run(ctx context.Context) {
	ctx, cancel := context.WithCancelCause(ctx)

	if err := s.server.ListenAndServe(); err != nil {
		cancel(err)
	}

	<-ctx.Done()
}

func (s *HttpServer) ShutDown(ctx context.Context) {
	s.server.Shutdown(ctx)
}

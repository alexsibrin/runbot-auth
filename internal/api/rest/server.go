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

type DepHttpServer struct {
	Config  *config.HttpServer
	Handler http.Handler
}

type HttpServer struct {
	config  *config.HttpServer
	handler http.Handler

	server *http.Server
}

func NewHttpServer(d *DepHttpServer) (*HttpServer, error) {
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

func (s *HttpServer) Run(ctx context.Context, ch chan error) {
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		ch<-s.server.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		s.ShutDown()
	case :
	}
}

func (s *HttpServer) ShutDown(ctx context.Context) {
	s.server.Shutdown(ctx)
}
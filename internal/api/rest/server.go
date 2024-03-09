package rest

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)

var (
	ErrDepServerIsNil        = errors.New("DepServer is nil")
	ErrDepServerConfigIsNil  = errors.New("DepServer.Config is nil")
	ErrDepServerHandlerIsNil = errors.New("DepServer.Handler is nil")
)

type Config struct {
	Host              string
	Port              string
	ReadTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
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
		Addr:              addr,
		Handler:           d.Handler,
		ReadTimeout:       d.Config.ReadTimeout,
		ReadHeaderTimeout: d.Config.ReadHeaderTimeout,
		WriteTimeout:      d.Config.WriteTimeout,
		IdleTimeout:       d.Config.IdleTimeout,
	}

	return &Server{
		config: d.Config,
		server: hs,
	}, nil
}

func (s *Server) Run(ctx context.Context) error {
	errchan := make(chan error)

	go func() {
		errchan <- s.server.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		return s.ShutDown(ctx)
	case err := <-errchan:
		return err
	}
}

func (s *Server) ShutDown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

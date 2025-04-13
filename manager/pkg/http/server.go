package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)

const (
	shutdownTimeout = 5 * time.Second
)

type Server struct {
	server *http.Server
}

type IHTTPHandler interface {
	Handler() http.Handler
}

func New(port int, handler IHTTPHandler) *Server {
	return &Server{
		server: &http.Server{
			Handler: handler.Handler(),
			Addr:    fmt.Sprintf(":%d", port),
		},
	}
}

func (s *Server) Run() error {
	if s.server == nil {
		return errors.New("http server is not initialized")
	}

	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (s *Server) Shutdown() error {
	if s.server == nil {
		return errors.New("http server is not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	return s.server.Shutdown(ctx)
}

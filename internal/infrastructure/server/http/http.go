package http

import (
	"context"
	"fmt"
	"github.com/jackvonhouse/car-enrichment/pkg/errors"
	"net/http"

	"github.com/jackvonhouse/car-enrichment/config"
)

type Server struct {
	server *http.Server
}

func New(
	handler http.Handler,
	config config.Server,
) Server {

	httpServer := http.Server{
		Addr:    fmt.Sprintf(":%d", config.Port),
		Handler: handler,
	}

	return Server{
		server: &httpServer,
	}
}

func (s Server) Run() error {
	err := s.server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (s Server) Shutdown(
	ctx context.Context,
) error {

	return s.server.Shutdown(ctx)
}

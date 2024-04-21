package app

import (
	"context"

	"github.com/jackvonhouse/car-enrichment/app/infrastructure"
	"github.com/jackvonhouse/car-enrichment/app/repository"
	"github.com/jackvonhouse/car-enrichment/app/service"
	"github.com/jackvonhouse/car-enrichment/app/transport"
	"github.com/jackvonhouse/car-enrichment/app/usecase"
	"github.com/jackvonhouse/car-enrichment/config"
	"github.com/jackvonhouse/car-enrichment/internal/infrastructure/server/http"
	"github.com/jackvonhouse/car-enrichment/pkg/log"
)

type App struct {
	infrastructure infrastructure.Infrastructure
	repository     repository.Repository
	service        service.Service
	useCase        usecase.UseCase
	transport      transport.Transport

	config config.Config
	logger log.Logger
	server http.Server
}

func New(
	ctx context.Context,
	config config.Config,
	logger log.Logger,
) (App, error) {

	i, err := infrastructure.New(ctx, config, logger)
	if err != nil {
		return App{}, err
	}

	r := repository.New(i, logger)
	s := service.New(r, config.API, logger)
	u := usecase.New(s, logger)
	t := transport.New(u, logger)

	httpServer := http.New(t.Router(), config.HTTP)

	return App{
		infrastructure: i,
		repository:     r,
		service:        s,
		useCase:        u,
		transport:      t,
		config:         config,
		logger:         logger,
		server:         httpServer,
	}, nil
}

func (a App) Run() error {
	a.logger.Info("running http server...")

	return a.server.Run()
}

func (a App) Shutdown(
	ctx context.Context,
) error {

	a.logger.Info("http server shutdowning..")

	if err := a.server.Shutdown(ctx); err != nil {
		return err
	}

	a.logger.Info("repository shutdowning..")

	if err := a.repository.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}

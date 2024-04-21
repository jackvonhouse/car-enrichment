package repository

import (
	"context"
	"github.com/jackvonhouse/car-enrichment/app/infrastructure"
	"github.com/jackvonhouse/car-enrichment/internal/infrastructure/postgres"
	"github.com/jackvonhouse/car-enrichment/internal/repository/car"
	"github.com/jackvonhouse/car-enrichment/internal/repository/owner"
	"github.com/jackvonhouse/car-enrichment/pkg/log"
)

type Repository struct {
	Car   car.Repository
	Owner owner.Repository

	Storage postgres.Database
}

func New(
	infrastructure infrastructure.Infrastructure,
	logger log.Logger,
) Repository {

	repositoryLogger := logger.WithField("layer", "repository")

	return Repository{
		Car:   car.New(infrastructure.Storage.Database(), repositoryLogger),
		Owner: owner.New(infrastructure.Storage.Database(), repositoryLogger),

		Storage: infrastructure.Storage,
	}
}

func (r Repository) Shutdown(
	_ context.Context,
) error {

	return r.Storage.Database().Close()
}

package service

import (
	"github.com/jackvonhouse/car-enrichment/app/repository"
	"github.com/jackvonhouse/car-enrichment/config"
	"github.com/jackvonhouse/car-enrichment/internal/service/car"
	"github.com/jackvonhouse/car-enrichment/internal/service/enrichment"
	"github.com/jackvonhouse/car-enrichment/internal/service/owner"
	"github.com/jackvonhouse/car-enrichment/pkg/log"
)

type Service struct {
	Car        car.Service
	Owner      owner.Service
	Enrichment enrichment.Service
}

func New(
	repository repository.Repository,
	config config.API,
	logger log.Logger,
) Service {

	serviceLogger := logger.WithField("layer", "service")

	return Service{
		Enrichment: enrichment.New(config, serviceLogger),
		Car:        car.New(repository.Car, serviceLogger),
		Owner:      owner.New(repository.Owner, serviceLogger),
	}
}

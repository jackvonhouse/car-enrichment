package usecase

import (
	"github.com/jackvonhouse/car-enrichment/app/service"
	"github.com/jackvonhouse/car-enrichment/internal/usecase/car"
	"github.com/jackvonhouse/car-enrichment/pkg/log"
)

type UseCase struct {
	Car car.UseCase
}

func New(
	service service.Service,
	logger log.Logger,
) UseCase {

	useCaseLogger := logger.WithField("layer", "usecase")

	return UseCase{
		Car: car.New(service.Car, service.Owner, service.Enrichment, useCaseLogger),
	}
}

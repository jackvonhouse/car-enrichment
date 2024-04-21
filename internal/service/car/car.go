package car

import (
	"context"
	"github.com/jackvonhouse/car-enrichment/internal/dto"
	"github.com/jackvonhouse/car-enrichment/pkg/log"
)

type carRepository interface {
	Create(context.Context, map[int64]dto.Car) error

	Get(context.Context, dto.Filter, dto.Pagination) ([]dto.Car, error)
	GetById(context.Context, int64) (dto.Car, error)

	Update(context.Context, dto.Car) error

	Delete(context.Context, dto.Car) error
}

type Service struct {
	car carRepository

	logger log.Logger
}

func New(
	car carRepository,
	logger log.Logger,
) Service {

	return Service{
		car:    car,
		logger: logger.WithField("unit", "car"),
	}
}

func (s Service) Create(
	ctx context.Context,
	cars map[int64]dto.Car,
) error {

	return s.car.Create(ctx, cars)
}

func (s Service) Get(
	ctx context.Context,
	filter dto.Filter,
	pagination dto.Pagination,
) ([]dto.Car, error) {

	return s.car.Get(ctx, filter, pagination)
}

func (s Service) GetById(
	ctx context.Context,
	id int64,
) (dto.Car, error) {

	return s.car.GetById(ctx, id)
}

func (s Service) Update(
	ctx context.Context,
	update dto.Car,
) error {

	_, err := s.car.GetById(ctx, update.ID)
	if err != nil {
		s.logger.Infof("can't get car by id (%d): %s", update.ID, err)

		return err
	}

	s.logger.Debug("car found")

	return s.car.Update(ctx, update)
}

func (s Service) Delete(
	ctx context.Context,
	carId int64,
) error {

	car, err := s.car.GetById(ctx, carId)
	if err != nil {
		s.logger.Infof("can't get car by id (%d): %s", carId, err)

		return err
	}

	s.logger.Debug("car found")

	return s.car.Delete(ctx, car)
}

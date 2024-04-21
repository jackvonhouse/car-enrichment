package owner

import (
	"context"
	"github.com/jackvonhouse/car-enrichment/internal/dto"
	"github.com/jackvonhouse/car-enrichment/pkg/log"
)

type ownerRepository interface {
	Create(context.Context, dto.CreateOwner) (int64, error)

	GetByCarId(context.Context, int64) (dto.Owner, error)
}

type Service struct {
	owner ownerRepository

	logger log.Logger
}

func New(
	owner ownerRepository,
	logger log.Logger,
) Service {

	return Service{
		owner:  owner,
		logger: logger.WithField("unit", "owner"),
	}
}

func (s Service) Create(
	ctx context.Context,
	create dto.CreateOwner,
) (int64, error) {

	return s.owner.Create(ctx, create)
}

func (s Service) GetByCarId(
	ctx context.Context,
	carId int64,
) (dto.Owner, error) {

	return s.owner.GetByCarId(ctx, carId)
}

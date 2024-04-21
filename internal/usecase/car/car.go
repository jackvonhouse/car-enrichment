package car

import (
	"context"
	"github.com/jackvonhouse/car-enrichment/internal/dto"
	"github.com/jackvonhouse/car-enrichment/internal/errors"
	errpkg "github.com/jackvonhouse/car-enrichment/pkg/errors"
	"github.com/jackvonhouse/car-enrichment/pkg/log"
	"sync"
)

type ownerService interface {
	Create(context.Context, dto.CreateOwner) (int64, error)

	GetByCarId(context.Context, int64) (dto.Owner, error)
}

type carService interface {
	Create(context.Context, map[int64]dto.Car) error

	Get(context.Context, dto.Filter, dto.Pagination) ([]dto.Car, error)
	GetById(context.Context, int64) (dto.Car, error)

	Update(context.Context, dto.Car) error

	Delete(context.Context, int64) error
}

type enrichmentService interface {
	Enrichment(context.Context, []string) (map[int64]dto.Car, error)
}

type UseCase struct {
	car   carService
	owner ownerService

	enrichment enrichmentService

	logger log.Logger
}

func New(
	car carService,
	owner ownerService,
	enrichment enrichmentService,
	logger log.Logger,
) UseCase {

	return UseCase{
		car:        car,
		owner:      owner,
		enrichment: enrichment,
		logger:     logger.WithField("unit", "car"),
	}
}

func (u UseCase) Create(
	ctx context.Context,
	create dto.CreateCar,
) (map[int64]string, error) {

	u.logger.Debug("starting enrichment")

	enrichmentCars, err := u.enrichment.Enrichment(ctx, create.RegNumbers)
	if err != nil {
		return map[int64]string{}, err
	}

	if len(enrichmentCars) == 0 {
		u.logger.Warnf("enrichment failed for cars")

		return map[int64]string{}, errors.ErrInternal.New("enrichment failed for cars").Wrap(err)
	}

	u.logger.Debug("enrichment finished")

	failedEnrichmentCars := map[int64]string{}

	if len(enrichmentCars) != len(create.RegNumbers) {
		failedEnrichmentCars = u.getFailedEnrichmentCars(create.RegNumbers, enrichmentCars)
		u.logger.Debugf("failed enrichment cars: %d", len(failedEnrichmentCars))

		u.logger.Debugf(
			"not all cars are enriched (%d), actual (%d)",
			len(create.RegNumbers),
			len(enrichmentCars),
		)
	}

	u.logger.Debug("starting create owners")

	enrichmentCarsWithOwners, err := u.createOrGetOwners(ctx, enrichmentCars)
	if err != nil {
		return failedEnrichmentCars, err
	}

	u.logger.Debugf(
		"enrichment cars after join owners: %d",
		len(enrichmentCarsWithOwners),
	)

	return failedEnrichmentCars, u.car.Create(ctx, enrichmentCarsWithOwners)
}

func (u UseCase) getFailedEnrichmentCars(
	regNumbers []string,
	enrichmentCars map[int64]dto.Car,
) map[int64]string {

	failed := map[int64]string{}

	for id, regNumber := range regNumbers {
		id := int64(id)

		if _, ok := enrichmentCars[id]; !ok {
			failed[id] = regNumber
		}
	}

	return failed
}

func (u UseCase) createOrGetOwners(
	ctx context.Context,
	enrichmentCars map[int64]dto.Car,
) (map[int64]dto.Car, error) {

	mu := &sync.Mutex{}
	wg := &sync.WaitGroup{}
	wg.Add(len(enrichmentCars))

	enrichmentCarsCopy := make(map[int64]dto.Car, len(enrichmentCars))

	for carId, car := range enrichmentCars {
		enrichmentCarsCopy[carId] = car
	}

	var ownerId int64

	for carId, car := range enrichmentCarsCopy {
		u.logger.Debug("starting join owner to car with id %d", carId)

		go func(carId int64, car dto.Car) {
			defer wg.Done()

			owner, err := u.owner.GetByCarId(ctx, carId)

			if err == nil {
				u.logger.Infof("owner for car id (%d) already exists", carId)

				ownerId = owner.ID
			} else {
				u.logger.Infof("can't get owner by car id (%d): %s", carId, err)
				u.logger.Debugf("try to create owner for car id (%d)", carId)

				ownerId, err = u.owner.Create(ctx, dto.CreateOwner{
					Name:       car.Owner.Name,
					Surname:    car.Owner.Surname,
					Patronymic: car.Owner.Patronymic,
				})

				if err != nil && !errpkg.TypeIs(err, errors.ErrAlreadyExists) {
					u.logger.Warnf("can't create owner for car id (%d): %s", carId, err)

					mu.Lock()
					delete(enrichmentCarsCopy, carId)
					mu.Unlock()
				}

				u.logger.Debugf("owner created for car id (%d)", carId)
			}

			mu.Lock()
			if c, ok := enrichmentCarsCopy[carId]; ok {
				c.Owner.ID = ownerId
				enrichmentCarsCopy[carId] = c
			}
			mu.Unlock()
		}(carId, car)

		u.logger.Debugf("joined owner to car with id %d", carId)
	}

	wg.Wait()

	return enrichmentCarsCopy, nil
}

func (u UseCase) Get(
	ctx context.Context,
	filter dto.Filter,
	pagination dto.Pagination,
) ([]dto.Car, error) {

	return u.car.Get(ctx, filter, pagination)
}

func (u UseCase) Update(
	ctx context.Context,
	car dto.Car,
) error {

	owner, err := u.owner.GetByCarId(ctx, car.ID)
	if err != nil {
		u.logger.Warnf("can't get owner by car id (%d): %s", car.ID, err)

		return err
	}

	car.Owner.ID = owner.ID

	return u.car.Update(ctx, car)
}

func (u UseCase) Delete(
	ctx context.Context,
	carId int64,
) error {

	return u.car.Delete(ctx, carId)
}

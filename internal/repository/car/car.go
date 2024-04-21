package car

import (
	"context"
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	pgerr "github.com/jackc/pgerrcode"
	"github.com/jackvonhouse/car-enrichment/internal/dto"
	"github.com/jackvonhouse/car-enrichment/internal/errors"
	errpkg "github.com/jackvonhouse/car-enrichment/pkg/errors"
	"github.com/jackvonhouse/car-enrichment/pkg/log"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"strings"
)

type Repository struct {
	db *sqlx.DB

	logger log.Logger
}

func New(
	db *sqlx.DB,
	logger log.Logger,
) Repository {

	return Repository{
		db:     db,
		logger: logger.WithField("unit", "car"),
	}
}

func (r Repository) Create(
	ctx context.Context,
	cars map[int64]dto.Car,
) error {

	insertBuilder := sq.
		Insert("car").
		Columns("regNum", "mark", "model", "year", "owner_id")

	for _, car := range cars {
		insertBuilder = insertBuilder.Values(
			car.RegNum, car.Mark, car.Model, car.Year, car.Owner.ID,
		)
	}

	query, args, err := insertBuilder.PlaceholderFormat(sq.Dollar).ToSql()

	logger := r.logger.WithFields(map[string]any{
		"request": map[string]any{
			"query": query,
			"args": map[string]any{
				"cars": cars,
			},
		},
	})

	if err != nil {
		logger.Warnf("error on create sql query: %s", err)

		return errors.ErrInternal.New("can't create car").Wrap(err)
	}

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		if errpkg.Is(err, sql.ErrNoRows) {
			logger.Warnf("can't create car: %s", err)

			return errors.ErrInternal.New("can't create car").Wrap(err)
		}

		if e, ok := err.(*pq.Error); ok {
			switch e.Code {

			case pgerr.UniqueViolation:
				logger.Warnf("car already exists: %s", err)

				return errors.ErrAlreadyExists.New("car already exists").Wrap(err)

			case pgerr.ForeignKeyViolation:
				logger.Warnf("owner not found: %s", err)

				return errors.ErrNotFound.New("owner not found").Wrap(err)

			default:
				logger.Warnf("can't create car: %s", err)

				return errors.ErrInternal.New("can't create car").Wrap(err)
			}
		}

		logger.Warnf("can't create car: %s", err)

		return errors.ErrInternal.New("can't create car").Wrap(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected != int64(len(cars)) {
		logger.Warnf("can't create car: %s", err)

		return errors.ErrInternal.New("can't create car").Wrap(err)
	}

	logger.Debugf("rows affected: %d", rowsAffected)

	return nil
}

func (r Repository) Get(
	ctx context.Context,
	filter dto.Filter,
	pagination dto.Pagination,
) ([]dto.Car, error) {

	var (
		offset = uint64(pagination.Offset)
		limit  = uint64(pagination.Limit)
	)

	selectBuilder := sq.
		Select(
			"car.id AS car_id",
			"car.regnum AS car_regnum",
			"car.mark AS car_mark",
			"car.model AS car_model",
			"car.year AS car_year",
			"owner.id AS owner_id",
			"owner.name AS owner_name",
			"owner.surname AS owner_surname",
			"owner.patronymic AS owner_patronymic",
		).
		From("car").
		LeftJoin("owner ON car.owner_id = owner.id").
		OrderBy("car.id DESC").
		Offset(offset).
		Limit(limit).
		PlaceholderFormat(sq.Dollar)

	selectBuilder = r.where(selectBuilder, filter)

	query, args, err := selectBuilder.ToSql()

	logger := r.logger.WithFields(map[string]any{
		"query": query,
		"args": map[string]any{
			"limit":  pagination.Limit,
			"offset": pagination.Offset,
		},
	})

	if err != nil {
		logger.Warnf("can't get cars: %s", err)

		return []dto.Car{}, errors.ErrInternal.New("can't get cars").Wrap(err)
	}

	type car struct {
		CarID           int64  `db:"car_id"`
		RegNum          string `db:"car_regnum"`
		Mark            string `db:"car_mark"`
		Model           string `db:"car_model"`
		Year            int    `db:"car_year"`
		OwnerID         int64  `db:"owner_id"`
		OwnerName       string `db:"owner_name"`
		OwnerSurname    string `db:"owner_surname"`
		OwnerPatronymic string `db:"owner_patronymic"`
	}

	rawCars := make([]car, 0)

	if err := r.db.SelectContext(ctx, &rawCars, query, args...); err != nil {
		if !errpkg.Is(err, sql.ErrNoRows) {
			logger.Warnf("can't get cars: %s", err)

			return []dto.Car{}, errors.ErrInternal.New("can't get cars").Wrap(err)
		}

		logger.Warnf("no cars: %s", err)

		return []dto.Car{}, errors.ErrNotFound.New("no cars").Wrap(err)
	}

	cars := make([]dto.Car, len(rawCars))
	for i, rawCar := range rawCars {
		cars[i] = dto.Car{
			ID:     rawCar.CarID,
			RegNum: rawCar.RegNum,
			Mark:   rawCar.Mark,
			Model:  rawCar.Model,
			Year:   rawCar.Year,
			Owner: dto.Owner{
				ID:         rawCar.OwnerID,
				Name:       rawCar.OwnerName,
				Surname:    rawCar.OwnerSurname,
				Patronymic: rawCar.OwnerPatronymic,
			},
		}
	}

	return cars, nil
}

func (r Repository) GetById(
	ctx context.Context,
	id int64,
) (dto.Car, error) {

	query, args, err := sq.
		Select(
			"car.id AS car_id",
			"car.regnum AS car_regnum",
			"car.mark AS car_mark",
			"car.model AS car_model",
			"car.year AS car_year",
			"owner.id AS owner_id",
			"owner.name AS owner_name",
			"owner.surname AS owner_surname",
			"owner.patronymic AS owner_patronymic",
		).
		From("car").
		LeftJoin("owner ON car.owner_id = owner.id").
		Where(sq.Eq{"car.id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	logger := r.logger.WithFields(map[string]any{
		"query": query,
		"args": map[string]any{
			"id": id,
		},
	})

	if err != nil {
		logger.Warnf("can't get car: %s", err)

		return dto.Car{}, errors.ErrInternal.New("can't get car").Wrap(err)
	}

	type car struct {
		CarID           int64  `db:"car_id"`
		RegNum          string `db:"car_regnum"`
		Mark            string `db:"car_mark"`
		Model           string `db:"car_model"`
		Year            int    `db:"car_year"`
		OwnerID         int64  `db:"owner_id"`
		OwnerName       string `db:"owner_name"`
		OwnerSurname    string `db:"owner_surname"`
		OwnerPatronymic string `db:"owner_patronymic"`
	}

	rawCar := car{}

	if err := r.db.GetContext(ctx, &rawCar, query, args...); err != nil {
		logger.Warnf("can't get car: %s", err)

		if !errpkg.Is(err, sql.ErrNoRows) {
			return dto.Car{}, errors.ErrInternal.New("can't get car").Wrap(err)
		}

		return dto.Car{}, errors.ErrNotFound.New("car not found").Wrap(err)
	}

	return dto.Car{
		ID:     rawCar.CarID,
		RegNum: rawCar.RegNum,
		Mark:   rawCar.Mark,
		Model:  rawCar.Model,
		Year:   rawCar.Year,
		Owner: dto.Owner{
			ID:         rawCar.OwnerID,
			Name:       rawCar.OwnerName,
			Surname:    rawCar.OwnerSurname,
			Patronymic: rawCar.OwnerPatronymic,
		},
	}, nil
}

func (r Repository) Update(
	ctx context.Context,
	update dto.Car,
) error {

	rollback := func(tx *sqlx.Tx) error {
		if err := tx.Rollback(); err != nil {
			r.logger.Warnf("unknown error on rollback: %s", err)

			return errors.ErrInternal.New("can't update car").Wrap(err)
		}

		return nil
	}

	commit := func(tx *sqlx.Tx) error {
		if err := tx.Commit(); err != nil {
			r.logger.Warnf("unknown error on commit: %s", err)

			return errors.ErrInternal.New("can't update car").Wrap(err)
		}

		return nil
	}

	r.logger.Info("starting update car transaction")

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		r.logger.Warnf("can't start transaction: %s", err)

		return errors.ErrInternal.New("can't update car").Wrap(err)
	}

	if err := r.updateCar(ctx, tx, update); err != nil {
		if err := rollback(tx); err != nil {
			return err
		}

		return err
	}

	if err := r.updateOwner(ctx, tx, update.Owner); err != nil {
		if err := rollback(tx); err != nil {
			return err
		}

		return err
	}

	return commit(tx)
}

func (r Repository) updateCar(
	ctx context.Context,
	tx *sqlx.Tx,
	update dto.Car,
) error {

	values := r.updateCarValues(update)

	query, args, err := sq.
		Update("car").
		SetMap(values).
		Where(sq.Eq{"id": update.ID}).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar).
		ToSql()

	logger := r.logger.WithFields(map[string]any{
		"request": map[string]any{
			"query": query,
			"args":  values,
		},
	})

	if err != nil {
		logger.Warnf("error on create sql query: %s", err)

		return nil
	}

	var carId int64

	if err := tx.GetContext(ctx, &carId, query, args...); err != nil {
		if errpkg.Is(err, sql.ErrNoRows) {
			logger.Warnf("can't update car: %s", err)

			return errors.ErrInternal.New("can't update car").Wrap(err)
		}

		if e, ok := err.(*pq.Error); ok {
			switch e.Code {

			case pgerr.UniqueViolation:
				logger.Warnf("car already exists: %s", err)

				return errors.ErrAlreadyExists.New("car already exists").Wrap(err)

			case pgerr.ForeignKeyViolation:
				logger.Warnf("owner not found: %s", err)

				return errors.ErrNotFound.New("owner not found").Wrap(err)

			default:
				logger.Warnf("can't update car: %s", err)

				return errors.ErrInternal.New("can't update car").Wrap(err)
			}
		}

		logger.Warnf("can't update car: %s", err)

		return errors.ErrInternal.New("can't update car").Wrap(err)
	}

	return nil
}

func (r Repository) updateOwner(
	ctx context.Context,
	tx *sqlx.Tx,
	owner dto.Owner,
) error {

	values := r.updateOwnerValues(owner)

	query, args, err := sq.
		Update("owner").
		SetMap(values).
		Where(sq.Eq{"id": owner.ID}).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar).
		ToSql()

	logger := r.logger.WithFields(map[string]any{
		"request": map[string]any{
			"query": query,
			"args":  values,
		},
	})

	if err != nil {
		logger.Warnf("error on create sql query: %s", err)

		return nil
	}

	var ownerId int64

	if err := tx.GetContext(ctx, &ownerId, query, args...); err != nil {
		if errpkg.Is(err, sql.ErrNoRows) {
			logger.Warnf("can't update owner: %s", err)

			return errors.ErrInternal.New("can't update owner").Wrap(err)
		}

		if e, ok := err.(*pq.Error); ok {
			switch e.Code {

			case pgerr.UniqueViolation:
				logger.Warnf("owner already exists: %s", err)

				return errors.ErrAlreadyExists.New("owner already exists").Wrap(err)

			default:
				logger.Warnf("can't update owner: %s", err)

				return errors.ErrInternal.New("can't update owner").Wrap(err)
			}
		}

		logger.Warnf("can't owner car: %s", err)

		return errors.ErrInternal.New("can't owner car").Wrap(err)
	}

	return nil
}

func (r Repository) updateCarValues(
	update dto.Car,
) map[string]any {

	u := map[string]any{}

	if len(update.Mark) != 0 {
		u["mark"] = update.Mark
	}

	if len(update.Model) != 0 {
		u["model"] = update.Model
	}

	if len(update.RegNum) != 0 {
		u["regNum"] = update.RegNum
	}

	if update.Year != 0 {
		u["year"] = update.Year
	}

	return u
}

func (r Repository) updateOwnerValues(
	owner dto.Owner,
) map[string]any {

	u := map[string]any{}

	if len(strings.TrimSpace(owner.Name)) != 0 {
		u["name"] = owner.Name
	}

	if len(strings.TrimSpace(owner.Surname)) != 0 {
		u["surname"] = owner.Surname
	}

	if len(strings.TrimSpace(owner.Patronymic)) != 0 {
		u["patronymic"] = owner.Patronymic
	}

	return u
}

func (r Repository) Delete(
	ctx context.Context,
	car dto.Car,
) error {

	query, args, err := sq.
		Delete("car").
		Where(sq.Eq{"id": car.ID}).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar).
		ToSql()

	logger := r.logger.WithFields(map[string]any{
		"query": query,
		"args": map[string]any{
			"car": car,
		},
	})

	if err != nil {
		logger.Warnf("can't get owner: %s", err)

		return errors.ErrInternal.New("can't delete owner").Wrap(err)
	}

	var carId int

	if err := r.db.GetContext(ctx, &carId, query, args...); err != nil {
		if !errpkg.Is(err, sql.ErrNoRows) {
			logger.Warnf("can't delete car: %s", err)

			return errors.ErrInternal.New("can't delete car").Wrap(err)
		}

		logger.Warnf("car not found: %s", err)

		return errors.ErrNotFound.New("car not found").Wrap(err)
	}

	return nil
}

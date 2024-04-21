package owner

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
	create dto.CreateOwner,
) (int64, error) {

	query, args, err := sq.
		Insert("owner").
		Columns("name", "surname", "patronymic").
		Values(create.Name, create.Surname, create.Patronymic).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar).
		ToSql()

	logger := r.logger.WithFields(map[string]any{
		"request": map[string]any{
			"query": query,
			"args": map[string]any{
				"owner": map[string]any{
					"name":       create.Name,
					"surname":    create.Surname,
					"patronymic": create.Patronymic,
				},
			},
		},
	})

	if err != nil {
		logger.Warnf("error on create sql query: %s", err)

		return 0, errors.ErrInternal.New("can't create car").Wrap(err)
	}

	var ownerId int64

	if err := r.db.GetContext(ctx, &ownerId, query, args...); err != nil {
		if errpkg.Is(err, sql.ErrNoRows) {
			logger.Warnf("can't create car owner: %s", err)

			return 0, errors.ErrInternal.New("can't create car owner").Wrap(err)
		}

		if e, ok := err.(*pq.Error); ok {
			switch e.Code {

			case pgerr.UniqueViolation:
				logger.Warnf("car owner already exists: %s", err)

				return 0, errors.ErrAlreadyExists.New("car owner already exists").Wrap(err)

			default:
				logger.Warnf("can't create car owner: %s", err)

				return 0, errors.ErrInternal.New("can't create car owner").Wrap(err)
			}
		}

		logger.Warnf("can't create car owner: %s", err)

		return 0, errors.ErrInternal.New("can't create car owner").Wrap(err)
	}

	return ownerId, nil
}

func (r Repository) GetByCarId(
	ctx context.Context,
	carId int64,
) (dto.Owner, error) {

	query, args, err := sq.
		Select("owner.*").
		From("car").
		LeftJoin("owner ON car.owner_id = owner.id").
		Where(sq.Eq{"car.id": carId}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	logger := r.logger.WithFields(map[string]any{
		"query": query,
		"args": map[string]any{
			"car": map[string]any{
				"id": carId,
			},
		},
	})

	if err != nil {
		logger.Warnf("can't get owner: %s", err)

		return dto.Owner{}, errors.ErrInternal.New("can't get owner").Wrap(err)
	}

	type car struct {
		ID         int64  `db:"id"`
		Name       string `db:"name"`
		Surname    string `db:"surname"`
		Patronymic string `db:"patronymic"`
	}

	rawOwner := car{}

	if err := r.db.GetContext(ctx, &rawOwner, query, args...); err != nil {
		logger.Warnf("can't get owner: %s", err)

		if !errpkg.Is(err, sql.ErrNoRows) {
			return dto.Owner{}, errors.ErrInternal.New("can't get owner").Wrap(err)
		}

		return dto.Owner{}, errors.ErrNotFound.New("owner not found").Wrap(err)
	}

	return dto.Owner{
		ID:         rawOwner.ID,
		Name:       rawOwner.Name,
		Surname:    rawOwner.Surname,
		Patronymic: rawOwner.Patronymic,
	}, nil
}

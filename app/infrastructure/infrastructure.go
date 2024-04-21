package infrastructure

import (
	"context"
	"github.com/jackvonhouse/car-enrichment/config"
	"github.com/jackvonhouse/car-enrichment/internal/infrastructure/postgres"
	"github.com/jackvonhouse/car-enrichment/pkg/log"
)

type Infrastructure struct {
	Storage postgres.Database
}

func New(
	ctx context.Context,
	config config.Config,
	logger log.Logger,
) (Infrastructure, error) {

	infrastructureLog := logger.WithField("layer", "infrastructure")

	db, err := postgres.New(ctx, config.Database, infrastructureLog)
	if err != nil {
		infrastructureLog.Warn(err)

		return Infrastructure{}, err
	}

	return Infrastructure{
		Storage: db,
	}, nil
}

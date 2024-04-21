package main

import (
	"context"
	"flag"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackvonhouse/car-enrichment/app"
	"github.com/jackvonhouse/car-enrichment/config"
	"github.com/jackvonhouse/car-enrichment/pkg/log"
	"github.com/jackvonhouse/car-enrichment/pkg/shutdown"
)

func getConfigPath() string {
	var configPath string

	flag.StringVar(
		&configPath,
		"config",
		"config/config.toml",
		"The path to the configuration file",
	)

	flag.Parse()

	return configPath
}

// @title			Каталог автомобилей
// @version		1.0
// @description	Простейшее API для каталога автомобилей
// @host		localhost:8081
// @BasePath	/api/v1
func main() {
	ctx, cancel := context.WithCancel(context.Background())

	logger := log.NewLogrusLogger()

	logger.Info("application started")
	defer func() {
		logger.Info("application stopped")
	}()

	configPath := getConfigPath()
	cfg, err := config.New(
		configPath,
		logger.WithField("layer", "config"),
	)

	if err != nil {
		logger.Warnf("failed to load config: %s", err)

		return
	}

	application, err := app.New(ctx, cfg, logger)
	if err != nil {
		logger.Warnf("failed to create application: %s", err)

		return
	}

	go application.Run()

	shutdown.Graceful(ctx, cancel, logger, application)
}

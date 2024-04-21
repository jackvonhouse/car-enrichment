package config

import (
	"fmt"
	"github.com/jackvonhouse/car-enrichment/internal/errors"
	"github.com/jackvonhouse/car-enrichment/pkg/log"
	"github.com/spf13/viper"
	"path/filepath"
	"strings"
)

type Database struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
	SSLMode  string
}

func (d Database) String() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		d.Username, d.Password,
		d.Host, d.Port,
		d.Database,
		d.SSLMode,
	)
}

type Server struct {
	Port int
}

type API struct {
	Url string
}

type Config struct {
	Database Database
	HTTP     Server
	API      API
}

func New(
	path string,
	logger log.Logger,
) (Config, error) {

	configType := strings.TrimPrefix(filepath.Ext(path), ".")

	configLogger := logger.WithField("path", path)
	configLogger.Infof("config type: %s", configType)

	viper.SetConfigType(configType)
	viper.SetConfigFile(path)

	if err := viper.ReadInConfig(); err != nil {
		configLogger.Warnf("failed to read config: %s", err)

		return Config{}, errors.ErrInvalid.New("invalid config").Wrap(err)
	}

	pgPrefix := "database.postgres"
	httpPrefix := "server.http"
	apiPrefix := "api"

	return Config{
		Database: Database{
			Host:     viper.GetString(fmt.Sprintf("%s.host", pgPrefix)),
			Port:     viper.GetInt(fmt.Sprintf("%s.port", pgPrefix)),
			Username: viper.GetString(fmt.Sprintf("%s.username", pgPrefix)),
			Password: viper.GetString(fmt.Sprintf("%s.password", pgPrefix)),
			Database: viper.GetString(fmt.Sprintf("%s.database", pgPrefix)),
			SSLMode:  viper.GetString(fmt.Sprintf("%s.ssl_mode", pgPrefix)),
		},

		HTTP: Server{
			Port: viper.GetInt(fmt.Sprintf("%s.port", httpPrefix)),
		},

		API: API{
			Url: viper.GetString(fmt.Sprintf("%s.url", apiPrefix)),
		},
	}, nil
}

package config

import (
	"gitlab.hotel.tools/backend-team/allocation-go/internal/common/adapter/server/http"
	"os"
	"time"

	"go.uber.org/config"
	"go.uber.org/zap"

	"gitlab.hotel.tools/backend-team/allocation-go/internal/common/adapter/db"
)

type Config struct {
	Location string `yaml:"location"`

	Name      string `yaml:"name"`
	Version   string `yaml:"version"`
	BuildDate string `yaml:"build_date"`

	Server *http.Config `yaml:"server"`
	Logger *zap.Config  `yaml:"logger"`
	DB     *db.Config   `yaml:"db"`
}

var cfg *Config

const defaultConfigYaml = "./config.yaml"

var (
	name      = "saga-svc"
	version   = "undefined/local"
	buildDate = time.Now().Format(time.RFC3339)
)

func NewConfig() (*Config, error) {
	configFile := os.Getenv("CONFIG_FILE")
	if configFile == "" {
		configFile = defaultConfigYaml
	}

	provider, err := NewProviderByOptions(config.File(configFile))
	if err != nil {
		return nil, err
	}
	if err = provider.Populate(&cfg); err != nil {
		panic(err)
	}

	cfg.Name = name
	cfg.Version = version
	cfg.BuildDate = buildDate

	return cfg, nil
}

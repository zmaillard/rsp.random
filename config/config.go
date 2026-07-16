package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	DBUser            string `envconfig:"DB_USER"`
	DBHost            string `envconfig:"DB_HOST"`
	DBPassword        string `envconfig:"DB_PASSWORD"`
	DBName            string `envconfig:"DB_NAME"`
	DBPort            string `envconfig:"DB_PORT"`
	BaseUrl           string `envconfig:"BASE_URL"`
	LoadDataAtStartup bool   `envconfig:"LOAD_DATA_STARTUP"`
	VersionNumber     string
}

func NewConfig() (*Config, error) {
	return NewConfigWithVersion("development")
}

func NewConfigWithVersion(version string) (*Config, error) {
	var config Config
	err := envconfig.Process("", &config)
	config.VersionNumber = version
	return &config, err
}

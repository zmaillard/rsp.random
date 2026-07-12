package config

import (
	"net/url"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	DBUser        string `envconfig:"DB_USER"`
	DBHost        string `envconfig:"DB_HOST"`
	DBPassword    string `envconfig:"DB_PASSWORD"`
	DBName        string `envconfig:"DB_NAME"`
	DBPort        string `envconfig:"DB_PORT"`
	SearchServer  string `envconfig:"SEARCH_SERVER"`
	SignIndex     string `envconfig:"SEARCH_SIGN_INDEX"`
	SearchApiKey  string `envconfig:"SEARCH_API_TOKEN"`
	BaseUrl       string `envconfig:"BASE_URL"`
	VersionNumber string
}

func (c *Config) GetSearchUrl() (string, error) {
	return url.JoinPath(c.SearchServer, "indexes", c.SignIndex, "search")
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

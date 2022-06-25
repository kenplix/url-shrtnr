package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

type (
	// Config -.
	Config struct {
		App    `yaml:"app"`
		HTTP   `yaml:"http"`
		Logger `yaml:"logger"`
	}

	// App -.
	App struct {
		Name    string `yaml:"name" env:"APP_NAME" env-required:"true"`
		Version string `yaml:"version" env:"APP_VERSION" env-required:"true"`
	}

	// HTTP -.
	HTTP struct {
		Port string `yaml:"port" env:"HTTP_PORT" env-required:"true"`
	}

	// Logger -.
	Logger struct {
		Level           string `yaml:"level" env:"LOG_LEVEL" env-required:"true"`
		TimestampFormat string `yaml:"timestamp_format" env:"TIMESTAMP_FORMAT" env-required:"true"`
	}
)

// NewConfig returns app config.
func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := cleanenv.ReadConfig("./config/config.yml", cfg); err != nil {
		return nil, fmt.Errorf("config reading error: %w", err)
	}

	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, fmt.Errorf("env reading error: %w", err)
	}

	return cfg, nil
}

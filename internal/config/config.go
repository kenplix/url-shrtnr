package config

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/Kenplix/url-shrtnr/internal/repository"
	"github.com/Kenplix/url-shrtnr/pkg/hash"
	"github.com/Kenplix/url-shrtnr/pkg/httpserver"
	"github.com/Kenplix/url-shrtnr/pkg/logger"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

const EnvPrefix = "URL_SHRTNR"

// Config -.
type Config struct {
	Environment string            `mapstructure:"environment"`
	HTTP        httpserver.Config `mapstructure:"http"`
	Database    repository.Config `mapstructure:"database"`
	Logger      logger.Config     `mapstructure:"logger"`
	Hasher      hash.Config       `mapstructure:"hasher"`
}

// Read -.
func Read(dir string) (Config, error) {
	viper.AutomaticEnv()
	viper.SetEnvPrefix(EnvPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	var cfg Config
	if err := load(&cfg); err != nil {
		return Config{}, errors.Wrap(err, "error loading config")
	}

	if err := read(dir); err != nil {
		return Config{}, errors.Wrap(err, "error reading config")
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return Config{}, errors.Wrap(err, "error unmarshalling config")
	}

	return cfg, nil
}

func load(cfg *Config) error {
	keys := map[string]interface{}{}
	if err := mapstructure.Decode(cfg, &keys); err != nil {
		return errors.Wrap(err, "error decoding config keys")
	}

	buf, err := json.Marshal(keys)
	if err != nil {
		return errors.Wrap(err, "error marshaling config")
	}

	viper.SetConfigType("json")
	err = viper.ReadConfig(bytes.NewReader(buf))

	return errors.Wrap(err, "error reading config")
}

func read(dir string) error {
	viper.AddConfigPath(dir)
	viper.SetConfigType("yaml")

	file := "development"
	if env := viper.GetString("ENVIRONMENT"); env != "" {
		file = env
	}

	viper.SetConfigName(file)

	err := viper.MergeInConfig()

	return errors.Wrapf(err, "error merging with %q config file", viper.ConfigFileUsed())
}

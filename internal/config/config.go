package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/kenplix/url-shrtnr/pkg/log"

	"github.com/kenplix/url-shrtnr/internal/service"

	"github.com/kenplix/url-shrtnr/pkg/cache/redis"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"github.com/kenplix/url-shrtnr/internal/repository"
	"github.com/kenplix/url-shrtnr/pkg/hash"
	"github.com/kenplix/url-shrtnr/pkg/httpserver"
)

const EnvPrefix = "URL_SHRTNR"

type Environment string

const (
	OmitEnvironment        Environment = ""
	TestingEnvironment     Environment = "testing"
	DevelopmentEnvironment Environment = "development"
	ProductionEnvironment  Environment = "production"
)

func (env Environment) Validate() error {
	switch env {
	case DevelopmentEnvironment:
	case ProductionEnvironment:
	case OmitEnvironment:
	case TestingEnvironment:
	default:
		return fmt.Errorf("unknown environment %q", env)
	}

	return nil
}

func (env Environment) ConfigName() string {
	if env == OmitEnvironment {
		return string(DevelopmentEnvironment)
	}

	return string(env)
}

// Config -.
type Config struct {
	Environment Environment              `mapstructure:"environment"`
	HTTP        httpserver.Config        `mapstructure:"http"`
	Database    repository.Config        `mapstructure:"database"`
	Logger      log.Config               `mapstructure:"logger"`
	Redis       redis.Config             `mapstructure:"redis"`
	Hasher      hash.Config              `mapstructure:"hasher"`
	JWT         service.JWTServiceConfig `mapstructure:"jwt"`
}

// Read -.
func Read(dir string) (Config, error) {
	viper.AutomaticEnv()
	viper.SetEnvPrefix(EnvPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	var cfg Config
	if err := load(&cfg); err != nil {
		return Config{}, errors.Wrap(err, "failed to load config")
	}

	if err := read(dir); err != nil {
		return Config{}, errors.Wrap(err, "failed to read config")
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return Config{}, errors.Wrap(err, "failed to unmarshall config")
	}

	return cfg, nil
}

func load(cfg *Config) error {
	keys := map[string]any{}
	if err := mapstructure.Decode(cfg, &keys); err != nil {
		return errors.Wrap(err, "failed to decode config keys")
	}

	buf, err := json.Marshal(keys)
	if err != nil {
		return errors.Wrap(err, "failed to marshall config keys")
	}

	viper.SetConfigType("json")

	err = viper.ReadConfig(bytes.NewReader(buf))
	if err != nil {
		return errors.Wrap(err, "failed to read config")
	}

	return nil
}

func read(dir string) error {
	viper.AddConfigPath(dir)
	viper.SetConfigType("yaml")

	env := Environment(viper.GetString("ENVIRONMENT"))
	if err := env.Validate(); err != nil {
		return err
	}

	viper.SetConfigName(env.ConfigName())

	err := viper.MergeInConfig()
	if err != nil {
		return errors.Wrapf(err, "failed to merge with %q config file", viper.ConfigFileUsed())
	}

	return nil
}

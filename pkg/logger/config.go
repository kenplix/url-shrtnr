package logger

import (
	"time"

	"github.com/sirupsen/logrus"
)

const (
	defaultLevel           = logrus.InfoLevel
	defaultTimestampFormat = time.Stamp
)

type Config struct {
	Level string `mapstructure:"level"`
}

func SetConfig(cfg Config) Option {
	return Preset(
		SetLevel(cfg.Level),
	)
}

package log

import "go.uber.org/zap"

const (
	consoleEncoding = "console"
	jsonEncoding    = "json"
)

const (
	developmentMode = "development"
	productionMode  = "production"
)

const (
	defaultLevel    = zap.InfoLevel
	defaultEncoding = consoleEncoding
	defaultMode     = developmentMode
)

type Config struct {
	Level    string `mapstructure:"level"`
	Encoding string `mapstructure:"encoding"`
	Mode     string `mapstructure:"mode"`
}

func SetConfig(cfg Config) Option {
	return Preset(
		SetLevel(cfg.Level),
		SetMode(cfg.Mode),
		SetEncoding(cfg.Encoding),
	)
}

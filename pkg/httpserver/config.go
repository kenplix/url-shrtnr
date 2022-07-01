package httpserver

import "time"

const (
	defaultPort            uint16 = 80
	defaultReadTimeout            = 5 * time.Second
	defaultWriteTimeout           = 5 * time.Second
	defaultShutdownTimeout        = 3 * time.Second
)

// Config structure is used to configure the Server
type Config struct {
	Port            uint16        `mapstructure:"port"`
	ReadTimeout     time.Duration `mapstructure:"readTimeout"`
	WriteTimeout    time.Duration `mapstructure:"writeTimeout"`
	IdleTimeout     time.Duration `mapstructure:"idleTimeout"`
	ShutdownTimeout time.Duration `mapstructure:"shutdownTimeout"`
}

func DefaultPreset() Option {
	return SetConfig(DefaultConfig())
}

func DefaultConfig() *Config {
	return &Config{
		Port:            defaultPort,
		ReadTimeout:     defaultReadTimeout,
		WriteTimeout:    defaultWriteTimeout,
		ShutdownTimeout: defaultShutdownTimeout,
	}
}

func SetConfig(cfg *Config) Option {
	return Preset(
		SetPort(cfg.Port),
		SetReadTimeout(cfg.ReadTimeout),
		SetIdleTimeout(cfg.IdleTimeout),
		SetWriteTimeout(cfg.WriteTimeout),
		SetShutdownTimeout(cfg.ShutdownTimeout),
	)
}

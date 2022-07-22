package httpserver

import "time"

const (
	defaultPort              = "80"
	defaultReadTimeout       = 5 * time.Second
	defaultReadHeaderTimeout = 2 * time.Second
	defaultWriteTimeout      = 5 * time.Second
	defaultShutdownTimeout   = 3 * time.Second
)

// Config structure is used to configure the Server
type Config struct {
	Port              string        `mapstructure:"port"`
	ReadTimeout       time.Duration `mapstructure:"readTimeout"`
	ReadHeaderTimeout time.Duration `mapstructure:"readHeaderTimeout"`
	WriteTimeout      time.Duration `mapstructure:"writeTimeout"`
	IdleTimeout       time.Duration `mapstructure:"idleTimeout"`
	ShutdownTimeout   time.Duration `mapstructure:"shutdownTimeout"`
}

func SetConfig(cfg Config) Option {
	return Preset(
		SetPort(cfg.Port),
		SetReadTimeout(cfg.ReadTimeout),
		SetReadHeaderTimeout(cfg.ReadHeaderTimeout),
		SetIdleTimeout(cfg.IdleTimeout),
		SetWriteTimeout(cfg.WriteTimeout),
		SetShutdownTimeout(cfg.ShutdownTimeout),
	)
}

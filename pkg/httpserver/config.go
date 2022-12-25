package httpserver

import "time"

const (
	defaultHost = ""
	defaultPort = "80"

	defaultReadTimeout       = 5 * time.Second
	defaultReadHeaderTimeout = 2 * time.Second
	defaultWriteTimeout      = 5 * time.Second
	defaultIdleTimeout       = 2 * time.Second
	defaultShutdownTimeout   = 3 * time.Second
)

// Config structure is used to configure the Server
type Config struct {
	Host              string        `mapstructure:"host"`
	Port              string        `mapstructure:"port"`
	ReadTimeout       time.Duration `mapstructure:"readTimeout"`
	ReadHeaderTimeout time.Duration `mapstructure:"readHeaderTimeout"`
	WriteTimeout      time.Duration `mapstructure:"writeTimeout"`
	IdleTimeout       time.Duration `mapstructure:"idleTimeout"`
	ShutdownTimeout   time.Duration `mapstructure:"shutdownTimeout"`
}

func SetConfig(cfg Config) Option {
	return Preset(
		SetAddr(cfg.Host, cfg.Port),
		SetReadTimeout(cfg.ReadTimeout),
		SetReadHeaderTimeout(cfg.ReadHeaderTimeout),
		SetWriteTimeout(cfg.WriteTimeout),
		SetIdleTimeout(cfg.IdleTimeout),
		SetShutdownTimeout(cfg.ShutdownTimeout),
	)
}

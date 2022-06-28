package httpserver

const DefaultAddr = ":80"

type Config struct {
	Addr string `mapstructure:"addr"`
}

func DefaultConfig() *Config {
	return &Config{
		Addr: DefaultAddr,
	}
}

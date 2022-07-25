package bcrypt

type Config struct {
	Cost int `mapstructure:"cost"`
}

func SetConfig(cfg Config) Option {
	return Preset(
		SetCost(cfg.Cost),
	)
}

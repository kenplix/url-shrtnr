package argon2

type Config struct {
	Memory      uint32 `mapstructure:"memory"`
	Iterations  uint32 `mapstructure:"iterations"`
	Parallelism uint8  `mapstructure:"parallelism"`
	SaltLength  uint32 `mapstructure:"saltLength"`
	KeyLength   uint32 `mapstructure:"keyLength"`
}

func SetConfig(cfg Config) Option {
	return Preset(
		SetMemory(cfg.Memory),
		SetIterations(cfg.Iterations),
		SetParallelism(cfg.Parallelism),
		SetSaltLength(cfg.SaltLength),
		SetKeyLength(cfg.KeyLength),
	)
}

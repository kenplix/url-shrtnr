package token

import "time"

type Config struct {
	PrivateKey string        `mapstructure:"privateKey"`
	PublicKey  string        `mapstructure:"publicKey"`
	TTL        time.Duration `mapstructure:"ttl"`
}

func SetConfig(cfg Config) Option {
	preset := Preset(
		SetPrivateKey(cfg.PrivateKey),
		SetPublicKey(cfg.PublicKey),
	)

	if cfg.TTL != 0 {
		preset = Preset(preset, SetTTL(cfg.TTL))
	}

	return preset
}

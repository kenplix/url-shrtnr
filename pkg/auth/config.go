package auth

import "time"

const (
	defaultAccessTokenTTL  = time.Minute
	defaultRefreshTokenTTL = time.Minute
)

// Config structure is used to configure the Service
type Config struct {
	AccessTokenSigningKey  string         `mapstructure:"accessTokenSigningKey"`
	AccessTokenTTL         *time.Duration `mapstructure:"accessTokenTTL"`
	RefreshTokenSigningKey string         `mapstructure:"refreshTokenSigningKey"`
	RefreshTokenTTL        *time.Duration `mapstructure:"refreshTokenTTL"`
}

func SetConfig(cfg Config) Option {
	preset := Preset(
		SetAccessTokenSigningKey(cfg.AccessTokenSigningKey),
		SetRefreshTokenSigningKey(cfg.RefreshTokenSigningKey),
	)

	if cfg.AccessTokenTTL != nil {
		preset = Preset(preset, SetAccessTokenTTL(*cfg.AccessTokenTTL))
	}

	if cfg.RefreshTokenTTL != nil {
		preset = Preset(preset, SetRefreshTokenTTL(*cfg.RefreshTokenTTL))
	}

	return preset
}

package auth

import "time"

const (
	defaultAccessTokenTTL  = 15 * time.Minute
	defaultRefreshTokenTTL = 24 * time.Hour
)

// Config structure is used to configure the TokensService
type Config struct {
	AccessToken  TokenConfig `mapstructure:"accessToken"`
	RefreshToken TokenConfig `mapstructure:"refreshToken"`
}

type TokenConfig struct {
	PrivateKey string         `mapstructure:"privateKey"`
	PublicKey  string         `mapstructure:"publicKey"`
	TTL        *time.Duration `mapstructure:"ttl"`
}

func SetConfig(cfg Config) Option {
	preset := Preset(
		SetAccessTokenPrivateKey(cfg.AccessToken.PrivateKey),
		SetAccessTokenPublicKey(cfg.AccessToken.PublicKey),
		SetRefreshTokenPrivateKey(cfg.RefreshToken.PrivateKey),
		SetRefreshTokenPublicKey(cfg.RefreshToken.PublicKey),
	)

	if ttl := cfg.AccessToken.TTL; ttl != nil {
		preset = Preset(preset, SetAccessTokenTTL(*ttl))
	}

	if ttl := cfg.RefreshToken.TTL; ttl != nil {
		preset = Preset(preset, SetRefreshTokenTTL(*ttl))
	}

	return preset
}

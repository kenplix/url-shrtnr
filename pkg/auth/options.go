package auth

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
)

// Option configures a TokensService.
type Option interface {
	apply(s *tokensService) error
}

type optionFunc func(s *tokensService) error

func (fn optionFunc) apply(s *tokensService) error {
	return fn(s)
}

// Preset turns a list of Option instances into an Option
func Preset(options ...Option) Option {
	return optionFunc(func(s *tokensService) error {
		for _, option := range options {
			if err := option.apply(s); err != nil {
				return err
			}
		}

		return nil
	})
}

func SetAccessTokenSigningKey(signingKey string) Option {
	return optionFunc(func(s *tokensService) error {
		if signingKey == "" {
			return errors.New("empty access token signing key")
		}

		s.accessTokenSigningKey = signingKey
		return nil
	})
}

func SetAccessTokenTTL(ttl time.Duration) Option {
	return optionFunc(func(s *tokensService) error {
		if ttl <= 0 {
			return fmt.Errorf("access token TTL can't be less or equal 0")
		}

		s.accessTokenTTL = ttl
		return nil
	})
}

func SetRefreshTokenSigningKey(signingKey string) Option {
	return optionFunc(func(s *tokensService) error {
		if signingKey == "" {
			return errors.New("empty refresh token signing key")
		}

		s.refreshTokenSigningKey = signingKey
		return nil
	})
}

func SetRefreshTokenTTL(ttl time.Duration) Option {
	return optionFunc(func(s *tokensService) error {
		if ttl <= 0 {
			return fmt.Errorf("refresh token TTL can't be less or equal 0")
		}

		s.refreshTokenTTL = ttl
		return nil
	})
}

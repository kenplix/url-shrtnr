package auth

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
)

// Option configures a Service.
type Option interface {
	apply(s *Service) error
}

type optionFunc func(s *Service) error

func (fn optionFunc) apply(s *Service) error {
	return fn(s)
}

// Preset turns a list of Option instances into an Option
func Preset(options ...Option) Option {
	return optionFunc(func(s *Service) error {
		for _, option := range options {
			if err := option.apply(s); err != nil {
				return err
			}
		}

		return nil
	})
}

func SetAccessTokenSigningKey(signingKey string) Option {
	return optionFunc(func(s *Service) error {
		if signingKey == "" {
			return errors.New("empty access token signing key")
		}

		s.accessTokenSigningKey = signingKey
		return nil
	})
}

func SetAccessTokenTTL(ttl time.Duration) Option {
	return optionFunc(func(s *Service) error {
		if ttl <= 0 {
			return fmt.Errorf("access token TTL can't be less or equal 0")
		}

		s.accessTokenTTL = ttl
		return nil
	})
}

func SetRefreshTokenSigningKey(signingKey string) Option {
	return optionFunc(func(s *Service) error {
		if signingKey == "" {
			return errors.New("empty refresh token signing key")
		}

		s.refreshTokenSigningKey = signingKey
		return nil
	})
}

func SetRefreshTokenTTL(ttl time.Duration) Option {
	return optionFunc(func(s *Service) error {
		if ttl <= 0 {
			return fmt.Errorf("refresh token TTL can't be less or equal 0")
		}

		s.refreshTokenTTL = ttl
		return nil
	})
}

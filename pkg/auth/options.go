package auth

import (
	"crypto/rsa"
	"encoding/base64"
	"github.com/golang-jwt/jwt"
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

func SetAccessTokenPrivateKey(privateKey string) Option {
	return optionFunc(func(s *tokensService) error {
		key, err := parseRSAPrivateKey(privateKey)
		if err != nil {
			return errors.Wrap(err, "access token")
		}

		s.accessServ.privateKey = key
		return nil
	})
}

func SetAccessTokenPublicKey(publicKey string) Option {
	return optionFunc(func(s *tokensService) error {
		key, err := parseRSAPublicKey(publicKey)
		if err != nil {
			return errors.Wrap(err, "access token")
		}

		s.accessServ.publicKey = key
		return nil
	})
}

func SetAccessTokenTTL(ttl time.Duration) Option {
	return optionFunc(func(s *tokensService) error {
		if ttl <= 0 {
			return errors.New("access token TTL can't be less or equal 0")
		} else if ttl >= s.refreshServ.ttl {
			return errors.New("access token TTL can't be greater or equal refresh token TTL")
		}

		s.accessServ.ttl = ttl
		return nil
	})
}

func SetRefreshTokenPrivateKey(privateKey string) Option {
	return optionFunc(func(s *tokensService) error {
		key, err := parseRSAPrivateKey(privateKey)
		if err != nil {
			return errors.Wrap(err, "refresh token")
		}

		s.refreshServ.privateKey = key
		return nil
	})
}

func SetRefreshTokenPublicKey(publicKey string) Option {
	return optionFunc(func(s *tokensService) error {
		key, err := parseRSAPublicKey(publicKey)
		if err != nil {
			return errors.Wrap(err, "refresh token")
		}

		s.refreshServ.publicKey = key
		return nil
	})
}

func SetRefreshTokenTTL(ttl time.Duration) Option {
	return optionFunc(func(s *tokensService) error {
		if ttl <= 0 {
			return errors.New("refresh token TTL can't be less or equal 0")
		} else if ttl <= s.accessServ.ttl {
			return errors.New("refresh token TTL can't be less or equal access token TTL")
		}

		s.refreshServ.ttl = ttl
		return nil
	})
}

func parseRSAPrivateKey(privateKey string) (*rsa.PrivateKey, error) {
	if privateKey == "" {
		return nil, errors.New("empty RSA private key")
	}

	decodedPrivateKey, err := base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode RSA private key")
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(decodedPrivateKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse RSA private key")
	}

	return key, nil
}

func parseRSAPublicKey(publicKey string) (*rsa.PublicKey, error) {
	if publicKey == "" {
		return nil, errors.New("empty RSA public key")
	}

	decodedPublicKey, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode RSA public key")
	}

	key, err := jwt.ParseRSAPublicKeyFromPEM(decodedPublicKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse RSA public key")
	}

	return key, nil
}

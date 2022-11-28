package token

import (
	"encoding/base64"
	"time"

	"github.com/golang-jwt/jwt"

	"github.com/pkg/errors"
)

// Option configures a TokensService.
type Option interface {
	apply(s *jwtService) error
}

type optionFunc func(s *jwtService) error

func (fn optionFunc) apply(s *jwtService) error {
	return fn(s)
}

// Preset turns a list of Option instances into an Option
func Preset(options ...Option) Option {
	return optionFunc(func(s *jwtService) error {
		for _, option := range options {
			if err := option.apply(s); err != nil {
				return err
			}
		}

		return nil
	})
}

func SetPrivateKey(privateKey string) Option {
	return optionFunc(func(s *jwtService) error {
		if privateKey == "" {
			return errors.New("empty RSA private key")
		}

		decodedPrivateKey, err := base64.StdEncoding.DecodeString(privateKey)
		if err != nil {
			return errors.Wrap(err, "failed to decode RSA private key")
		}

		key, err := jwt.ParseRSAPrivateKeyFromPEM(decodedPrivateKey)
		if err != nil {
			return errors.Wrap(err, "failed to parse RSA private key")
		}

		s.privateKey = key
		return nil
	})
}

func SetPublicKey(publicKey string) Option {
	return optionFunc(func(s *jwtService) error {
		if publicKey == "" {
			return errors.New("empty RSA public key")
		}

		decodedPublicKey, err := base64.StdEncoding.DecodeString(publicKey)
		if err != nil {
			return errors.Wrap(err, "failed to decode RSA public key")
		}

		key, err := jwt.ParseRSAPublicKeyFromPEM(decodedPublicKey)
		if err != nil {
			return errors.Wrap(err, "failed to parse RSA public key")
		}

		s.publicKey = key
		return nil
	})
}

func SetTTL(ttl time.Duration) Option {
	return optionFunc(func(s *jwtService) error {
		if ttl <= 0 {
			return errors.New("ttl must be positive")
		}

		s.ttl = ttl
		return nil
	})
}

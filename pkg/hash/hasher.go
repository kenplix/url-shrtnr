package hash

import (
	"fmt"

	"github.com/kenplix/url-shrtnr/pkg/hash/argon2"
	"github.com/kenplix/url-shrtnr/pkg/hash/bcrypt"
)

// HasherService provides hashing logic to securely store passwords.
//
//go:generate mockery --dir . --name HasherService --output ./mocks
type HasherService interface {
	HashPassword(password string) (string, error)
	VerifyPassword(password, hash string) bool
}

type Config struct {
	Use    string        `mapstructure:"use"`
	Bcrypt bcrypt.Config `mapstructure:"bcrypt"`
	Argon2 argon2.Config `mapstructure:"argon2"`
}

func NewHasherService(cfg Config) (HasherService, error) {
	switch cfg.Use {
	case "bcrypt":
		hasher := bcrypt.NewHasherService(bcrypt.SetConfig(cfg.Bcrypt))
		return hasher, nil
	case "argon2":
		hasher := argon2.NewHasherService(argon2.SetConfig(cfg.Argon2))
		return hasher, nil
	default:
		return nil, fmt.Errorf("unknown hasher %q", cfg.Use)
	}
}

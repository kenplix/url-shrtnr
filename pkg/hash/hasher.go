package hash

import (
	"fmt"

	"github.com/Kenplix/url-shrtnr/pkg/hash/argon2"
	"github.com/Kenplix/url-shrtnr/pkg/hash/bcrypt"
)

// Hasher provides hashing logic to securely store passwords.
//go:generate mockery --dir . --name Hasher --output ./mocks
type Hasher interface {
	HashPassword(password string) (string, error)
	CheckPasswordHash(password, hash string) bool
}

type Config struct {
	Use    string        `mapstructure:"use"`
	Bcrypt bcrypt.Config `mapstructure:"bcrypt"`
	Argon2 argon2.Config `mapstructure:"argon2"`
}

func New(cfg Config) (Hasher, error) {
	switch cfg.Use {
	case "bcrypt":
		hasher := bcrypt.NewHasher(bcrypt.SetConfig(cfg.Bcrypt))
		return hasher, nil
	case "argon2":
		hasher := argon2.NewHasher(argon2.SetConfig(cfg.Argon2))
		return hasher, nil
	default:
		return nil, fmt.Errorf("unknown hasher %q", cfg.Use)
	}
}

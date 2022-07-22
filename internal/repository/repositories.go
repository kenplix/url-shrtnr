package repository

import (
	"context"
	"fmt"

	"github.com/Kenplix/url-shrtnr/internal/entity"
	"github.com/Kenplix/url-shrtnr/internal/repository/mongodb"
)

// UsersRepository is a store for users
//go:generate mockery --dir . --name UsersRepository --output ./mocks
type UsersRepository interface {
	Create(ctx context.Context, user entity.User) error
	GetByCredentials(ctx context.Context, email, password string) (entity.User, error)
}

type Config struct {
	Use     string         `mapstructure:"use"`
	MongoDB mongodb.Config `mapstructure:"mongodb"`
}

// Repositories -.
type Repositories struct {
	Users UsersRepository
}

func New(ctx context.Context, cfg Config) (*Repositories, error) {
	switch cfg.Use {
	case "mongodb":
		db, err := mongodb.New(ctx, cfg.MongoDB)
		if err != nil {
			return nil, err
		}

		repo := &Repositories{
			Users: db.UsersRepository(),
		}

		return repo, nil
	default:
		return nil, fmt.Errorf("unknown database %q", cfg.Use)
	}
}

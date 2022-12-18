package repository

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/kenplix/url-shrtnr/internal/entity"
	"github.com/kenplix/url-shrtnr/internal/repository/mongodb"
)

// UsersRepository is a store for users
//
//go:generate mockery --dir . --name UsersRepository --output ./mocks
type UsersRepository interface {
	Create(ctx context.Context, user entity.User) error
	FindByID(ctx context.Context, userID primitive.ObjectID) (entity.User, error)
	FindByUsername(ctx context.Context, username string) (entity.User, error)
	FindByEmail(ctx context.Context, email string) (entity.User, error)
	FindByLogin(ctx context.Context, login string) (entity.User, error)
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

		usersRepo, err := db.UsersRepository(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create users repository")
		}

		r := &Repositories{
			Users: usersRepo,
		}

		return r, nil
	default:
		return nil, fmt.Errorf("unknown database %q", cfg.Use)
	}
}

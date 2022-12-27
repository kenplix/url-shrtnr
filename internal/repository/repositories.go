package repository

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/kenplix/url-shrtnr/internal/entity"
)

// UsersRepository is a store for users
//
//go:generate mockery --dir . --name UsersRepository --output ./mocks
type UsersRepository interface {
	Create(ctx context.Context, user entity.UserModel) error
	FindByID(ctx context.Context, userID primitive.ObjectID) (entity.UserModel, error)
	FindByUsername(ctx context.Context, username string) (entity.UserModel, error)
	FindByEmail(ctx context.Context, email string) (entity.UserModel, error)
	FindByLogin(ctx context.Context, login string) (entity.UserModel, error)
	ChangePassword(ctx context.Context, userID primitive.ObjectID, passwordHash string) error
}

type Config struct {
	Use     string        `mapstructure:"use"`
	MongoDB MongoDBConfig `mapstructure:"mongodb"`
	FileDB  FileDBConfig  `mapstructure:"filedb"`
}

// Repositories -.
type Repositories struct {
	Users UsersRepository
	close func(ctx context.Context) error
}

func New(ctx context.Context, cfg Config) (*Repositories, error) {
	f, err := createDatabaseFactory(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create database factory")
	}

	db, err := f.make(ctx)
	if err != nil {
		return nil, err
	}

	r := &Repositories{
		Users: db.getUsersRepository(),
		close: db.close,
	}

	return r, nil
}

type database interface {
	getUsersRepository() UsersRepository
	close(ctx context.Context) error
}

type databaseMaker interface {
	make(ctx context.Context) (database, error)
}

func createDatabaseFactory(cfg Config) (databaseMaker, error) {
	switch cfg.Use {
	case "mongodb":
		return &mongoDBMaker{config: cfg.MongoDB}, nil
	case "filedb":
		return &fileDBMaker{config: cfg.FileDB}, nil
	default:
		return nil, fmt.Errorf("unknown database %q", cfg.Use)
	}
}

type mongoDBMaker struct {
	config MongoDBConfig
}

func (m *mongoDBMaker) make(ctx context.Context) (database, error) {
	db, err := newMongoDB(ctx, m.config)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create mongodb")
	}

	err = db.createUsersRepository(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create users repository")
	}

	return db, nil
}

type fileDBMaker struct {
	config FileDBConfig
}

func (m *fileDBMaker) make(_ context.Context) (database, error) {
	db, err := newFileDB(m.config)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create filedb")
	}

	err = db.createUsersRepository()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create users repository")
	}

	return db, nil
}

func (r *Repositories) Close(ctx context.Context) error {
	return r.close(ctx)
}

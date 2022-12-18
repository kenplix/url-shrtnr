package repository

import (
	"context"

	"github.com/pkg/errors"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/kenplix/url-shrtnr/pkg/database/mongodb"
)

type MongoDBConfig struct {
	URI      string `mapstructure:"uri"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
}

type mongoDB struct {
	client *mongo.Client
	db     *mongo.Database
	users  UsersRepository
}

func newMongoDB(ctx context.Context, cfg MongoDBConfig) (*mongoDB, error) {
	client, err := mongodb.NewClient(ctx, cfg.URI, cfg.Username, cfg.Password)
	if err != nil {
		return nil, err
	}

	db := &mongoDB{
		client: client,
		db:     client.Database(cfg.Database),
	}

	return db, nil
}

func (m *mongoDB) close(ctx context.Context) error {
	if err := m.client.Disconnect(ctx); err != nil {
		return errors.Wrap(err, "failed to disconnect mongodb client")
	}

	return nil
}

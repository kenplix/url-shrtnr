package mongodb

import (
	"context"
	"github.com/pkg/errors"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/Kenplix/url-shrtnr/pkg/database/mongodb"
)

type Config struct {
	URI      string `mapstructure:"uri"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
}

type MongoDB struct {
	db *mongo.Database
}

func New(ctx context.Context, cfg Config) (*MongoDB, error) {
	client, err := mongodb.NewClient(ctx, cfg.URI, cfg.Username, cfg.Password)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create MongoDB client")
	}

	return &MongoDB{db: client.Database(cfg.Database)}, nil
}

package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/kenplix/url-shrtnr/pkg/database/mongodb"
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
		return nil, err
	}

	return &MongoDB{db: client.Database(cfg.Database)}, nil
}

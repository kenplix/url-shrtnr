package mongodb

import (
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/Kenplix/url-shrtnr/pkg/database/mongodb"
)

type Config struct {
	URI      string `mapstructure:"uri"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
}

type MongoDB struct {
	db *mongo.Database
}

func New(cfg Config) (*MongoDB, error) {
	client, err := mongodb.NewClient(cfg.URI, cfg.User, cfg.Password)
	if err != nil {
		return nil, err
	}

	return &MongoDB{db: client.Database(cfg.Name)}, nil
}

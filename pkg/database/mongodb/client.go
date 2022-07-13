package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	defaultConnectTimeout = 10 * time.Second
	defaultPingTimeout    = 2 * time.Second
)

// NewClient establish connection with MongoDB instance using provided URI and auth credentials.
func NewClient(ctx context.Context, uri, username, password string) (*mongo.Client, error) {
	opts := options.Client().ApplyURI(uri)
	if username != "" && password != "" {
		opts.SetAuth(options.Credential{
			Username: username,
			Password: password,
		})
	}

	client, err := mongo.NewClient(opts)
	if err != nil {
		return nil, err
	}

	connCtx, connCancel := context.WithTimeout(ctx, defaultConnectTimeout)
	defer connCancel()

	err = client.Connect(connCtx)
	if err != nil {
		return nil, err
	}

	pingCtx, pingCancel := context.WithTimeout(ctx, defaultPingTimeout)
	defer pingCancel()

	err = client.Ping(pingCtx, nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}

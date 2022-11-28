package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/pkg/errors"
)

const defaultPingTimeout = 2 * time.Second

type Config struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

func NewClient(ctx context.Context, cfg Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
	})

	pingCtx, pingCancel := context.WithTimeout(ctx, defaultPingTimeout)
	defer pingCancel()

	_, err := client.Ping(pingCtx).Result()
	if err != nil {
		return nil, errors.Wrap(err, "failed to ping redis server")
	}

	return client, nil
}

package app

import (
	"context"
	"net/http"

	transport "github.com/kenplix/url-shrtnr/internal/controller/http"

	"github.com/kenplix/url-shrtnr/pkg/log"

	"go.uber.org/zap"

	"github.com/kenplix/url-shrtnr/pkg/hash"

	"github.com/pkg/errors"

	"github.com/kenplix/url-shrtnr/pkg/cache/redis"

	"github.com/kenplix/url-shrtnr/internal/config"
	"github.com/kenplix/url-shrtnr/internal/repository"
	"github.com/kenplix/url-shrtnr/internal/service"
	"github.com/kenplix/url-shrtnr/pkg/httpserver"
)

// Run -.
func Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.Read("configs")
	if err != nil {
		return errors.Wrap(err, "failed to create config")
	}

	logger, err := log.NewLogger(log.SetConfig(cfg.Logger))
	if err != nil {
		return errors.Wrapf(err, "failed to init zap logger")
	}

	defer logger.Sync()

	repos, err := repository.New(ctx, cfg.Database)
	if err != nil {
		return errors.Wrap(err, "failed to create repositories")
	}
	defer repos.Close(context.TODO())

	cache, err := redis.NewClient(ctx, cfg.Redis)
	if err != nil {
		return errors.Wrap(err, "failed to create redis client")
	}

	hasherServ, err := hash.NewHasherService(cfg.Hasher)
	if err != nil {
		return errors.Wrapf(err, "failed to create hasher service")
	}

	services, err := service.NewServices(service.Dependencies{
		Cache:            cache,
		Repos:            repos,
		HasherService:    hasherServ,
		JWTServiceConfig: cfg.JWT,
	})
	if err != nil {
		return errors.Wrapf(err, "failed to create services")
	}

	handler, err := transport.NewHandler(logger, services)
	if err != nil {
		return errors.Wrap(err, "failed to create handler")
	}

	httpServer := httpserver.New(
		handler.InitEngine(),
		httpserver.SetConfig(cfg.HTTP),
	)

	addr := httpServer.Start()
	logger.Info("started HTTP server",
		zap.String("addr", addr),
	)

	if err = <-httpServer.Notify(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error("error occurred while running HTTP server",
			zap.Error(err),
		)
	}

	err = httpServer.Shutdown()

	return errors.Wrap(err, "failed to shutdown HTTP server")
}

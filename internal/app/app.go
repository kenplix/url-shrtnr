package app

import (
	"context"
	"net/http"

	"github.com/pkg/errors"

	"github.com/Kenplix/url-shrtnr/internal/config"
	"github.com/Kenplix/url-shrtnr/internal/controller/http/v1"
	"github.com/Kenplix/url-shrtnr/internal/repository"
	"github.com/Kenplix/url-shrtnr/internal/service"
	"github.com/Kenplix/url-shrtnr/pkg/auth"
	"github.com/Kenplix/url-shrtnr/pkg/hash"
	"github.com/Kenplix/url-shrtnr/pkg/httpserver"
	"github.com/Kenplix/url-shrtnr/pkg/logger"
)

// Run -.
func Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.Read("configs")
	if err != nil {
		return errors.Wrap(err, "failed to create config")
	}

	log, err := logger.New(logger.SetConfig(cfg.Logger))
	if err != nil {
		return errors.Wrapf(err, "failed to create logger")
	}

	repos, err := repository.New(ctx, cfg.Database)
	if err != nil {
		return errors.Wrap(err, "failed to create repositories")
	}

	hasherServ, err := hash.NewHasherService(cfg.Hasher)
	if err != nil {
		return errors.Wrapf(err, "failed to create hasher service")
	}

	tokensServ, err := auth.NewTokensService(cfg.Authorization)
	if err != nil {
		return errors.Wrap(err, "failed to create tokens service")
	}

	services, err := service.NewServices(service.Dependencies{
		Repos:         repos,
		HasherService: hasherServ,
		TokensService: tokensServ,
	})
	if err != nil {
		return errors.Wrapf(err, "failed to create services")
	}

	handler, err := v1.NewHandler(services)
	if err != nil {
		return errors.Wrap(err, "failed to create v1 handler")
	}

	httpServer := httpserver.New(
		handler.Init(),
		httpserver.SetConfig(cfg.HTTP),
	)

	log.Infof("HTTP server started at port %s", cfg.HTTP.Port)
	httpServer.Start()

	if err = <-httpServer.Notify(); !errors.Is(err, http.ErrServerClosed) {
		log.Errorf("error occurred while running HTTP server: %s", err)
	}

	err = httpServer.Shutdown()

	return errors.Wrap(err, "failed to shutdown HTTP server")
}

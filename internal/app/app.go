package app

import (
	"context"
	"net/http"

	"github.com/pkg/errors"

	"github.com/Kenplix/url-shrtnr/internal/config"
	v1 "github.com/Kenplix/url-shrtnr/internal/controller/http/v1"
	"github.com/Kenplix/url-shrtnr/internal/repository"
	"github.com/Kenplix/url-shrtnr/internal/usecase"
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
		return errors.Wrap(err, "could not create config")
	}

	log, err := logger.New(logger.SetConfig(cfg.Logger))
	if err != nil {
		return errors.Wrapf(err, "could not create logger")
	}

	hasher, err := hash.New(cfg.Hasher)
	if err != nil {
		return errors.Wrapf(err, "could not create hasher")
	}

	repo, err := repository.New(ctx, cfg.Database)
	if err != nil {
		return errors.Wrap(err, "could not create repository")
	}

	tokenService, err := auth.NewTokenService(cfg.Authorization)
	if err != nil {
		return errors.Wrap(err, "could not create token service")
	}

	manager := usecase.NewManager(usecase.Dependencies{
		Repos:        repo,
		Hasher:       hasher,
		TokenService: tokenService,
	})

	httpServer := httpserver.New(
		v1.NewHandler(manager, tokenService),
		httpserver.SetConfig(cfg.HTTP),
	)

	log.Infof("HTTP server started at port %s", cfg.HTTP.Port)
	httpServer.Start()

	if err = <-httpServer.Notify(); !errors.Is(err, http.ErrServerClosed) {
		log.Errorf("error occurred while running HTTP server: %s", err)
	}

	err = httpServer.Shutdown()

	return errors.Wrap(err, "could not shutdown HTTP server")
}

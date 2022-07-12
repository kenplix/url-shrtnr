package app

import (
	"net/http"

	"github.com/pkg/errors"

	"github.com/Kenplix/url-shrtnr/internal/config"
	v1 "github.com/Kenplix/url-shrtnr/internal/controller/http/v1"
	"github.com/Kenplix/url-shrtnr/internal/repository"
	"github.com/Kenplix/url-shrtnr/internal/usecase"
	"github.com/Kenplix/url-shrtnr/pkg/httpserver"
	"github.com/Kenplix/url-shrtnr/pkg/logger"
)

// Run -.
func Run() error {
	cfg, err := config.Read("configs")
	if err != nil {
		return errors.Wrap(err, "could not create config")
	}

	log, err := logger.New(cfg.Logger)
	if err != nil {
		return errors.Wrapf(err, "could not create logger")
	}

	repo, err := repository.New(cfg.Database)
	if err != nil {
		return errors.Wrap(err, "could not create repository")
	}

	manager := usecase.NewManager(repo)

	httpServer := httpserver.New(
		v1.NewHandler(manager),
		httpserver.SetConfig(cfg.HTTP),
	)

	log.Infof("HTTP server started at port %d", cfg.HTTP.Port)
	httpServer.Start()

	if err = <-httpServer.Notify(); !errors.Is(err, http.ErrServerClosed) {
		log.Errorf("error occurred while running HTTP server: %s", err)
	}

	err = httpServer.Shutdown()

	return errors.Wrap(err, "could not shutdown HTTP server")
}

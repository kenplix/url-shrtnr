package app

import (
	"github.com/Kenplix/url-shrtnr/internal/config"
	"github.com/Kenplix/url-shrtnr/internal/controller/http/v1"
	"github.com/Kenplix/url-shrtnr/pkg/httpserver"
	"github.com/Kenplix/url-shrtnr/pkg/logger"
	"github.com/pkg/errors"
	"net/http"
)

// Run -.
func Run() error {
	cfg, err := config.New("configs")
	if err != nil {
		return errors.Wrap(err, "could not create config")
	}

	log, err := logger.New(cfg.Logger)
	if err != nil {
		return errors.Wrapf(err, "could not create logger")
	}

	httpServer := httpserver.New(
		v1.NewHandler(log),
		httpserver.SetConfig(cfg.HTTP),
	)

	log.Info("HTTP server started at port %d", cfg.HTTP.Port)
	httpServer.Start()

	if err = <-httpServer.Notify(); !errors.Is(err, http.ErrServerClosed) {
		log.Error("error occurred while running HTTP server: %s", err)
	}

	err = httpServer.Shutdown()
	return errors.Wrap(err, "could not shutdown HTTP server")
}

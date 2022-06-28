package app

import (
	"github.com/Kenplix/url-shrtnr/internal/config"
	"github.com/Kenplix/url-shrtnr/pkg/logger"
	"github.com/pkg/errors"
)

func Run() error {
	cfg, err := config.New("configs")
	if err != nil {
		return errors.Wrap(err, "could not create config")
	}

	log, err := logger.New(cfg.Logger)
	if err != nil {
		return errors.Wrapf(err, "could not create logger")
	}

	log.Info("logger created successfully")
	return nil
}

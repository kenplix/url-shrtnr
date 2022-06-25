package app

import (
	"github.com/Kenplix/url-shrtnr/config"
	"github.com/Kenplix/url-shrtnr/pkg/logger"
	"github.com/pkg/errors"
)

func Run(cfg *config.Config) error {
	l, err := logger.New(cfg.Logger)
	if err != nil {
		return errors.Wrapf(err, "could not create logger")
	}

	l.Info("logger created successfully")
	return nil
}

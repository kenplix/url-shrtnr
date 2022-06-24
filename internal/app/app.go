package app

import (
	"github.com/Kenplix/url-shrtnr/config"
	"github.com/sirupsen/logrus"
)

func Run(cfg *config.Config) error {
	logrus.Infof("read config: %#v", cfg)
	return nil
}

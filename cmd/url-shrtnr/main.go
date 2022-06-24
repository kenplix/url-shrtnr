package main

import (
	"github.com/Kenplix/url-shrtnr/config"
	"github.com/Kenplix/url-shrtnr/internal/app"
	log "github.com/sirupsen/logrus"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("config error: %s", err)
	}

	if err := app.Run(cfg); err != nil {
		log.Fatalf("application error: %s", err)
	}
}

package main

import (
	"github.com/Kenplix/url-shrtnr/internal/app"
	log "github.com/sirupsen/logrus"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

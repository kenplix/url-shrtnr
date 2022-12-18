package main

import (
	"log"

	"github.com/kenplix/url-shrtnr/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

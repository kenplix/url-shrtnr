package main

import (
	"log"

	"github.com/Kenplix/url-shrtnr/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

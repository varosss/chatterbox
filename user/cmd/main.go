package main

import (
	"chatterbox/user/internal/infrastructure/app"
	"log"
)

func main() {
	app, err := app.New()
	if err != nil {
		log.Fatalf("failed to create app: %s", err.Error())
	}

	if err := app.Run(); err != nil {
		log.Fatalf("failed to run app: %s", err.Error())
	}
}

package main

import (
	"chatterbox/notification/internal/infrastructure/app"
	"chatterbox/notification/internal/infrastructure/config"
	"log"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %s", err.Error())
	}

	app, err := app.New(cfg)
	if err != nil {
		log.Fatalf("failed to create app: %s", err.Error())
	}

	if err := app.Run(); err != nil {
		log.Fatalf("failed to run app: %s", err.Error())
	}
}

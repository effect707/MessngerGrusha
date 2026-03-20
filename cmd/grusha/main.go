package main

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"

	"github.com/effect707/MessngerGrusha/internal/app"
	"github.com/effect707/MessngerGrusha/internal/config"
)

func main() {
	var cfg config.Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatalf("failed to read config: %v", err)
	}

	application, err := app.New(&cfg)
	if err != nil {
		log.Fatalf("failed to create application: %v", err)
	}

	if err := application.Run(); err != nil {
		log.Fatalf("application error: %v", err)
	}
}

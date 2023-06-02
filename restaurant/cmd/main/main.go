package main

import (
	"github.com/caarlos0/env"
	"github.com/comp1x/final-task/restaurant/internal/app"
	"github.com/comp1x/final-task/restaurant/internal/config"
	"log"
)

func main() {
	cfg := config.Config{}

	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("failed to retrieve env variables, %v", err)
	}

	if err := app.Run(cfg); err != nil {
		log.Fatal("error running gateway server ", err)
	}
}

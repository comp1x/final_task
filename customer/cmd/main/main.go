package main

import (
	"log"

	"github.com/caarlos0/env"

	"final-task/customer/internal/app"
	"final-task/customer/internal/config"
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

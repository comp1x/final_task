package main

import (
	"log"

	"github.com/caarlos0/env"
	"github.com/comp1x/final-task/customer/pkg/app"
	"github.com/comp1x/final-task/customer/pkg/config"
)

func main() {
	cfg := config.Config{}

	if err := env.Parse(&cfg.DB); err != nil {
		log.Fatalf("failed to retrieve env variables, %v", err)
	}
	if err := env.Parse(&cfg.Restaurant); err != nil {
		log.Fatalf("failed to retrieve env variables, %v", err)
	}
	if err := env.Parse(&cfg.Customer); err != nil {
		log.Fatalf("failed to retrieve env variables, %v", err)
	}
	if err := env.Parse(&cfg.Kafka); err != nil {
		log.Fatalf("failed to retrieve env variables, %v", err)
	}

	if err := app.Run(cfg); err != nil {
		log.Fatal("error running gateway server ", err)
	}
}

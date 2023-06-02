package main

import (
	"fmt"
	"github.com/caarlos0/env"
	"github.com/comp1x/final-task/customer/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"

	"github.com/comp1x/final-task/customer/internal/models"
)

func InitGormDB(cfg config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.PgHost, cfg.PgUser, cfg.PgPwd, cfg.PgDBName, cfg.PgPort,
	)
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func main() {
	cfg := config.Config{}

	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("failed to retrieve env variables, %v", err)
	}

	GormDB, err := InitGormDB(cfg)

	if err != nil {
		log.Fatal("? Could not connect to db", err)
	}

	GormDB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	err = GormDB.AutoMigrate(&models.Office{}, &models.User{})
	if err != nil {
		log.Fatal("problem with migration ", err)
	} else {
		fmt.Println("Migration complete")
	}
}

package config

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDb() *gorm.DB {
	var (
		dbUser, dbPass, dbName, dbHost, dbPort string
		err                                    error
	)

	if dbUser, err = GetEnv("DB_USER"); err != nil {
		log.Fatalf("DB_USER: %w", err)
	}
	if dbPass, err = GetEnv("DB_PASSWORD"); err != nil {
		log.Fatalf("DB_PASSWORD: %w", err)
	}
	if dbName, err = GetEnv("DB_NAME"); err != nil {
		log.Fatalf("DB_NAME: %w", err)
	}
	if dbHost, err = GetEnv("DB_HOST"); err != nil {
		log.Fatalf("DB_HOST: %w", err)
	}
	if dbPort, err = GetEnv("DB_PORT"); err != nil {
		log.Fatalf("DB_PORT: %w", err)
	}

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPass, dbName,
	)

	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db
}

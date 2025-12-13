package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"github.com/Candoo/thornton-pickard-api/internal/models"
)

func Initialize() (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return nil, err
	}

	// Auto-migrate models (added User model)
	if err := db.AutoMigrate(
		&models.Camera{},
		&models.Ephemera{},
		&models.Manufacturer{},
		&models.User{}, // NEW
	); err != nil {
		return nil, err
	}

	log.Println("Database connected and migrated successfully")
	return db, nil
}
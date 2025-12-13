package database

import (
	"encoding/json"
	"log"
	"os"

	"gorm.io/gorm"
	"github.com/Candoo/thornton-pickard-api/internal/models"
)

func SeedDatabase(db *gorm.DB) error {
	log.Println("Starting database seeding...")

	// Seed admin user
	if err := seedAdminUser(db); err != nil {
		return err
	}

	// Seed manufacturers
	if err := seedManufacturers(db); err != nil {
		return err
	}

	// Seed cameras
	if err := seedCameras(db); err != nil {
		return err
	}

	log.Println("Database seeding completed successfully")
	return nil
}

func seedAdminUser(db *gorm.DB) error {
	var count int64
	db.Model(&models.User{}).Count(&count)
	if count > 0 {
		log.Println("Users already exist, skipping admin user seed")
		return nil
	}

	admin := models.User{
		Email: "admin@thorntonpickard.com",
		Role:  "admin",
	}

	if err := admin.HashPassword("admin123"); err != nil {
		return err
	}

	if err := db.Create(&admin).Error; err != nil {
		return err
	}

	log.Println("✓ Admin user created (email: admin@thorntonpickard.com, password: admin123)")
	return nil
}

func seedManufacturers(db *gorm.DB) error {
	var count int64
	db.Model(&models.Manufacturer{}).Count(&count)
	if count > 0 {
		log.Println("Manufacturers already exist, skipping seed")
		return nil
	}

	manufacturers := []models.Manufacturer{
		{
			Name:        "Thornton-Pickard",
			Founded:     1888,
			Defunct:     intPtr(1939),
			Country:     "United Kingdom",
			Description: "British camera manufacturer known for quality plate cameras and shutters",
		},
	}

	if err := db.Create(&manufacturers).Error; err != nil {
		return err
	}

	log.Printf("✓ Seeded %d manufacturers", len(manufacturers))
	return nil
}

func seedCameras(db *gorm.DB) error {
	var count int64
	db.Model(&models.Camera{}).Count(&count)
	if count > 0 {
		log.Println("Cameras already exist, skipping seed")
		return nil
	}

	// Try to read from JSON file if it exists
	if data, err := os.ReadFile("seeds/cameras.json"); err == nil {
		var cameras []models.Camera
		if err := json.Unmarshal(data, &cameras); err == nil {
			if err := db.Create(&cameras).Error; err != nil {
				return err
			}
			log.Printf("✓ Seeded %d cameras from JSON file", len(cameras))
			return nil
		}
	}

	// Default seed data
	cameras := []models.Camera{
		{
			Name:             "Ruby Reflex",
			Manufacturer:     "Thornton-Pickard",
			YearIntroduced:   1909,
			YearDiscontinued: intPtr(1926),
			Format:           "Plate",
			PlateSizes:       `["4x5", "5x7"]`,
			Lens:             "Various",
			Shutter:          "Focal Plane",
			Features:         `["Reflex viewing", "Tilting back", "Rising front"]`,
			Description:      "Professional reflex camera popular with press photographers",
			Rarity:           "Uncommon",
			EstimatedValueMin: float64Ptr(500),
			EstimatedValueMax: float64Ptr(800),
		},
		{
			Name:             "Imperial Triple Extension",
			Manufacturer:     "Thornton-Pickard",
			YearIntroduced:   1895,
			Format:           "Plate",
			PlateSizes:       `["Half-plate", "Whole-plate"]`,
			Shutter:          "Time Shutter",
			Features:         `["Triple extension bellows", "Mahogany construction"]`,
			Description:      "High-quality field camera with extensive movements",
			Rarity:           "Rare",
			EstimatedValueMin: float64Ptr(300),
			EstimatedValueMax: float64Ptr(600),
		},
	}

	if err := db.Create(&cameras).Error; err != nil {
		return err
	}

	log.Printf("✓ Seeded %d cameras", len(cameras))
	return nil
}

func intPtr(i int) *int {
	return &i
}

func float64Ptr(f float64) *float64 {
	return &f
}
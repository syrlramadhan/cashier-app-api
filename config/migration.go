package config

import (
	"log"

	"github.com/syrlramadhan/cashier-app/models"
	"golang.org/x/crypto/bcrypt"
)

func RunMigration() {
	err := DB.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.Product{},
		&models.Transaction{},
		&models.TransactionItem{},
		&models.Setting{},
	)

	if err != nil {
		log.Fatal("Migration failed:", err)
	}

	log.Println("Database migration completed successfully")

	// Seed default data
	seedDefaultData()
}

func seedDefaultData() {
	// Seed default admin user
	var userCount int64
	DB.Model(&models.User{}).Count(&userCount)
	if userCount == 0 {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		admin := models.User{
			Name:     "Admin",
			Email:    "admin@kasir.com",
			Password: string(hashedPassword),
			Role:     "admin",
			IsActive: true,
		}
		DB.Create(&admin)

		log.Println("Default users seeded")
	}

	// Seed default categories
	var categoryCount int64
	DB.Model(&models.Category{}).Count(&categoryCount)
	if categoryCount == 0 {
		categories := []models.Category{
			{Name: "Makanan"},
			{Name: "Minuman"},
			{Name: "Snack"},
		}
		DB.Create(&categories)
		log.Println("Default categories seeded")
	}
}

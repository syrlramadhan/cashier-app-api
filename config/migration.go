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

	// Seed default products (menu default)
	var productCount int64
	DB.Model(&models.Product{}).Count(&productCount)
	if productCount == 0 {
		// Get Minuman category ID
		var minumanCategory models.Category
		DB.Where("name = ?", "Minuman").First(&minumanCategory)

		defaultProducts := []models.Product{
			{
				Name:       "Avo Coffee",
				Price:      25000,
				Stock:      100,
				CategoryID: minumanCategory.ID,
				Image:      "/uploads/products/avocoffe.png",
			},
			{
				Name:       "Cappucino",
				Price:      22000,
				Stock:      100,
				CategoryID: minumanCategory.ID,
				Image:      "/uploads/products/cappucino.png",
			},
			{
				Name:       "Chococa",
				Price:      20000,
				Stock:      100,
				CategoryID: minumanCategory.ID,
				Image:      "/uploads/products/chococa.png",
			},
			{
				Name:       "Green Tea",
				Price:      18000,
				Stock:      100,
				CategoryID: minumanCategory.ID,
				Image:      "/uploads/products/greentea.png",
			},
			{
				Name:       "Macachino",
				Price:      24000,
				Stock:      100,
				CategoryID: minumanCategory.ID,
				Image:      "/uploads/products/macachino.png",
			},
			{
				Name:       "Machiato",
				Price:      26000,
				Stock:      100,
				CategoryID: minumanCategory.ID,
				Image:      "/uploads/products/machiato.png",
			},
			{
				Name:       "Vallate Coffee Milk",
				Price:      23000,
				Stock:      100,
				CategoryID: minumanCategory.ID,
				Image:      "/uploads/products/vallate_coffemilk.png",
			},
		}

		DB.Create(&defaultProducts)
		log.Println("Default products seeded")
	}

	// Seed default settings
	var settingCount int64
	DB.Model(&models.Setting{}).Count(&settingCount)
	if settingCount == 0 {
		defaultSettings := []models.Setting{
			// Store settings
			{Key: "store_name", Value: "Kasir POS"},
			{Key: "store_address", Value: "Jl. Contoh No. 123, Jakarta"},
			{Key: "store_phone", Value: "021-12345678"},
			{Key: "store_logo", Value: ""},
			// Payment settings
			{Key: "tax_rate", Value: "11"},
			{Key: "currency", Value: "IDR"},
			{Key: "payment_cash_enabled", Value: "true"},
			{Key: "payment_card_enabled", Value: "true"},
			{Key: "payment_qris_enabled", Value: "true"},
			// Printer settings
			{Key: "printer_type", Value: "thermal"},
			{Key: "receipt_footer", Value: "Terima kasih atas kunjungan Anda!"},
			{Key: "auto_print", Value: "true"},
			{Key: "print_logo", Value: "true"},
			{Key: "print_duplicate", Value: "false"},
			// General settings
			{Key: "enable_sound", Value: "true"},
			{Key: "enable_notifications", Value: "true"},
			{Key: "auto_logout", Value: "30"},
		}
		DB.Create(&defaultSettings)
		log.Println("Default settings seeded")
	}
}

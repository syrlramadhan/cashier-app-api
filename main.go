package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/syrlramadhan/cashier-app/config"
	"github.com/syrlramadhan/cashier-app/controllers"
	"github.com/syrlramadhan/cashier-app/repositories"
	"github.com/syrlramadhan/cashier-app/routes"
	"github.com/syrlramadhan/cashier-app/services"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	// Initialize database connection
	config.ConnectDatabase()
	db := config.GetDB()

	// Run migrations
	config.RunMigration()

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	categoryRepo := repositories.NewCategoryRepository(db)
	productRepo := repositories.NewProductRepository(db)
	transactionRepo := repositories.NewTransactionRepository(db)
	transactionItemRepo := repositories.NewTransactionItemRepository(db)
	settingRepo := repositories.NewSettingRepository(db)

	// Initialize services
	userService := services.NewUserService(userRepo)
	categoryService := services.NewCategoryService(categoryRepo)
	productService := services.NewProductService(productRepo, categoryRepo)
	transactionService := services.NewTransactionService(transactionRepo, transactionItemRepo, productRepo)
	settingService := services.NewSettingService(settingRepo)
	reportService := services.NewReportService(transactionRepo, transactionItemRepo, productRepo, categoryRepo)

	// Initialize controllers
	userController := controllers.NewUserController(userService)
	categoryController := controllers.NewCategoryController(categoryService)
	productController := controllers.NewProductController(productService)
	transactionController := controllers.NewTransactionController(transactionService)
	settingController := controllers.NewSettingController(settingService)
	reportController := controllers.NewReportController(reportService)

	// Initialize routes
	r := routes.NewRoutes(
		userController,
		categoryController,
		productController,
		transactionController,
		settingController,
		reportController,
	)

	// Setup router
	router := r.SetupRouter()

	// Get port from environment or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Cashier API Server starting on port %s", port)
	log.Printf("📝 API Documentation: http://localhost:%s/api/v1", port)
	log.Printf("❤️  Health Check: http://localhost:%s/health", port)

	// Start server
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

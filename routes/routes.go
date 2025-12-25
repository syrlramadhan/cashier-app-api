package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/syrlramadhan/cashier-app/controllers"
	"github.com/syrlramadhan/cashier-app/middleware"
)

type Routes struct {
	userController        *controllers.UserController
	categoryController    *controllers.CategoryController
	productController     *controllers.ProductController
	transactionController *controllers.TransactionController
	settingController     *controllers.SettingController
	reportController      *controllers.ReportController
}

func NewRoutes(
	userController *controllers.UserController,
	categoryController *controllers.CategoryController,
	productController *controllers.ProductController,
	transactionController *controllers.TransactionController,
	settingController *controllers.SettingController,
	reportController *controllers.ReportController,
) *Routes {
	return &Routes{
		userController:        userController,
		categoryController:    categoryController,
		productController:     productController,
		transactionController: transactionController,
		settingController:     settingController,
		reportController:      reportController,
	}
}

func (r *Routes) SetupRouter() *gin.Engine {
	router := gin.Default()

	// Global middleware
	router.Use(middleware.CORSMiddleware())

	// Health check endpoint
	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"status":  "ok",
			"message": "Cashier API is running",
		})
	})

	// API v1 group
	api := router.Group("/api/v1")
	{
		// Auth routes (public)
		auth := api.Group("/auth")
		{
			auth.POST("/login", r.userController.Login)
			auth.POST("/register", r.userController.Register)
		}

		// Protected routes
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			// User routes
			users := protected.Group("/users")
			{
				users.GET("/profile", r.userController.GetProfile)
				users.GET("", middleware.AdminOnly(), r.userController.GetAllUsers)
				users.GET("/:id", middleware.AdminOnly(), r.userController.GetUserByID)
				users.PUT("/:id", middleware.AdminOnly(), r.userController.UpdateUser)
				users.DELETE("/:id", middleware.AdminOnly(), r.userController.DeleteUser)
			}

			// Category routes
			categories := protected.Group("/categories")
			{
				categories.GET("", r.categoryController.GetAllCategories)
				categories.GET("/:id", r.categoryController.GetCategoryByID)
				categories.POST("", middleware.ManagerOrAdmin(), r.categoryController.CreateCategory)
				categories.PUT("/:id", middleware.ManagerOrAdmin(), r.categoryController.UpdateCategory)
				categories.DELETE("/:id", middleware.AdminOnly(), r.categoryController.DeleteCategory)
			}

			// Product routes
			products := protected.Group("/products")
			{
				products.GET("", r.productController.GetAllProducts)
				products.GET("/:id", r.productController.GetProductByID)
				products.GET("/category/:category_id", r.productController.GetProductsByCategory)
				products.POST("", middleware.ManagerOrAdmin(), r.productController.CreateProduct)
				products.PUT("/:id", middleware.ManagerOrAdmin(), r.productController.UpdateProduct)
				products.PATCH("/:id/stock", middleware.ManagerOrAdmin(), r.productController.UpdateStock)
				products.DELETE("/:id", middleware.AdminOnly(), r.productController.DeleteProduct)
			}

			// Transaction routes
			transactions := protected.Group("/transactions")
			{
				transactions.GET("", r.transactionController.GetAllTransactions)
				transactions.GET("/today", r.transactionController.GetTodayTransactions)
				transactions.GET("/:id", r.transactionController.GetTransactionByID)
				transactions.GET("/code/:code", r.transactionController.GetTransactionByCode)
				transactions.GET("/user/:user_id", r.transactionController.GetTransactionsByUser)
				transactions.POST("", r.transactionController.CreateTransaction)
				transactions.POST("/:id/cancel", middleware.ManagerOrAdmin(), r.transactionController.CancelTransaction)
			}

			// Setting routes
			settings := protected.Group("/settings")
			{
				settings.GET("", r.settingController.GetAllSettings)
				settings.GET("/store", r.settingController.GetStoreSettings)
				settings.GET("/payment", r.settingController.GetPaymentSettings)
				settings.GET("/:key", r.settingController.GetSettingByKey)
				settings.PUT("", middleware.AdminOnly(), r.settingController.UpdateSetting)
				settings.PUT("/batch", middleware.AdminOnly(), r.settingController.UpdateSettings)
			}

			// Report routes
			reports := protected.Group("/reports")
			{
				reports.GET("/dashboard", r.reportController.GetDashboard)
				reports.GET("/revenue/daily", r.reportController.GetDailyRevenue)
				reports.GET("/revenue/range", r.reportController.GetRevenueByDateRange)
				reports.GET("/payment-distribution", r.reportController.GetPaymentDistribution)
				reports.GET("/products/top", r.reportController.GetTopProducts)
				reports.GET("/summary/monthly", r.reportController.GetMonthlySummary)
				reports.GET("/export/transactions", middleware.ManagerOrAdmin(), r.reportController.ExportTransactions)
			}
		}
	}

	return router
}

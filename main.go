package main

import (
	"doAnHTTT_go/config"
	"doAnHTTT_go/controllers"
	"doAnHTTT_go/middlewares"
	"doAnHTTT_go/services"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	config.ConnectDatabase()
	go services.CleanupUnverifiedUsersAndOTP()

	router := gin.Default()

	router.Static("/uploads", "./uploads")

	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "API is running",
		})
	})

	router.HEAD("/", func(c *gin.Context) {
		c.Status(200)
	})

	healthHandler := func(c *gin.Context) {
		// sqlDB, err := config.DB.DB()
		// dbStatus := "connected"

		// if err != nil || sqlDB.Ping() != nil {
		// 	dbStatus = "disconnected"
		// }

		c.JSON(200, gin.H{
			"status": "ok",
			// "database": dbStatus,
			"message": "Hệ thống đang hoạt động",
		})
	}

	router.GET("/health", healthHandler)
	router.HEAD("/health", healthHandler)

	// ================= ROUTES =================

	authRoutes := router.Group("/api/auth")
	{
		authRoutes.POST("/register", controllers.RegisterOwnerHandler)
		authRoutes.POST("/verify-otp", controllers.VerifyRegistrationOTPHandler)
		authRoutes.POST("/login", controllers.LoginHandler)
		authRoutes.POST("/request-password-reset", controllers.RequestPasswordResetHandler)
		authRoutes.POST("/confirm-password", controllers.ConfirmNewPasswordHandler)
		authRoutes.POST("/refresh", controllers.RefreshTokenHandler)
	}

	ownerRoutes := router.Group("/api/admin")
	ownerRoutes.Use(middlewares.RequireRole("OWNER"))
	{
		ownerRoutes.POST("/houses", controllers.CreateHouseHandler)
		ownerRoutes.PUT("/houses/:id", controllers.UpdateHouseHandler)
		ownerRoutes.DELETE("/houses/:id", controllers.DeleteHouseHandler)

		ownerRoutes.POST("/rooms/:id/services", controllers.AssignServicesToRoomHandler)

		ownerRoutes.POST("/rooms", controllers.CreateRoomHandler)
		ownerRoutes.PUT("/rooms/:id", controllers.UpdateRoomHandler)
		ownerRoutes.DELETE("/rooms/:id", controllers.DeleteRoomHandler)

		ownerRoutes.POST("/services", controllers.CreateServiceHandler)
		ownerRoutes.PUT("/services/:id", controllers.UpdateServiceHandler)
		ownerRoutes.DELETE("/services/:id", controllers.DeleteServiceHandler)

		ownerRoutes.GET("/users", controllers.GetAllUserHandler)
		ownerRoutes.POST("/users", controllers.CreateUserHandler)
		ownerRoutes.PUT("/users/:id", controllers.UpdateUserHandler)
		ownerRoutes.DELETE("/users/:id", controllers.DeleteUserHandler)

		ownerRoutes.DELETE("/tenants/:id", controllers.DeleteTenantHandler)

		ownerRoutes.PUT("/transactions/:id", controllers.UpdateTransactionHandler)
		ownerRoutes.DELETE("/transactions/:id", controllers.DeleteTransactionHandler)

		ownerRoutes.GET("/reports/profit-loss", controllers.GetProfitLossHandler)
	}

	generalRoutes := router.Group("/api/general")
	generalRoutes.Use(middlewares.RequireRole("OWNER", "STAFF"))
	{
		generalRoutes.GET("/houses", controllers.GetAllHousesHandler)
		generalRoutes.GET("/rooms", controllers.GetAllRoomHandler)
		generalRoutes.GET("/services", controllers.GetAllServiceHandler)

		generalRoutes.GET("/tenants", controllers.GetAllTenantHandler)
		generalRoutes.POST("/tenants", controllers.CreateTenantHandler)
		generalRoutes.PUT("/tenants/:id", controllers.UpdateTenantHandler)

		generalRoutes.GET("/contracts", controllers.GetAllContractHandler)
		generalRoutes.POST("/contracts", controllers.CreateContractHandler)
		generalRoutes.PUT("/contracts/:id/terminate", controllers.TerminateContractHandler)

		generalRoutes.GET("/rooms/:id/services", controllers.GetServicesOfRoomHandler)

		generalRoutes.GET("/tasks", controllers.GetAllTaskHandler)
		generalRoutes.POST("/tasks", controllers.CreateTaskHandler)
		generalRoutes.PUT("/tasks/:id", controllers.UpdateTaskHandler)
		generalRoutes.DELETE("/tasks/:id", controllers.DeleteTaskHandler)

		generalRoutes.POST("/tasks/upload-image", controllers.UploadTaskImageHandler)

		generalRoutes.POST("/meter-readings", controllers.AddMeterReadingHandler)
		generalRoutes.GET("/meter-readings", controllers.GetMeterReadingsHandler)
		generalRoutes.PUT("/meter-readings/:id", controllers.UpdateMeterReadingHandler)
		generalRoutes.DELETE("/meter-readings/:id", controllers.DeleteMeterReadingHandler)
		generalRoutes.GET("/meter-readings/latest-index", controllers.GetLatestIndexHandler)

		generalRoutes.GET("/invoices", controllers.GetAllInvoicesHandler)
		generalRoutes.POST("/invoices/generate", controllers.TriggerGenerateInvoices)
		generalRoutes.POST("/invoices/pay", controllers.PayInvoiceHandler)
		generalRoutes.DELETE("/invoices/:id", controllers.DeleteInvoiceHandler)

		generalRoutes.GET("/transactions", controllers.GetAllTransactionsHandler)
		generalRoutes.POST("/transactions", controllers.AddTransactionHandler)

		generalRoutes.GET("/profile/me", controllers.GetMyProfileHandler)
		generalRoutes.PUT("/profile/me", controllers.UpdateMyProfileHandler)
		generalRoutes.PUT("/profile/password", controllers.ChangeMyPasswordHandler)
		generalRoutes.POST("/upload", controllers.UploadImageHandler)

		generalRoutes.POST("/upload-multiple", controllers.UploadMultipleImagesHandler)
		generalRoutes.POST("/delete-file", controllers.DeleteImageHandler)

	}

	webhookRoutes := router.Group("/api/webhooks")
	{
		webhookRoutes.POST("/bank-transfer", controllers.WebhookBankTransferHandler)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server đang chạy tại cổng %s", port)
	log.Fatal(router.Run(":" + port))
}

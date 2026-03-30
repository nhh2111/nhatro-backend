package main

import (
	"doAnHTTT_go/config"
	"doAnHTTT_go/controllers"
	"doAnHTTT_go/middlewares"
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
	router := gin.Default()

	// 1. CẤU HÌNH CORS CHO PHÉP MỌI DOMAIN (Cần thiết khi đưa lên Vercel)
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true, // Thay đổi ở đây: Mở cửa cho Vercel gọi API
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// NHÓM 1: XÁC THỰC (Không cần Token)
	authRoutes := router.Group("/api/auth")
	{
		authRoutes.POST("/register", controllers.RegisterOwnerHandler)
		authRoutes.POST("/verify-otp", controllers.VerifyRegistrationOTPHandler)
		authRoutes.POST("/login", controllers.LoginHandler)
		authRoutes.POST("/request-password-reset", controllers.RequestPasswordResetHandler)
		authRoutes.POST("/confirm-password", controllers.ConfirmNewPasswordHandler)
		authRoutes.POST("/refresh", controllers.RefreshTokenHandler)
	}

	// NHÓM 2: QUYỀN CỦA CHỦ TRỌ (OWNER)
	ownerRoutes := router.Group("/api/admin")
	ownerRoutes.Use(middlewares.RequireRole("OWNER"))
	{
		// ... (Toàn bộ các route của Admin giữ nguyên) ...
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

	// NHÓM 3: QUYỀN CHUNG (OWNER & STAFF)
	generalRoutes := router.Group("/api/general")
	generalRoutes.Use(middlewares.RequireRole("OWNER", "STAFF"))
	{
		// ... (Toàn bộ các route của General giữ nguyên) ...
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

		generalRoutes.POST("/meter-readings", controllers.AddMeterReadingHandler)
		generalRoutes.GET("/meter-readings", controllers.GetMeterReadingsHandler)
		generalRoutes.PUT("/meter-readings/:id", controllers.UpdateMeterReadingHandler)
		generalRoutes.DELETE("/meter-readings/:id", controllers.DeleteMeterReadingHandler)

		generalRoutes.GET("/invoices", controllers.GetAllInvoicesHandler)
		generalRoutes.POST("/invoices/generate", controllers.TriggerGenerateInvoices)
		generalRoutes.POST("/invoices/pay", controllers.PayInvoiceHandler)
		generalRoutes.DELETE("/invoices/:id", controllers.DeleteInvoiceHandler)

		generalRoutes.GET("/transactions", controllers.GetAllTransactionsHandler)
		generalRoutes.POST("/transactions", controllers.AddTransactionHandler)

		generalRoutes.GET("/profile/me", controllers.GetMyProfileHandler)
		generalRoutes.PUT("/profile/me", controllers.UpdateMyProfileHandler)
	}

	// 2. CHẠY SERVER BẰNG PORT ĐỘNG (Xóa đoạn cứng 8080 đi)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Dự phòng chạy localhost
	}

	log.Printf("Server đang chạy tại cổng %s", port)
	errRun := router.Run(":" + port)

	if errRun != nil {
		log.Fatalf("Lỗi khi khởi chạy server: %v", errRun)
	}
}

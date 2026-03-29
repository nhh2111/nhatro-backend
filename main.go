package main

import (
	"doAnHTTT_go/config"
	"doAnHTTT_go/controllers"
	"doAnHTTT_go/middlewares"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	config.ConnectDatabase()
	router := gin.Default()

	// CẤU HÌNH CORS CHO PHÉP ANGULAR (Cổng 4200) KẾT NỐI ĐẾN GOLANG (Cổng 8080)
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:4200"}, // Đổi thành domain thực tế sau này
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// =====================================================================
	// NHÓM 1: XÁC THỰC (Không cần Token)
	// =====================================================================
	authRoutes := router.Group("/api/auth")
	{
		authRoutes.POST("/register", controllers.RegisterOwnerHandler)
		authRoutes.POST("/verify-otp", controllers.VerifyRegistrationOTPHandler)
		authRoutes.POST("/login", controllers.LoginHandler)
		authRoutes.POST("/request-password-reset", controllers.RequestPasswordResetHandler)
		authRoutes.POST("/confirm-password", controllers.ConfirmNewPasswordHandler)
		authRoutes.POST("/refresh", controllers.RefreshTokenHandler)
	}

	// =====================================================================
	// NHÓM 2: QUYỀN CỦA CHỦ TRỌ (OWNER)
	// (Đảm nhận các thao tác Xóa nhạy cảm và Thêm/Sửa Cấu hình hệ thống)
	// =====================================================================
	ownerRoutes := router.Group("/api/admin")
	ownerRoutes.Use(middlewares.RequireRole("OWNER"))
	{
		// 1. Quản lý Nhà, Phòng, Dịch vụ (Thêm, Sửa, Xóa)
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

		// 2. Tài khoản nhân viên (OWNER được Thêm, Sửa, Xóa)
		ownerRoutes.GET("/users", controllers.GetAllUserHandler)
		ownerRoutes.POST("/users", controllers.CreateUserHandler)
		ownerRoutes.PUT("/users/:id", controllers.UpdateUserHandler)
		ownerRoutes.DELETE("/users/:id", controllers.DeleteUserHandler)

		// 3. Khách thuê (Chỉ OWNER được Xóa)
		ownerRoutes.DELETE("/tenants/:id", controllers.DeleteTenantHandler)

		// 4. Phiếu Thu/Chi (Chỉ OWNER được Sửa, Xóa)
		ownerRoutes.PUT("/transactions/:id", controllers.UpdateTransactionHandler)
		ownerRoutes.DELETE("/transactions/:id", controllers.DeleteTransactionHandler)

		// 5. Báo cáo Lời Lỗ (Chỉ OWNER được xem)
		ownerRoutes.GET("/reports/profit-loss", controllers.GetProfitLossHandler)
	}

	// =====================================================================
	// NHÓM 3: QUYỀN CHUNG (OWNER & STAFF)
	// (Vận hành hàng ngày, nhập liệu cơ bản)
	// =====================================================================
	generalRoutes := router.Group("/api/general")
	generalRoutes.Use(middlewares.RequireRole("OWNER", "STAFF"))
	{
		// 1. Xem danh sách Nhà, Phòng, Dịch vụ (STAFF chỉ được Xem)
		generalRoutes.GET("/houses", controllers.GetAllHousesHandler)
		generalRoutes.GET("/rooms", controllers.GetAllRoomHandler)
		generalRoutes.GET("/services", controllers.GetAllServiceHandler)

		// 2. Khách & Hợp đồng (STAFF được Thêm, Sửa)
		generalRoutes.GET("/tenants", controllers.GetAllTenantHandler)
		generalRoutes.POST("/tenants", controllers.CreateTenantHandler)
		generalRoutes.PUT("/tenants/:id", controllers.UpdateTenantHandler)

		generalRoutes.GET("/contracts", controllers.GetAllContractHandler)
		generalRoutes.POST("/contracts", controllers.CreateContractHandler)
		generalRoutes.PUT("/contracts/:id/terminate", controllers.TerminateContractHandler)

		generalRoutes.GET("/rooms/:id/services", controllers.GetServicesOfRoomHandler)

		// 3. Công việc/Sự cố (Ai cũng được cập nhật)
		generalRoutes.GET("/tasks", controllers.GetAllTaskHandler)
		generalRoutes.POST("/tasks", controllers.CreateTaskHandler)
		generalRoutes.PUT("/tasks/:id", controllers.UpdateTaskHandler)
		generalRoutes.DELETE("/tasks/:id", controllers.DeleteTaskHandler)

		// 4. Chỉ số Điện/Nước (STAFF được Thêm, Sửa - Logic khóa tháng cũ nằm ở Service)
		generalRoutes.POST("/meter-readings", controllers.AddMeterReadingHandler)
		generalRoutes.GET("/meter-readings", controllers.GetMeterReadingsHandler)
		generalRoutes.PUT("/meter-readings/:id", controllers.UpdateMeterReadingHandler)
		generalRoutes.DELETE("/meter-readings/:id", controllers.DeleteMeterReadingHandler)

		// 5. Hóa Đơn & Thanh toán (STAFF được Khởi tạo và Thu tiền)
		generalRoutes.GET("/invoices", controllers.GetAllInvoicesHandler)
		generalRoutes.POST("/invoices/generate", controllers.TriggerGenerateInvoices)
		generalRoutes.POST("/invoices/pay", controllers.PayInvoiceHandler)
		generalRoutes.DELETE("/invoices/:id", controllers.DeleteInvoiceHandler)

		// 6. Phiếu Thu/Chi (STAFF chỉ được Thêm)
		generalRoutes.GET("/transactions", controllers.GetAllTransactionsHandler)
		generalRoutes.POST("/transactions", controllers.AddTransactionHandler)

		// =====================================================================
		// 7. Profile (Thông tin cá nhân)
		// =====================================================================
		generalRoutes.GET("/profile/me", controllers.GetMyProfileHandler)
		generalRoutes.PUT("/profile/me", controllers.UpdateMyProfileHandler)
	}

	log.Println("Server đang chạy tại http://localhost:8080")
	errRun := router.Run(":8080")

	if errRun != nil {
		log.Fatalf("Lỗi khi khởi chạy server: %v", errRun)
	}
}

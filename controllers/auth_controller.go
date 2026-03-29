package controllers

import (
	"doAnHTTT_go/config"
	"doAnHTTT_go/models"
	"doAnHTTT_go/services"
	"doAnHTTT_go/utils"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type RegisterReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	FullName string `json:"full_name" binding:"required"`
}

type OTPReq struct {
	Email   string `json:"email" binding:"required,email"`
	OTPCode string `json:"otp_code" binding:"required"`
}

type LoginReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
type RefreshReq struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func RefreshTokenHandler(c *gin.Context) {
	var req RefreshReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Thiếu refresh token"})
		return
	}

	// Giải mã Refresh Token
	secretKey := []byte(os.Getenv("JWT_SECRET_KEY"))
	token, err := jwt.Parse(req.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token không hợp lệ hoặc đã hết hạn"})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi đọc dữ liệu token"})
		return
	}

	userID := uint(claims["user_id"].(float64))

	// Truy vấn DB để lấy lại Role của User (Vì refresh token của bạn không lưu role)
	var user models.User // Đảm bảo bạn gọi đúng model tài khoản của bạn (User hoặc Owner)
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Không tìm thấy người dùng"})
		return
	}

	// Gọi hàm GenerateTokens của bạn để cấp cặp khóa mới
	newAccessToken, newRefreshToken, errToken := utils.GenerateTokens(userID, user.Role)
	if errToken != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi tạo token mới"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken,
	})
}

func RegisterOwnerHandler(c *gin.Context) {
	var req RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ"})
		return
	}

	errService := services.RegisterOwner(req.Email, req.Password, req.FullName)
	if errService != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errService.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Đăng ký thành công, vui lòng kiểm tra email lấy mã OTP"})
}

func VerifyRegistrationOTPHandler(c *gin.Context) {
	var req OTPReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ"})
		return
	}

	errService := services.VerifyRegistrationOTP(req.Email, req.OTPCode)
	if errService != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errService.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Kích hoạt tài khoản thành công! Bạn có thể đăng nhập."})
}

func LoginHandler(c *gin.Context) {
	var req LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Vui lòng nhập email và mật khẩu"})
		return
	}

	accessToken, refreshToken, isFirstLogin, errService := services.LoginUser(req.Email, req.Password)

	if isFirstLogin {
		c.JSON(http.StatusForbidden, gin.H{
			"error":          errService.Error(),
			"is_first_login": true,
			"email":          req.Email,
		})
		return
	}

	if errService != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": errService.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Đăng nhập thành công",
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// Yêu cầu gửi OTP đổi mật khẩu
func RequestPasswordResetHandler(c *gin.Context) {
	var request struct {
		Email string `json:"email"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ"})
		return
	}

	if err := services.RequestPasswordChangeOTP(request.Email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Mã OTP đã được gửi về email"})
}

// Xác nhận OTP và đặt mật khẩu mới
func ConfirmNewPasswordHandler(c *gin.Context) {
	var request struct {
		Email       string `json:"email"`
		OTPCode     string `json:"otp_code"`
		NewPassword string `json:"new_password"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ"})
		return
	}

	if err := services.ConfirmNewPassword(request.Email, request.OTPCode, request.NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Đổi mật khẩu thành công"})
}

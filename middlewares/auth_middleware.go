package middlewares

import (
	"doAnHTTT_go/config"
	"doAnHTTT_go/models"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		authHeader := ginContext.GetHeader("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			ginContext.JSON(http.StatusUnauthorized, gin.H{"error": "Không tìm thấy token xác thực"})
			ginContext.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		secretKey := []byte(os.Getenv("JWT_SECRET_KEY"))
		token, errParse := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})

		if errParse != nil || !token.Valid {
			ginContext.JSON(http.StatusUnauthorized, gin.H{"error": "Token không hợp lệ hoặc đã hết hạn"})
			ginContext.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			ginContext.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi đọc dữ liệu token"})
			ginContext.Abort()
			return
		}

		userRole := claims["role"].(string)
		userID := uint(claims["user_id"].(float64))

		if !checkRoleExists(userRole, allowedRoles) {
			ginContext.JSON(http.StatusForbidden, gin.H{"error": "Bạn không có quyền thực hiện thao tác này"})
			ginContext.Abort()
			return
		}

		var ownerID uint
		if userRole == "STAFF" {
			var staff models.User
			if err := config.DB.Select("employer_id").First(&staff, userID).Error; err == nil {
				ownerID = uint(staff.EmployerID)
			}
		} else {
			ownerID = userID
		}

		ginContext.Set("userID", userID)
		ginContext.Set("ownerID", ownerID)
		ginContext.Set("userRole", userRole)
		ginContext.Next()
	}
}

func checkRoleExists(role string, allowedRoles []string) bool {
	for _, allowedRole := range allowedRoles {
		if role == allowedRole {
			return true
		}
	}
	return false
}

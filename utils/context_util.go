package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetOwnerIDFromContext(c *gin.Context) (uint, bool) {
	value, exists := c.Get("ownerID")
	if !exists {
		return 0, false
	}

	switch v := value.(type) {
	case uint:
		return v, true
	case int:
		if v < 0 {
			return 0, false
		}
		return uint(v), true
	case int64:
		if v < 0 {
			return 0, false
		}
		return uint(v), true
	case float64:
		if v < 0 {
			return 0, false
		}
		return uint(v), true
	default:
		return 0, false
	}
}

func RequireOwnerID(c *gin.Context) (uint, bool) {
	ownerID, ok := GetOwnerIDFromContext(c)
	if !ok {
		ErrorResponse(c, http.StatusUnauthorized, 401, "Không xác định được danh tính người dùng")
		return 0, false
	}
	return ownerID, true
}

func GetUserIDFromContext(c *gin.Context) (uint, bool) {
	value, exists := c.Get("userID")
	if !exists {
		return 0, false
	}
	return CoerceUint(value)
}

func RequireUserID(c *gin.Context) (uint, bool) {
	userID, ok := GetUserIDFromContext(c)
	if !ok || userID == 0 {
		ErrorResponse(c, http.StatusUnauthorized, 401, "Không xác thực được danh tính")
		return 0, false
	}
	return userID, true
}

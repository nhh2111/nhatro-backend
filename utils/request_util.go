package utils

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ParseUintParam(c *gin.Context, paramName string) (uint, error) {
	raw := c.Param(paramName)
	if raw == "" {
		return 0, errors.New("thiếu tham số " + paramName)
	}

	id64, err := strconv.ParseUint(raw, 10, 64)
	if err != nil || id64 == 0 {
		return 0, errors.New("tham số " + paramName + " không hợp lệ")
	}
	return uint(id64), nil
}

func GetUserRoleFromContext(c *gin.Context) (string, bool) {
	value, exists := c.Get("userRole")
	if !exists {
		return "", false
	}
	role, ok := value.(string)
	if !ok || role == "" {
		return "", false
	}
	return role, true
}

func RequireUserRole(c *gin.Context) (string, bool) {
	role, ok := GetUserRoleFromContext(c)
	if !ok {
		ErrorResponse(c, http.StatusUnauthorized, 401, "Không xác định được quyền người dùng")
		return "", false
	}
	return role, true
}

package controllers

import (
	"doAnHTTT_go/services"
	"doAnHTTT_go/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetAllUserHandler(ginContext *gin.Context) {
	page, pageSize := utils.GetPaginationParams(ginContext)
	search := ginContext.Query("search")

	resultData, err := services.GetAllStaffs(page, pageSize, search)

	if err != nil {
		utils.ErrorResponse(ginContext, http.StatusInternalServerError, 500, "Lỗi lấy danh sách nhân viên: "+err.Error())
		return
	}

	utils.SuccessResponse(ginContext, http.StatusOK, resultData)
}

func CreateUserHandler(ginContext *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		FullName string `json:"full_name" binding:"required"`
	}

	if errBind := ginContext.ShouldBindJSON(&req); errBind != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "Vui lòng nhập email và họ tên hợp lệ"})
		return
	}

	errCreate := services.CreateStaffAccount(req.Email, req.FullName)
	if errCreate != nil {
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": errCreate.Error()})
		return
	}

	ginContext.JSON(http.StatusCreated, gin.H{"message": "Tạo nhân viên mới thành công. Mật khẩu mặc định là 123456"})
}

func UpdateUserHandler(ginContext *gin.Context) {
	userID, errParse := strconv.Atoi(ginContext.Param("id"))
	if errParse != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "ID nhân viên không hợp lệ"})
		return
	}

	var updateData map[string]interface{}
	if errBind := ginContext.ShouldBindJSON(&updateData); errBind != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu cập nhật không hợp lệ"})
		return
	}

	errService := services.UpdateStaffs(uint(userID), updateData)
	if errService != nil {
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": errService.Error()})
		return
	}

	ginContext.JSON(http.StatusOK, gin.H{"message": "Cập nhật nhân viên thành công"})
}

func DeleteUserHandler(ginContext *gin.Context) {
	userID, errParse := strconv.Atoi(ginContext.Param("id"))
	if errParse != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "ID nhân viên không hợp lệ"})
		return
	}

	errService := services.DeleteUser(uint(userID))
	if errService != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": errService.Error()})
		return
	}

	ginContext.JSON(http.StatusOK, gin.H{"message": "Xóa nhân viên thành công"})
}

func GetMyProfileHandler(c *gin.Context) {
	userIDValue, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, 401, "Không xác thực được danh tính")
		return
	}

	var userID uint
	switch v := userIDValue.(type) {
	case float64:
		userID = uint(v)
	case uint:
		userID = v
	default:
		utils.ErrorResponse(c, http.StatusInternalServerError, 500, "Lỗi định dạng ID người dùng")
		return
	}

	resultData, err := services.GetMyProfile(userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, 404, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, resultData)
}

func UpdateMyProfileHandler(c *gin.Context) {
	userIDValue, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, 401, "Không xác thực được danh tính")
		return
	}

	var userID uint
	switch v := userIDValue.(type) {
	case float64:
		userID = uint(v)
	case uint:
		userID = v
	}

	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, 400, "Dữ liệu cập nhật không hợp lệ")
		return
	}

	if err := services.UpdateMyProfile(userID, updateData); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, 500, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"errorCode":    200,
		"errorMessage": "Cập nhật thông tin cá nhân thành công",
		"result":       nil,
	})
}

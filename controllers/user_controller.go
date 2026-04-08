package controllers

import (
	"doAnHTTT_go/services"
	"doAnHTTT_go/utils"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

func GetAllUserHandler(ginContext *gin.Context) {
	page, pageSize := utils.GetPaginationParams(ginContext)
	search := ginContext.Query("search")

	ownerID, ok := utils.RequireOwnerID(ginContext)
	if !ok {
		return
	}

	resultData, err := services.GetAllStaffs(ownerID, page, pageSize, search)

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
		utils.ErrorResponse(ginContext, http.StatusBadRequest, 400, "Vui lòng nhập email và họ tên hợp lệ")
		return
	}

	ownerID, ok := utils.RequireOwnerID(ginContext)
	if !ok {
		return
	}

	errCreate := services.CreateStaffAccount(ownerID, req.Email, req.FullName)
	if errCreate != nil {
		utils.ErrorResponse(ginContext, http.StatusInternalServerError, 500, errCreate.Error())
		return
	}

	utils.SuccessResponse(ginContext, http.StatusCreated, gin.H{"message": "Tạo nhân viên mới thành công. Mật khẩu mặc định là 123456"})
}

func UpdateUserHandler(ginContext *gin.Context) {
	userID, err := utils.ParseUintParam(ginContext, "id")
	if err != nil {
		utils.ErrorResponse(ginContext, http.StatusBadRequest, 400, "ID nhân viên không hợp lệ")
		return
	}

	var updateData map[string]interface{}
	if errBind := ginContext.ShouldBindJSON(&updateData); errBind != nil {
		utils.ErrorResponse(ginContext, http.StatusBadRequest, 400, "Dữ liệu cập nhật không hợp lệ")
		return
	}

	ownerID, ok := utils.RequireOwnerID(ginContext)
	if !ok {
		return
	}

	errService := services.UpdateStaffs(ownerID, userID, updateData)
	if errService != nil {
		utils.ErrorResponse(ginContext, http.StatusInternalServerError, 500, errService.Error())
		return
	}

	utils.SuccessResponse(ginContext, http.StatusOK, gin.H{"message": "Cập nhật nhân viên thành công"})
}

func DeleteUserHandler(ginContext *gin.Context) {
	userID, err := utils.ParseUintParam(ginContext, "id")
	if err != nil {
		utils.ErrorResponse(ginContext, http.StatusBadRequest, 400, "ID nhân viên không hợp lệ")
		return
	}

	ownerID, ok := utils.RequireOwnerID(ginContext)
	if !ok {
		return
	}

	errService := services.DeleteUser(ownerID, userID)
	if errService != nil {
		utils.ErrorResponse(ginContext, http.StatusBadRequest, 400, errService.Error())
		return
	}

	utils.SuccessResponse(ginContext, http.StatusOK, gin.H{"message": "Xóa nhân viên thành công"})
}

func GetMyProfileHandler(c *gin.Context) {
	userID, ok := utils.RequireUserID(c)
	if !ok {
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
	userID, ok := utils.RequireUserID(c)
	if !ok {
		return
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

	utils.SuccessResponse(c, http.StatusOK, gin.H{"message": "Cập nhật thông tin cá nhân thành công"})
}

func ChangeMyPasswordHandler(c *gin.Context) {
	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, 400, "Dữ liệu không hợp lệ")
		return
	}

	userID, ok := utils.RequireUserID(c)
	if !ok {
		return
	}

	if err := services.ChangePassword(userID, req.OldPassword, req.NewPassword); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, map[string]string{"message": "Đổi mật khẩu thành công"})
}

func UploadAvatarHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, 400, "Không tìm thấy file tải lên")
		return
	}

	uploadDir := "./uploads/avatars"
	os.MkdirAll(uploadDir, os.ModePerm)

	filename := fmt.Sprintf("%d_%s", time.Now().Unix(), filepath.Base(file.Filename))
	savePath := filepath.Join(uploadDir, filename)

	if err := c.SaveUploadedFile(file, savePath); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, 500, "Không thể lưu file ảnh")
		return
	}

	fileURL := "/uploads/avatars/" + filename
	utils.SuccessResponse(c, http.StatusOK, map[string]string{"url": fileURL})
}

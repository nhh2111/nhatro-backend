package controllers

import (
	"doAnHTTT_go/models"
	"doAnHTTT_go/services"
	"doAnHTTT_go/utils"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func GetAllTaskHandler(ginContext *gin.Context) {
	page, pageSize := utils.GetPaginationParams(ginContext)
	search := ginContext.Query("search")

	ownerID, ok := utils.RequireOwnerID(ginContext)
	if !ok {
		return
	}

	resultData, err := services.GetAllTask(ownerID, page, pageSize, search)

	if err != nil {
		utils.ErrorResponse(ginContext, http.StatusInternalServerError, 500, "Lỗi lấy danh sách nhiệm vụ: "+err.Error())
		return
	}

	utils.SuccessResponse(ginContext, http.StatusOK, resultData)
}

func CreateTaskHandler(ginContext *gin.Context) {
	var newTask models.Task
	if err := ginContext.ShouldBindJSON(&newTask); err != nil {
		utils.ErrorResponse(ginContext, http.StatusBadRequest, 400, "Dữ liệu không hợp lệ")
		return
	}

	ownerID, ok := utils.RequireOwnerID(ginContext)
	if !ok {
		return
	}

	if err := services.CreateNewTask(ownerID, &newTask); err != nil {
		utils.ErrorResponse(ginContext, http.StatusInternalServerError, 500, err.Error())
		return
	}

	utils.SuccessResponse(ginContext, http.StatusCreated, gin.H{"message": "Tạo nhiệm vụ thành công", "data": newTask})
}

func UpdateTaskHandler(ginContext *gin.Context) {
	taskID, err := utils.ParseUintParam(ginContext, "id")
	if err != nil {
		utils.ErrorResponse(ginContext, http.StatusBadRequest, 400, "ID không hợp lệ")
		return
	}

	var updateData map[string]interface{}
	if err := ginContext.ShouldBindJSON(&updateData); err != nil {
		utils.ErrorResponse(ginContext, http.StatusBadRequest, 400, "Dữ liệu không hợp lệ")
		return
	}

	ownerID, ok := utils.RequireOwnerID(ginContext)
	if !ok {
		return
	}

	if err := services.UpdateTask(ownerID, taskID, updateData); err != nil {
		utils.ErrorResponse(ginContext, http.StatusInternalServerError, 500, err.Error())
		return
	}

	utils.SuccessResponse(ginContext, http.StatusOK, gin.H{"message": "Cập nhật thành công"})
}

func DeleteTaskHandler(ginContext *gin.Context) {
	taskID, err := utils.ParseUintParam(ginContext, "id")
	if err != nil {
		utils.ErrorResponse(ginContext, http.StatusBadRequest, 400, "ID không hợp lệ")
		return
	}

	ownerID, ok := utils.RequireOwnerID(ginContext)
	if !ok {
		return
	}

	if err := services.DeleteTask(ownerID, taskID); err != nil {
		utils.ErrorResponse(ginContext, http.StatusBadRequest, 400, err.Error())
		return
	}
	utils.SuccessResponse(ginContext, http.StatusOK, gin.H{"message": "Xóa nhiệm vụ thành công"})
}

func UploadTaskImageHandler(ginContext *gin.Context) {
	file, err := ginContext.FormFile("image")
	if err != nil {
		utils.ErrorResponse(ginContext, http.StatusBadRequest, 400, "Không tìm thấy file ảnh đính kèm")
		return
	}

	uploadDir := "uploads/tasks"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		utils.ErrorResponse(ginContext, http.StatusInternalServerError, 500, "Lỗi server không thể tạo thư mục lưu ảnh")
		return
	}

	filename := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)
	filepath := fmt.Sprintf("%s/%s", uploadDir, filename)

	if err := ginContext.SaveUploadedFile(file, filepath); err != nil {
		utils.ErrorResponse(ginContext, http.StatusInternalServerError, 500, "Không thể lưu file ảnh")
		return
	}

	imageUrl := "/" + filepath
	utils.SuccessResponse(ginContext, http.StatusOK, gin.H{
		"message":   "Upload ảnh thành công",
		"image_url": imageUrl,
	})
}

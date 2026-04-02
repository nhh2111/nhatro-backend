package controllers

import (
	"doAnHTTT_go/models"
	"doAnHTTT_go/services"
	"doAnHTTT_go/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetAllTaskHandler(ginContext *gin.Context) {
	page, pageSize := utils.GetPaginationParams(ginContext)
	search := ginContext.Query("search")

	ownerIDVal, _ := ginContext.Get("ownerID")
	ownerID := ownerIDVal.(uint)

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
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ"})
		return
	}

	ownerIDVal, _ := ginContext.Get("ownerID")
	ownerID := ownerIDVal.(uint)

	if err := services.CreateNewTask(ownerID, &newTask); err != nil {
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ginContext.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Tạo nhiệm vụ thành công",
		"data":    newTask,
	})
}

func UpdateTaskHandler(ginContext *gin.Context) {
	taskID, _ := strconv.Atoi(ginContext.Param("id"))

	var updateData map[string]interface{}
	if err := ginContext.ShouldBindJSON(&updateData); err != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ"})
		return
	}

	ownerIDVal, _ := ginContext.Get("ownerID")
	ownerID := ownerIDVal.(uint)

	if err := services.UpdateTask(ownerID, uint(taskID), updateData); err != nil {
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ginContext.JSON(http.StatusOK, gin.H{"message": "Cập nhật thành công"})
}

func DeleteTaskHandler(ginContext *gin.Context) {
	taskID, _ := strconv.Atoi(ginContext.Param("id"))

	ownerIDVal, _ := ginContext.Get("ownerID")
	ownerID := ownerIDVal.(uint)

	if err := services.DeleteTask(ownerID, uint(taskID)); err != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ginContext.JSON(http.StatusOK, gin.H{"message": "Xóa nhiệm vụ thành công"})
}

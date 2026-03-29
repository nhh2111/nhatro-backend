package controllers

import (
	"doAnHTTT_go/models"
	"doAnHTTT_go/services"
	"doAnHTTT_go/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetAllServiceHandler(ginContext *gin.Context) {
	page, pageSize := utils.GetPaginationParams(ginContext)
	search := ginContext.Query("search")

	resultData, err := services.GetAllService(page, pageSize, search)

	if err != nil {
		utils.ErrorResponse(ginContext, http.StatusInternalServerError, 500, "Lỗi lấy danh sách dịch vụ: "+err.Error())
		return
	}

	utils.SuccessResponse(ginContext, http.StatusOK, resultData)
}

func CreateServiceHandler(ginContext *gin.Context) {
	var newService models.Service

	errBind := ginContext.ShouldBindJSON(&newService)
	if errBind != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Dữ liệu đầu vào không hợp lệ",
			"detail":  errBind.Error(),
		})
		return
	}
	errCreate := services.CreateNewService(&newService)
	if errCreate != nil {
		ginContext.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Không thể tạo dịch vụ mới",
			"detail":  errCreate.Error(),
		})
		return
	}

	ginContext.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Tạo phòng mới thành công",
		"data":    newService,
	})
}

func UpdateServiceHandler(ginContext *gin.Context) {
	serviceID, errParse := strconv.Atoi(ginContext.Param("id"))
	if errParse != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "ID dịch vụ không hợp lệ"})
		return
	}

	var updateData map[string]interface{}
	if errBind := ginContext.ShouldBindJSON(&updateData); errBind != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu cập nhật không hợp lệ"})
		return
	}

	errService := services.UpdateService(uint(serviceID), updateData)
	if errService != nil {
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": errService.Error()})
		return
	}

	ginContext.JSON(http.StatusOK, gin.H{"message": "Cập nhật dịch vụ thành công"})
}

func DeleteServiceHandler(ginContext *gin.Context) {
	serviceID, errParse := strconv.Atoi(ginContext.Param("id"))
	if errParse != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "ID dịch vụ không hợp lệ"})
		return
	}

	errService := services.DeleteService(uint(serviceID))
	if errService != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": errService.Error()})
		return
	}

	ginContext.JSON(http.StatusOK, gin.H{"message": "Xóa dịch vụ thành công"})
}

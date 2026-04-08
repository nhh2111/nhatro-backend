package controllers

import (
	"doAnHTTT_go/models"
	"doAnHTTT_go/services"
	"doAnHTTT_go/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAllServiceHandler(ginContext *gin.Context) {
	page, pageSize := utils.GetPaginationParams(ginContext)
	search := ginContext.Query("search")

	ownerID, ok := utils.RequireOwnerID(ginContext)
	if !ok {
		return
	}

	resultData, err := services.GetAllService(ownerID, page, pageSize, search)

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
		utils.ErrorResponse(ginContext, http.StatusBadRequest, 400, "Dữ liệu đầu vào không hợp lệ: "+errBind.Error())
		return
	}

	ownerID, ok := utils.RequireOwnerID(ginContext)
	if !ok {
		return
	}
	newService.OwnerID = ownerID

	errCreate := services.CreateNewService(&newService)
	if errCreate != nil {
		utils.ErrorResponse(ginContext, http.StatusInternalServerError, 500, "Không thể tạo dịch vụ mới: "+errCreate.Error())
		return
	}

	utils.SuccessResponse(ginContext, http.StatusCreated, gin.H{
		"message": "Tạo dịch vụ mới thành công",
		"data":    newService,
	})
}

func UpdateServiceHandler(ginContext *gin.Context) {
	serviceID, err := utils.ParseUintParam(ginContext, "id")
	if err != nil {
		utils.ErrorResponse(ginContext, http.StatusBadRequest, 400, "ID dịch vụ không hợp lệ")
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

	errService := services.UpdateService(ownerID, serviceID, updateData)
	if errService != nil {
		utils.ErrorResponse(ginContext, http.StatusInternalServerError, 500, errService.Error())
		return
	}

	utils.SuccessResponse(ginContext, http.StatusOK, gin.H{"message": "Cập nhật dịch vụ thành công"})
}

func DeleteServiceHandler(ginContext *gin.Context) {
	serviceID, err := utils.ParseUintParam(ginContext, "id")
	if err != nil {
		utils.ErrorResponse(ginContext, http.StatusBadRequest, 400, "ID dịch vụ không hợp lệ")
		return
	}

	ownerID, ok := utils.RequireOwnerID(ginContext)
	if !ok {
		return
	}

	errService := services.DeleteService(ownerID, serviceID)
	if errService != nil {
		utils.ErrorResponse(ginContext, http.StatusBadRequest, 400, errService.Error())
		return
	}

	utils.SuccessResponse(ginContext, http.StatusOK, gin.H{"message": "Xóa dịch vụ thành công"})
}

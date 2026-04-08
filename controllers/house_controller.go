package controllers

import (
	"doAnHTTT_go/models"
	"doAnHTTT_go/services"
	"doAnHTTT_go/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAllHousesHandler(ginContext *gin.Context) {
	page, pageSize := utils.GetPaginationParams(ginContext)
	search := ginContext.Query("search")

	ownerID, ok := utils.RequireOwnerID(ginContext)
	if !ok {
		return
	}

	resultData, err := services.GetAllHouses(ownerID, page, pageSize, search)
	if err != nil {
		utils.ErrorResponse(ginContext, http.StatusInternalServerError, 500, "Lỗi lấy danh sách nhà: "+err.Error())
		return
	}

	utils.SuccessResponse(ginContext, http.StatusOK, resultData)
}

func CreateHouseHandler(ginContext *gin.Context) {
	var newHouse models.House

	errBind := ginContext.ShouldBindJSON(&newHouse)
	if errBind != nil {
		utils.ErrorResponse(ginContext, http.StatusBadRequest, 400, "Dữ liệu đầu vào không hợp lệ: "+errBind.Error())
		return
	}

	ownerID, ok := utils.RequireOwnerID(ginContext)
	if !ok {
		return
	}
	newHouse.OwnerID = ownerID

	errCreate := services.CreateNewHouse(&newHouse)
	if errCreate != nil {
		utils.ErrorResponse(ginContext, http.StatusInternalServerError, 500, "Không thể tạo nhà mới: "+errCreate.Error())
		return
	}

	utils.SuccessResponse(ginContext, http.StatusCreated, gin.H{
		"message": "Tạo nhà mới thành công",
		"data":    newHouse,
	})
}

func UpdateHouseHandler(ginContext *gin.Context) {
	houseID, err := utils.ParseUintParam(ginContext, "id")
	if err != nil {
		utils.ErrorResponse(ginContext, http.StatusBadRequest, 400, "ID nhà không hợp lệ")
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

	errService := services.UpdateHouse(ownerID, houseID, updateData)
	if errService != nil {
		utils.ErrorResponse(ginContext, http.StatusInternalServerError, 500, errService.Error())
		return
	}

	utils.SuccessResponse(ginContext, http.StatusOK, gin.H{"message": "Cập nhật nhà thành công"})
}

func DeleteHouseHandler(ginContext *gin.Context) {
	houseID, err := utils.ParseUintParam(ginContext, "id")
	if err != nil {
		utils.ErrorResponse(ginContext, http.StatusBadRequest, 400, "ID nhà không hợp lệ")
		return
	}

	ownerID, ok := utils.RequireOwnerID(ginContext)
	if !ok {
		return
	}

	errService := services.DeleteHouse(ownerID, houseID)
	if errService != nil {
		utils.ErrorResponse(ginContext, http.StatusBadRequest, 400, errService.Error())
		return
	}

	utils.SuccessResponse(ginContext, http.StatusOK, gin.H{"message": "Xóa nhà thành công"})
}

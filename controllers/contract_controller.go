package controllers

import (
	"doAnHTTT_go/models"
	"doAnHTTT_go/services"
	"doAnHTTT_go/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAllContractHandler(ginContext *gin.Context) {
	page, pageSize := utils.GetPaginationParams(ginContext)
	search := ginContext.Query("search")

	ownerID, ok := utils.RequireOwnerID(ginContext)
	if !ok {
		return
	}

	resultData, err := services.GetAllContracts(ownerID, page, pageSize, search)

	if err != nil {
		utils.ErrorResponse(ginContext, http.StatusInternalServerError, 500, "Lỗi lấy danh sách hợp đồng: "+err.Error())
		return
	}

	utils.SuccessResponse(ginContext, http.StatusOK, resultData)
}

func CreateContractHandler(c *gin.Context) {
	var contract models.Contract
	if err := c.ShouldBindJSON(&contract); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, 400, "Dữ liệu hợp đồng không hợp lệ")
		return
	}

	ownerID, ok := utils.RequireOwnerID(c)
	if !ok {
		return
	}

	if err := services.CreateContract(ownerID, &contract); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, gin.H{"message": "Lập hợp đồng thành công! Phòng đã được chuyển sang trạng thái Đang Thuê."})
}

func TerminateContractHandler(c *gin.Context) {
	contractID, err := utils.ParseUintParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, 400, "ID không hợp lệ")
		return
	}

	ownerID, ok := utils.RequireOwnerID(c)
	if !ok {
		return
	}

	if err := services.TerminateContract(ownerID, contractID); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, 500, err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, gin.H{"message": "Thanh lý hợp đồng thành công! Phòng đã được trống."})
}

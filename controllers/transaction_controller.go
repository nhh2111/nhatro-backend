package controllers

import (
	"doAnHTTT_go/dto"
	"doAnHTTT_go/services"
	"doAnHTTT_go/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddTransactionHandler(ginContext *gin.Context) {
	var requestData dto.CreateTransactionDTO

	errBind := ginContext.ShouldBindJSON(&requestData)
	if errBind != nil {
		utils.ErrorResponse(ginContext, http.StatusBadRequest, 400, "Dữ liệu phiếu thu/chi không hợp lệ")
		return
	}

	ownerID, ok := utils.RequireOwnerID(ginContext)
	if !ok {
		return
	}

	errService := services.CreateNewTransaction(ownerID, requestData)
	if errService != nil {
		utils.ErrorResponse(ginContext, http.StatusInternalServerError, 500, errService.Error())
		return
	}

	utils.SuccessResponse(ginContext, http.StatusCreated, gin.H{"message": "Ghi nhận phiếu thu/chi thành công"})
}

func GetAllTransactionsHandler(ginContext *gin.Context) {
	page, pageSize := utils.GetPaginationParams(ginContext)
	search := ginContext.Query("search")
	monthYear := ginContext.Query("month_year") // <-- BỔ SUNG DÒNG NÀY

	ownerID, ok := utils.RequireOwnerID(ginContext)
	if !ok {
		return
	}

	// TRUYỀN THÊM monthYear VÀO SERVICE
	resultData, err := services.GetAllTransactions(ownerID, page, pageSize, search, monthYear)

	if err != nil {
		utils.ErrorResponse(ginContext, http.StatusInternalServerError, 500, "Lỗi lấy danh sách giao dịch: "+err.Error())
		return
	}

	utils.SuccessResponse(ginContext, http.StatusOK, resultData)
}

func UpdateTransactionHandler(ginContext *gin.Context) {
	txID, err := utils.ParseUintParam(ginContext, "id")
	if err != nil {
		utils.ErrorResponse(ginContext, http.StatusBadRequest, 400, "ID không hợp lệ")
		return
	}

	var updateData map[string]interface{}
	if errBind := ginContext.ShouldBindJSON(&updateData); errBind != nil {
		utils.ErrorResponse(ginContext, http.StatusBadRequest, 400, "Dữ liệu không hợp lệ")
		return
	}

	ownerID, ok := utils.RequireOwnerID(ginContext)
	if !ok {
		return
	}

	errService := services.UpdateTransaction(ownerID, txID, updateData)
	if errService != nil {
		utils.ErrorResponse(ginContext, http.StatusInternalServerError, 500, errService.Error())
		return
	}
	utils.SuccessResponse(ginContext, http.StatusOK, gin.H{"message": "Cập nhật phiếu thu chi thành công"})
}

func DeleteTransactionHandler(ginContext *gin.Context) {
	txID, err := utils.ParseUintParam(ginContext, "id")
	if err != nil {
		utils.ErrorResponse(ginContext, http.StatusBadRequest, 400, "ID không hợp lệ")
		return
	}

	ownerID, ok := utils.RequireOwnerID(ginContext)
	if !ok {
		return
	}

	errService := services.DeleteTransaction(ownerID, txID)
	if errService != nil {
		utils.ErrorResponse(ginContext, http.StatusBadRequest, 400, errService.Error())
		return
	}
	utils.SuccessResponse(ginContext, http.StatusOK, gin.H{"message": "Xóa phiếu thu chi thành công"})
}

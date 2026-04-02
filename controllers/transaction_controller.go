package controllers

import (
	"doAnHTTT_go/dto"
	"doAnHTTT_go/services"
	"doAnHTTT_go/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AddTransactionHandler(ginContext *gin.Context) {
	var requestData dto.CreateTransactionDTO

	errBind := ginContext.ShouldBindJSON(&requestData)
	if errBind != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu phiếu thu/chi không hợp lệ"})
		return
	}

	ownerIDVal, _ := ginContext.Get("ownerID")
	ownerID := ownerIDVal.(uint)

	errService := services.CreateNewTransaction(ownerID, requestData)
	if errService != nil {
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": errService.Error()})
		return
	}

	ginContext.JSON(http.StatusCreated, gin.H{"message": "Ghi nhận phiếu thu/chi thành công"})
}

func GetAllTransactionsHandler(ginContext *gin.Context) {
	page, pageSize := utils.GetPaginationParams(ginContext)
	search := ginContext.Query("search")

	ownerIDVal, _ := ginContext.Get("ownerID")
	ownerID := ownerIDVal.(uint)

	resultData, err := services.GetAllTransactions(ownerID, page, pageSize, search)

	if err != nil {
		utils.ErrorResponse(ginContext, http.StatusInternalServerError, 500, "Lỗi lấy danh sách giao dịch: "+err.Error())
		return
	}

	utils.SuccessResponse(ginContext, http.StatusOK, resultData)
}

func UpdateTransactionHandler(ginContext *gin.Context) {
	txID, errParse := strconv.Atoi(ginContext.Param("id"))
	if errParse != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
		return
	}

	var updateData map[string]interface{}
	if errBind := ginContext.ShouldBindJSON(&updateData); errBind != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ"})
		return
	}

	ownerIDVal, _ := ginContext.Get("ownerID")
	ownerID := ownerIDVal.(uint)

	errService := services.UpdateTransaction(ownerID, uint(txID), updateData)
	if errService != nil {
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": errService.Error()})
		return
	}
	ginContext.JSON(http.StatusOK, gin.H{"message": "Cập nhật phiếu thu chi thành công"})
}

func DeleteTransactionHandler(ginContext *gin.Context) {
	txID, errParse := strconv.Atoi(ginContext.Param("id"))
	if errParse != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
		return
	}

	ownerIDVal, _ := ginContext.Get("ownerID")
	ownerID := ownerIDVal.(uint)

	errService := services.DeleteTransaction(ownerID, uint(txID))
	if errService != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": errService.Error()})
		return
	}
	ginContext.JSON(http.StatusOK, gin.H{"message": "Xóa phiếu thu chi thành công"})
}

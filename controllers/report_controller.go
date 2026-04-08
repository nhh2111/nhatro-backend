package controllers

import (
	"doAnHTTT_go/services"
	"doAnHTTT_go/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetProfitLossHandler(ginContext *gin.Context) {
	month := ginContext.Query("month")
	year := ginContext.Query("year")

	if month == "" || year == "" {
		utils.ErrorResponse(ginContext, http.StatusBadRequest, 400, "Vui lòng cung cấp tháng và năm cần thống kê")
		return
	}

	ownerID, ok := utils.RequireOwnerID(ginContext)
	if !ok {
		return
	}

	reportData, errService := services.GetProfitLossReport(ownerID, month, year)
	if errService != nil {
		utils.ErrorResponse(ginContext, http.StatusInternalServerError, 500, errService.Error())
		return
	}

	utils.SuccessResponse(ginContext, http.StatusOK, gin.H{"data": reportData})
}

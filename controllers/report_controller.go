package controllers

import (
	"doAnHTTT_go/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetProfitLossHandler(ginContext *gin.Context) {
	month := ginContext.Query("month")
	year := ginContext.Query("year")

	if month == "" || year == "" {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "Vui lòng cung cấp tháng và năm cần thống kê"})
		return
	}

	ownerIDVal, _ := ginContext.Get("ownerID")
	ownerID := ownerIDVal.(uint)

	reportData, errService := services.GetProfitLossReport(ownerID, month, year)
	if errService != nil {
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": errService.Error()})
		return
	}

	ginContext.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   reportData,
	})
}

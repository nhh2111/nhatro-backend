package controllers

import (
	"doAnHTTT_go/dto"
	"doAnHTTT_go/services"
	"doAnHTTT_go/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type InvoiceReq struct {
	MonthYear string `json:"month_year" binding:"required"`
}

func GetAllInvoicesHandler(ginContext *gin.Context) {
	page, pageSize := utils.GetPaginationParams(ginContext)
	monthYear := ginContext.Query("month_year")

	ownerIDVal, _ := ginContext.Get("ownerID")
	ownerID := ownerIDVal.(uint)

	resultData, err := services.GetAllInvoices(ownerID, page, pageSize, monthYear)

	if err != nil {
		utils.ErrorResponse(ginContext, http.StatusInternalServerError, 500, "Lỗi tải hóa đơn: "+err.Error())
		return
	}

	utils.SuccessResponse(ginContext, http.StatusOK, resultData)
}

func TriggerGenerateInvoices(c *gin.Context) {
	var req InvoiceReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Thiếu tháng cần chốt sổ"})
		return
	}

	ownerIDVal, _ := c.Get("ownerID")
	ownerID := ownerIDVal.(uint)

	errService := services.AutoGenerateInvoices(ownerID, req.MonthYear)
	if errService != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errService.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Đã chốt sổ hóa đơn thành công cho tháng " + req.MonthYear})
}

func DeleteInvoiceHandler(ginContext *gin.Context) {
	invoiceIDParam := ginContext.Param("id")
	invoiceID, errParse := strconv.Atoi(invoiceIDParam)

	if errParse != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "ID hóa đơn không hợp lệ"})
		return
	}

	roleVal, _ := ginContext.Get("userRole")
	userRole := roleVal.(string)

	ownerIDVal, _ := ginContext.Get("ownerID")
	ownerID := ownerIDVal.(uint)

	errService := services.DeleteInvoice(ownerID, uint(invoiceID), userRole)
	if errService != nil {
		ginContext.JSON(http.StatusForbidden, gin.H{"error": errService.Error()})
		return
	}

	ginContext.JSON(http.StatusOK, gin.H{"message": "Đã xóa hóa đơn thành công"})
}

func PayInvoiceHandler(ginContext *gin.Context) {
	var requestData dto.PayInvoiceDTO

	if errBind := ginContext.ShouldBindJSON(&requestData); errBind != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu thanh toán không hợp lệ"})
		return
	}

	ownerIDVal, _ := ginContext.Get("ownerID")
	ownerID := ownerIDVal.(uint)

	errService := services.PayInvoice(ownerID, requestData)
	if errService != nil {
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": errService.Error()})
		return
	}

	ginContext.JSON(http.StatusOK, gin.H{"message": "Đã thu tiền và cập nhật hóa đơn thành công"})
}

package controllers

import (
	"doAnHTTT_go/dto"
	"doAnHTTT_go/services"
	"doAnHTTT_go/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type InvoiceReq struct {
	MonthYear string `json:"month_year" binding:"required"`
}

func GetAllInvoicesHandler(ginContext *gin.Context) {
	page, pageSize := utils.GetPaginationParams(ginContext)
	monthYear := ginContext.Query("month_year")

	ownerID, ok := utils.RequireOwnerID(ginContext)
	if !ok {
		return
	}

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
		utils.ErrorResponse(c, http.StatusBadRequest, 400, "Thiếu tháng cần chốt sổ")
		return
	}

	ownerID, ok := utils.RequireOwnerID(c)
	if !ok {
		return
	}

	errService := services.AutoGenerateInvoices(ownerID, req.MonthYear)
	if errService != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, 500, errService.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, gin.H{"message": "Đã chốt sổ hóa đơn thành công cho tháng " + req.MonthYear})
}

func DeleteInvoiceHandler(ginContext *gin.Context) {
	invoiceID, err := utils.ParseUintParam(ginContext, "id")
	if err != nil {
		utils.ErrorResponse(ginContext, http.StatusBadRequest, 400, "ID hóa đơn không hợp lệ")
		return
	}

	userRole, ok := utils.RequireUserRole(ginContext)
	if !ok {
		return
	}

	ownerID, ok := utils.RequireOwnerID(ginContext)
	if !ok {
		return
	}

	errService := services.DeleteInvoice(ownerID, invoiceID, userRole)
	if errService != nil {
		utils.ErrorResponse(ginContext, http.StatusForbidden, 403, errService.Error())
		return
	}

	utils.SuccessResponse(ginContext, http.StatusOK, gin.H{"message": "Đã xóa hóa đơn thành công"})
}

func PayInvoiceHandler(ginContext *gin.Context) {
	var requestData dto.PayInvoiceDTO

	if errBind := ginContext.ShouldBindJSON(&requestData); errBind != nil {
		utils.ErrorResponse(ginContext, http.StatusBadRequest, 400, "Dữ liệu thanh toán không hợp lệ")
		return
	}

	ownerID, ok := utils.RequireOwnerID(ginContext)
	if !ok {
		return
	}

	errService := services.PayInvoice(ownerID, requestData)
	if errService != nil {
		utils.ErrorResponse(ginContext, http.StatusInternalServerError, 500, errService.Error())
		return
	}

	utils.SuccessResponse(ginContext, http.StatusOK, gin.H{"message": "Đã thu tiền và cập nhật hóa đơn thành công"})
}

type BankWebhookPayload struct {
	Gateway            string  `json:"gateway"`
	TransactionDate    string  `json:"transactionDate"`
	AccountNumber      string  `json:"accountNumber"`
	AmountIn           float64 `json:"amountIn"`
	TransactionContent string  `json:"transactionContent"`
	ReferenceCode      string  `json:"referenceCode"`
}

func WebhookBankTransferHandler(c *gin.Context) {
	var payload BankWebhookPayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, 400, "Payload không hợp lệ")
		return
	}

	if payload.AmountIn <= 0 {
		utils.SuccessResponse(c, http.StatusOK, gin.H{"message": "Bỏ qua giao dịch chi tiền"})
		return
	}

	err := services.ProcessBankWebhook(payload.AmountIn, payload.TransactionContent)
	if err != nil {
		utils.ErrorResponse(c, http.StatusOK, 200, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, gin.H{"success": true, "message": "Gạch nợ hóa đơn thành công"})
}

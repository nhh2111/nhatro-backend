package controllers

import (
	"doAnHTTT_go/dto"
	"doAnHTTT_go/services"
	"doAnHTTT_go/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func AddMeterReadingHandler(ginContext *gin.Context) {
	var requestData dto.CreateMeterReadingDTO
	errBind := ginContext.ShouldBindJSON(&requestData)
	if errBind != nil {
		utils.ErrorResponse(ginContext, http.StatusBadRequest, 400, "Dữ liệu gửi lên không đúng định dạng")
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

	errService := services.CreateNewMeterReading(ownerID, requestData, userRole)
	if errService != nil {
		utils.ErrorResponse(ginContext, http.StatusBadRequest, 400, errService.Error())
		return
	}

	utils.SuccessResponse(ginContext, http.StatusCreated, gin.H{"message": "Ghi nhận chỉ số điện/nước thành công!"})
}

func GetMeterReadingsHandler(ginContext *gin.Context) {
	month := ginContext.Query("month")
	if month == "" {
		month = time.Now().Format("2006-01")
	}

	ownerID, ok := utils.RequireOwnerID(ginContext)
	if !ok {
		return
	}

	readings, err := services.GetMeterReadingsByMonth(ownerID, month)
	if err != nil {
		utils.ErrorResponse(ginContext, http.StatusInternalServerError, 500, err.Error())
		return
	}

	utils.SuccessResponse(ginContext, http.StatusOK, gin.H{"data": readings})
}

func UpdateMeterReadingHandler(ginContext *gin.Context) {
	id, err := utils.ParseUintParam(ginContext, "id")
	if err != nil {
		utils.ErrorResponse(ginContext, http.StatusBadRequest, 400, "ID không hợp lệ")
		return
	}
	var requestData dto.CreateMeterReadingDTO

	if err := ginContext.ShouldBindJSON(&requestData); err != nil {
		utils.ErrorResponse(ginContext, http.StatusBadRequest, 400, "Dữ liệu không hợp lệ")
		return
	}

	ownerID, ok := utils.RequireOwnerID(ginContext)
	if !ok {
		return
	}

	if err := services.UpdateMeterReading(ownerID, id, requestData); err != nil {
		utils.ErrorResponse(ginContext, http.StatusBadRequest, 400, err.Error())
		return
	}
	utils.SuccessResponse(ginContext, http.StatusOK, gin.H{"message": "Cập nhật chỉ số thành công!"})
}

func DeleteMeterReadingHandler(ginContext *gin.Context) {
	id, err := utils.ParseUintParam(ginContext, "id")
	if err != nil {
		utils.ErrorResponse(ginContext, http.StatusBadRequest, 400, "ID không hợp lệ")
		return
	}

	ownerID, ok := utils.RequireOwnerID(ginContext)
	if !ok {
		return
	}

	if err := services.DeleteMeterReading(ownerID, id); err != nil {
		utils.ErrorResponse(ginContext, http.StatusBadRequest, 400, err.Error())
		return
	}
	utils.SuccessResponse(ginContext, http.StatusOK, gin.H{"message": "Xóa chỉ số thành công!"})
}

func GetLatestIndexHandler(ginContext *gin.Context) {
	roomIDStr := ginContext.Query("room_id")
	serviceIDStr := ginContext.Query("service_id")
	date := ginContext.Query("date")

	roomID, _ := strconv.ParseUint(roomIDStr, 10, 32)
	serviceID, _ := strconv.ParseUint(serviceIDStr, 10, 32)

	if roomID == 0 || serviceID == 0 || date == "" {
		utils.ErrorResponse(ginContext, http.StatusBadRequest, 400, "Thiếu tham số phòng, dịch vụ hoặc ngày ghi")
		return
	}

	// Gọi service lấy số
	oldIndex := services.GetLatestOldIndex(uint(roomID), uint(serviceID), date)

	utils.SuccessResponse(ginContext, http.StatusOK, gin.H{"old_index": oldIndex})
}

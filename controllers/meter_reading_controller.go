package controllers

import (
	"doAnHTTT_go/dto"
	"doAnHTTT_go/services"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func AddMeterReadingHandler(ginContext *gin.Context) {
	var requestData dto.CreateMeterReadingDTO
	errBind := ginContext.ShouldBindJSON(&requestData)
	if errBind != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  "Dữ liệu gửi lên không đúng định dạng",
		})
		return
	}

	roleVal, exists := ginContext.Get("userRole")
	if !exists {
		ginContext.JSON(http.StatusUnauthorized, gin.H{"error": "Không xác định được quyền người dùng"})
		return
	}
	userRole := roleVal.(string)

	ownerIDVal, _ := ginContext.Get("ownerID")
	ownerID := ownerIDVal.(uint)

	errService := services.CreateNewMeterReading(ownerID, requestData, userRole)
	if errService != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  errService.Error(),
		})
		return
	}

	ginContext.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Ghi nhận chỉ số điện/nước thành công!",
	})
}

func GetMeterReadingsHandler(ginContext *gin.Context) {
	month := ginContext.Query("month")
	if month == "" {
		month = time.Now().Format("2006-01")
	}

	ownerIDVal, _ := ginContext.Get("ownerID")
	ownerID := ownerIDVal.(uint)

	readings, err := services.GetMeterReadingsByMonth(ownerID, month)
	if err != nil {
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ginContext.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   readings,
	})
}

func UpdateMeterReadingHandler(ginContext *gin.Context) {
	id, _ := strconv.Atoi(ginContext.Param("id"))
	var requestData dto.CreateMeterReadingDTO

	if err := ginContext.ShouldBindJSON(&requestData); err != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ"})
		return
	}

	ownerIDVal, _ := ginContext.Get("ownerID")
	ownerID := ownerIDVal.(uint)

	if err := services.UpdateMeterReading(ownerID, uint(id), requestData); err != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ginContext.JSON(http.StatusOK, gin.H{"message": "Cập nhật chỉ số thành công!"})
}

func DeleteMeterReadingHandler(ginContext *gin.Context) {
	id, _ := strconv.Atoi(ginContext.Param("id"))

	ownerIDVal, _ := ginContext.Get("ownerID")
	ownerID := ownerIDVal.(uint)

	if err := services.DeleteMeterReading(ownerID, uint(id)); err != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ginContext.JSON(http.StatusOK, gin.H{"message": "Xóa chỉ số thành công!"})
}

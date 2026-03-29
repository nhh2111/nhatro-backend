package controllers

import (
	"doAnHTTT_go/models"
	"doAnHTTT_go/services"
	"doAnHTTT_go/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetAllHousesHandler(ginContext *gin.Context) {
	page, pageSize := utils.GetPaginationParams(ginContext)
	search := ginContext.Query("search")

	resultData, err := services.GetAllHouses(page, pageSize, search)

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
		ginContext.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Dữ liệu đầu vào không hợp lệ",
			"detail":  errBind.Error(),
		})
		return
	}

	userIDVal, exists := ginContext.Get("userID")
	if !exists {
		ginContext.JSON(http.StatusUnauthorized, gin.H{"error": "Không xác định được danh tính người dùng"})
		return
	}
	newHouse.OwnerID = userIDVal.(uint)

	errCreate := services.CreateNewHouse(&newHouse)
	if errCreate != nil {
		ginContext.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Không thể tạo nhà mới",
			"detail":  errCreate.Error(),
		})
		return
	}

	ginContext.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Tạo nhà mới thành công",
		"data":    newHouse,
	})
}

func UpdateHouseHandler(ginContext *gin.Context) {
	houseID, errParse := strconv.Atoi(ginContext.Param("id"))
	if errParse != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "ID nhà không hợp lệ"})
		return
	}

	var updateData map[string]interface{}
	if errBind := ginContext.ShouldBindJSON(&updateData); errBind != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu cập nhật không hợp lệ"})
		return
	}

	errService := services.UpdateHouse(uint(houseID), updateData)
	if errService != nil {
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": errService.Error()})
		return
	}

	ginContext.JSON(http.StatusOK, gin.H{"message": "Cập nhật nhà thành công"})
}

func DeleteHouseHandler(ginContext *gin.Context) {
	houseID, errParse := strconv.Atoi(ginContext.Param("id"))
	if errParse != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "ID nhà không hợp lệ"})
		return
	}

	errService := services.DeleteHouse(uint(houseID))
	if errService != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": errService.Error()})
		return
	}

	ginContext.JSON(http.StatusOK, gin.H{"message": "Xóa nhà thành công"})
}

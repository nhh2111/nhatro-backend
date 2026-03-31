package controllers

import (
	"doAnHTTT_go/models"
	"doAnHTTT_go/services"
	"doAnHTTT_go/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetAllRoomHandler(ginContext *gin.Context) {
	page, pageSize := utils.GetPaginationParams(ginContext)
	search := ginContext.Query("search")
	houseIdStr := ginContext.Query("house_id") // Bắt thêm tham số house_id

	var houseId uint
	if houseIdStr != "" {
		id, err := strconv.Atoi(houseIdStr)
		if err == nil {
			houseId = uint(id)
		}
	}

	// Truyền thêm houseId vào service
	resultData, err := services.GetAllRooms(page, pageSize, search, houseId)

	if err != nil {
		utils.ErrorResponse(ginContext, http.StatusInternalServerError, 500, "Lỗi lấy danh sách phòng: "+err.Error())
		return
	}

	utils.SuccessResponse(ginContext, http.StatusOK, resultData)
}

func CreateRoomHandler(ginContext *gin.Context) {
	var newRoom models.Room

	errBind := ginContext.ShouldBindJSON(&newRoom)
	if errBind != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Dữ liệu đầu vào không hợp lệ",
			"detail":  errBind.Error(),
		})
		return
	}
	errCreate := services.CreateNewRoom(&newRoom)
	if errCreate != nil {
		ginContext.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Không thể tạo phòng mới",
			"detail":  errCreate.Error(),
		})
		return
	}

	ginContext.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Tạo phòng mới thành công",
		"data":    newRoom,
	})
}

func UpdateRoomHandler(ginContext *gin.Context) {
	roomID, errParse := strconv.Atoi(ginContext.Param("id"))
	if errParse != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "ID phòng không hợp lệ"})
		return
	}

	var updateData map[string]interface{}
	if errBind := ginContext.ShouldBindJSON(&updateData); errBind != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu cập nhật không hợp lệ"})
		return
	}

	errService := services.UpdateRoom(uint(roomID), updateData)
	if errService != nil {
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": errService.Error()})
		return
	}

	ginContext.JSON(http.StatusOK, gin.H{"message": "Cập nhật phòng thành công"})
}

func DeleteRoomHandler(ginContext *gin.Context) {
	roomID, errParse := strconv.Atoi(ginContext.Param("id"))
	if errParse != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "ID phòng không hợp lệ"})
		return
	}

	errService := services.DeleteRoom(uint(roomID))
	if errService != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": errService.Error()})
		return
	}

	ginContext.JSON(http.StatusOK, gin.H{"message": "Xóa phòng thành công"})
}

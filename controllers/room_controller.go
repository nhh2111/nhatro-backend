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
	houseIdStr := ginContext.Query("house_id")

	var houseId uint
	if houseIdStr != "" {
		id, err := strconv.Atoi(houseIdStr)
		if err == nil {
			houseId = uint(id)
		}
	}

	ownerID, ok := utils.RequireOwnerID(ginContext)
	if !ok {
		return
	}

	resultData, err := services.GetAllRooms(ownerID, page, pageSize, search, houseId)
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
		utils.ErrorResponse(ginContext, http.StatusBadRequest, 400, "Dữ liệu đầu vào không hợp lệ: "+errBind.Error())
		return
	}

	ownerID, ok := utils.RequireOwnerID(ginContext)
	if !ok {
		return
	}

	errCreate := services.CreateNewRoom(ownerID, &newRoom)
	if errCreate != nil {
		utils.ErrorResponse(ginContext, http.StatusInternalServerError, 500, "Không thể tạo phòng mới: "+errCreate.Error())
		return
	}

	utils.SuccessResponse(ginContext, http.StatusCreated, gin.H{
		"message": "Tạo phòng mới và gán dịch vụ mặc định thành công",
		"data":    newRoom,
	})
}

func UpdateRoomHandler(ginContext *gin.Context) {
	roomID, errParse := strconv.Atoi(ginContext.Param("id"))
	if errParse != nil {
		utils.ErrorResponse(ginContext, http.StatusBadRequest, 400, "ID phòng không hợp lệ")
		return
	}

	var updateData map[string]interface{}
	if errBind := ginContext.ShouldBindJSON(&updateData); errBind != nil {
		utils.ErrorResponse(ginContext, http.StatusBadRequest, 400, "Dữ liệu cập nhật không hợp lệ")
		return
	}

	ownerID, ok := utils.RequireOwnerID(ginContext)
	if !ok {
		return
	}

	errService := services.UpdateRoom(ownerID, uint(roomID), updateData)
	if errService != nil {
		utils.ErrorResponse(ginContext, http.StatusBadRequest, 400, errService.Error())
		return
	}

	utils.SuccessResponse(ginContext, http.StatusOK, gin.H{"message": "Cập nhật phòng thành công"})
}

func DeleteRoomHandler(ginContext *gin.Context) {
	roomID, errParse := strconv.Atoi(ginContext.Param("id"))
	if errParse != nil {
		utils.ErrorResponse(ginContext, http.StatusBadRequest, 400, "ID phòng không hợp lệ")
		return
	}

	ownerID, ok := utils.RequireOwnerID(ginContext)
	if !ok {
		return
	}

	errService := services.DeleteRoom(ownerID, uint(roomID))
	if errService != nil {
		utils.ErrorResponse(ginContext, http.StatusBadRequest, 400, errService.Error())
		return
	}

	utils.SuccessResponse(ginContext, http.StatusOK, gin.H{"message": "Xóa phòng thành công"})
}

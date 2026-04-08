package controllers

import (
	"doAnHTTT_go/config"
	"doAnHTTT_go/models"
	"doAnHTTT_go/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetServicesOfRoomHandler(c *gin.Context) {
	roomID, err := utils.ParseUintParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, 400, "ID phòng không hợp lệ")
		return
	}

	ownerID, ok := utils.RequireOwnerID(c)
	if !ok {
		return
	}

	// CHẶN BẢO MẬT: Kiểm tra phòng này có phải của chủ trọ không
	var count int64
	config.DB.Table("rooms").Joins("JOIN houses ON rooms.house_id = houses.id").
		Where("rooms.id = ? AND houses.owner_id = ?", roomID, ownerID).Count(&count)

	if count == 0 {
		utils.ErrorResponse(c, http.StatusForbidden, 403, "Phòng không hợp lệ hoặc bạn không có quyền")
		return
	}

	var roomServices []models.RoomService
	if err := config.DB.Where("room_id = ?", roomID).Find(&roomServices).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, 500, "Lỗi khi lấy dịch vụ của phòng")
		return
	}

	var serviceIDs []uint
	for _, rs := range roomServices {
		serviceIDs = append(serviceIDs, rs.ServiceID)
	}

	utils.SuccessResponse(c, http.StatusOK, serviceIDs)
}

func AssignServicesToRoomHandler(c *gin.Context) {
	roomID, err := utils.ParseUintParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, 400, "ID phòng không hợp lệ")
		return
	}

	var request struct {
		ServiceIDs []uint `json:"service_ids"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, 400, "Dữ liệu không hợp lệ")
		return
	}

	ownerID, ok := utils.RequireOwnerID(c)
	if !ok {
		return
	}

	// CHẶN BẢO MẬT
	var count int64
	config.DB.Table("rooms").Joins("JOIN houses ON rooms.house_id = houses.id").
		Where("rooms.id = ? AND houses.owner_id = ?", roomID, ownerID).Count(&count)

	if count == 0 {
		utils.ErrorResponse(c, http.StatusForbidden, 403, "Phòng không hợp lệ hoặc bạn không có quyền")
		return
	}

	tx := config.DB.Begin()

	if err := tx.Where("room_id = ?", roomID).Delete(&models.RoomService{}).Error; err != nil {
		tx.Rollback()
		utils.ErrorResponse(c, http.StatusInternalServerError, 500, "Lỗi khi dọn dẹp dịch vụ cũ")
		return
	}

	for _, srvID := range request.ServiceIDs {
		newRS := models.RoomService{
			RoomID:    roomID,
			ServiceID: srvID,
		}
		if err := tx.Create(&newRS).Error; err != nil {
			tx.Rollback()
			utils.ErrorResponse(c, http.StatusInternalServerError, 500, "Lỗi khi gán dịch vụ mới")
			return
		}
	}

	tx.Commit()
	utils.SuccessResponse(c, http.StatusOK, gin.H{"message": "Cập nhật dịch vụ cho phòng thành công"})
}

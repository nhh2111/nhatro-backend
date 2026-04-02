package controllers

import (
	"doAnHTTT_go/config"
	"doAnHTTT_go/models"
	"doAnHTTT_go/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetServicesOfRoomHandler(c *gin.Context) {
	roomID := c.Param("id")

	ownerIDVal, _ := c.Get("ownerID")
	ownerID := ownerIDVal.(uint)

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
	roomIDStr := c.Param("id")
	roomID, errParse := strconv.Atoi(roomIDStr)
	if errParse != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID phòng không hợp lệ"})
		return
	}

	var request struct {
		ServiceIDs []uint `json:"service_ids"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ"})
		return
	}

	ownerIDVal, _ := c.Get("ownerID")
	ownerID := ownerIDVal.(uint)

	// CHẶN BẢO MẬT
	var count int64
	config.DB.Table("rooms").Joins("JOIN houses ON rooms.house_id = houses.id").
		Where("rooms.id = ? AND houses.owner_id = ?", roomID, ownerID).Count(&count)

	if count == 0 {
		c.JSON(http.StatusForbidden, gin.H{"error": "Phòng không hợp lệ hoặc bạn không có quyền"})
		return
	}

	tx := config.DB.Begin()

	if err := tx.Where("room_id = ?", roomID).Delete(&models.RoomService{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi dọn dẹp dịch vụ cũ"})
		return
	}

	for _, srvID := range request.ServiceIDs {
		newRS := models.RoomService{
			RoomID:    uint(roomID),
			ServiceID: srvID,
		}
		if err := tx.Create(&newRS).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi gán dịch vụ mới"})
			return
		}
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"message": "Cập nhật dịch vụ cho phòng thành công"})
}

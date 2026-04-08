package services

import (
	"doAnHTTT_go/config"
	"doAnHTTT_go/models"
	"doAnHTTT_go/utils"
	"errors"
)

func GetAllRooms(ownerID uint, page int, pageSize int, search string, houseId uint) (map[string]interface{}, error) {
	var roomList []models.Room
	var totalRecords int64
	countQuery := config.DB.Table("rooms").
		Joins("JOIN houses ON houses.id = rooms.house_id").
		Where("houses.owner_id = ?", ownerID)

	if search != "" {
		searchKeyword := "%" + search + "%"
		countQuery = countQuery.Where("(rooms.room_number LIKE ? OR rooms.description LIKE ?)", searchKeyword, searchKeyword)
	}
	if houseId > 0 {
		countQuery = countQuery.Where("rooms.house_id = ?", houseId)
	}
	countQuery.Count(&totalRecords)

	query := config.DB.Table("rooms").
		Select("rooms.*, COALESCE((SELECT COUNT(id) FROM contracts WHERE contracts.room_id = rooms.id AND contracts.status = 'ACTIVE'), 0) AS current_occupants").
		Joins("JOIN houses ON houses.id = rooms.house_id").
		Where("houses.owner_id = ?", ownerID)

	if search != "" {
		searchKeyword := "%" + search + "%"
		query = query.Where("(rooms.room_number LIKE ? OR houses.name LIKE ?)", searchKeyword, searchKeyword)
	}
	if houseId > 0 {
		query = query.Where("rooms.house_id = ?", houseId)
	}

	pageCount := utils.GetPageCount(totalRecords, pageSize)
	offset := utils.GetOffset(page, pageSize)

	result := query.Offset(offset).Limit(pageSize).Order("rooms.id DESC").Find(&roomList)
	if result.Error != nil {
		return nil, result.Error
	}

	return map[string]interface{}{
		"recordCount": totalRecords,
		"pageCount":   pageCount,
		"currentPage": page,
		"pageSize":    pageSize,
		"records":     roomList,
	}, nil
}

func CreateNewRoom(ownerID uint, newRoom *models.Room) error {
	var houseCount int64
	config.DB.Model(&models.House{}).Where("id = ? AND owner_id = ?", newRoom.HouseID, ownerID).Count(&houseCount)
	if houseCount == 0 {
		return errors.New("khu trọ không tồn tại hoặc bạn không có quyền thêm phòng vào đây")
	}

	tx := config.DB.Begin()

	if err := tx.Create(newRoom).Error; err != nil {
		tx.Rollback()
		return err
	}

	defaultServices := []struct {
		Name        string
		ServiceType string
		Unit        string
	}{
		{"Điện", "METER", "kWh"},
		{"Nước", "METER", "Khối"},
		{"Rác & Vệ sinh", "FIXED", "Tháng"},
	}

	for _, ds := range defaultServices {
		var srv models.Service
		err := tx.Where("owner_id = ? AND name = ?", ownerID, ds.Name).First(&srv).Error

		if err != nil {
			srv = models.Service{
				Name:        ds.Name,
				ServiceType: ds.ServiceType,
				UnitPrice:   0,
				Unit:        ds.Unit,
				OwnerID:     ownerID,
			}
			if err := tx.Create(&srv).Error; err != nil {
				tx.Rollback()
				return errors.New("lỗi khi khởi tạo dịch vụ mặc định")
			}
		}

		roomService := models.RoomService{
			RoomID:    newRoom.ID,
			ServiceID: srv.ID,
		}
		if err := tx.Create(&roomService).Error; err != nil {
			tx.Rollback()
			return errors.New("lỗi khi tự động gán dịch vụ cho phòng")
		}
	}

	tx.Commit()
	return nil
}

func UpdateRoom(ownerID uint, roomID uint, updatedData map[string]interface{}) error {
	var room models.Room

	errFind := config.DB.Joins("JOIN houses ON houses.id = rooms.house_id").
		Where("rooms.id = ? AND houses.owner_id = ?", roomID, ownerID).
		First(&room).Error
	if errFind != nil {
		return errors.New("không tìm thấy dữ liệu phòng hoặc bạn không có quyền sửa")
	}

	if newHouseID, ok := updatedData["house_id"]; ok {
		var houseCount int64
		houseIdToFind, ok := utils.CoerceUint(newHouseID)
		if !ok || houseIdToFind == 0 {
			return errors.New("khu trọ đích không hợp lệ")
		}
		config.DB.Model(&models.House{}).Where("id = ? AND owner_id = ?", houseIdToFind, ownerID).Count(&houseCount)
		if houseCount == 0 {
			return errors.New("khu trọ đích không hợp lệ")
		}
	}

	if newStatus, ok := updatedData["status"].(string); ok && newStatus != room.Status {
		if newStatus == "AVAILABLE" || newStatus == "MAINTENANCE" {
			var activeContractCount int64
			config.DB.Table("contracts").
				Where("room_id = ? AND status = 'ACTIVE'", roomID).
				Count(&activeContractCount)

			if activeContractCount > 0 {
				return errors.New("không thể đổi trạng thái thành Trống/Bảo trì vì phòng đang có Khách thuê (Hợp đồng đang hiệu lực)")
			}
		}
	}

	errUpdate := config.DB.Model(&room).Updates(updatedData).Error
	if errUpdate != nil {
		return errors.New("lỗi khi cập nhật thông tin phòng")
	}

	return nil
}

func DeleteRoom(ownerID uint, roomID uint) error {
	var room models.Room

	errFind := config.DB.Joins("JOIN houses ON houses.id = rooms.house_id").
		Where("rooms.id = ? AND houses.owner_id = ?", roomID, ownerID).
		First(&room).Error
	if errFind != nil {
		return errors.New("không tìm thấy dữ liệu phòng hoặc bạn không có quyền xóa")
	}

	if room.Status == "OCCUPIED" || room.Status == "RENTED" {
		return errors.New("không thể xóa phòng đang có khách thuê")
	}

	errDelete := config.DB.Delete(&room).Error
	if errDelete != nil {
		return errors.New("lỗi khi xóa phòng khỏi hệ thống")
	}

	return nil
}

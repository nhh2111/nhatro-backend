package services

import (
	"doAnHTTT_go/config"
	"doAnHTTT_go/models"
	"doAnHTTT_go/utils"
	"errors"

	"gorm.io/gorm"
)

func GetAllContracts(ownerID uint, page int, pageSize int, search string) (map[string]interface{}, error) {
	var contractList []models.Contract
	var totalRecords int64

	// BẢO MẬT: CHỈ LẤY HỢP ĐỒNG THUỘC NHÀ CỦA OWNER NÀY
	query := config.DB.Model(&models.Contract{}).
		Preload("Room").Preload("Tenant").Preload("Room.House").
		Joins("JOIN rooms ON contracts.room_id = rooms.id").
		Joins("JOIN houses ON rooms.house_id = houses.id").
		Where("houses.owner_id = ?", ownerID)

	if search != "" {
		searchKeyword := "%" + search + "%"
		query = query.Where("(rooms.room_number LIKE ? OR tenants.full_name LIKE ?)", searchKeyword, searchKeyword)
	}

	query.Count(&totalRecords)

	pageCount := utils.GetPageCount(totalRecords, pageSize)
	offset := utils.GetOffset(page, pageSize)

	// Thêm tiền tố contracts.id để tránh lỗi mơ hồ (ambiguous) khi JOIN
	result := query.Offset(offset).Limit(pageSize).Order("contracts.id DESC").Find(&contractList)
	if result.Error != nil {
		return nil, result.Error
	}

	return map[string]interface{}{
		"recordCount": totalRecords,
		"pageCount":   pageCount,
		"currentPage": page,
		"pageSize":    pageSize,
		"records":     contractList,
	}, nil
}

func CreateContract(ownerID uint, contract *models.Contract) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		var room models.Room

		// KIỂM TRA BẢO MẬT: Phòng này có thuộc về nhà của chủ trọ đang thao tác không?
		if err := tx.Joins("JOIN houses ON rooms.house_id = houses.id").
			Where("rooms.id = ? AND houses.owner_id = ?", contract.RoomID, ownerID).
			First(&room).Error; err != nil {
			return errors.New("phòng không hợp lệ hoặc bạn không có quyền lập hợp đồng cho phòng này")
		}

		if room.Status == "MAINTENANCE" {
			return errors.New("phòng này đang bảo trì, không thể thuê")
		}

		var currentTenants int64
		tx.Model(&models.Contract{}).Where("room_id = ? AND status = ?", contract.RoomID, "ACTIVE").Count(&currentTenants)

		if currentTenants >= int64(room.MaxOccupants) {
			return errors.New("phòng này đã đủ số lượng người tối đa, không thể ghép thêm")
		}

		contract.Status = "ACTIVE"
		if err := tx.Create(contract).Error; err != nil {
			return errors.New("không thể tạo hợp đồng")
		}

		if err := SyncRoomStatus(tx, contract.RoomID); err != nil {
			return errors.New("lỗi cập nhật trạng thái phòng")
		}
		return nil
	})
}

func TerminateContract(ownerID uint, contractID uint) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		var contract models.Contract

		// KIỂM TRA BẢO MẬT: Hợp đồng này có nằm trong nhà của chủ trọ đang thao tác không?
		if err := tx.Joins("JOIN rooms ON contracts.room_id = rooms.id").
			Joins("JOIN houses ON rooms.house_id = houses.id").
			Where("contracts.id = ? AND houses.owner_id = ?", contractID, ownerID).
			First(&contract).Error; err != nil {
			return errors.New("không tìm thấy hợp đồng hoặc bạn không có quyền thao tác")
		}

		if err := tx.Model(&contract).Update("status", "TERMINATED").Error; err != nil {
			return errors.New("không thể chốt hợp đồng")
		}

		if err := SyncRoomStatus(tx, contract.RoomID); err != nil {
			return errors.New("lỗi cập nhật trạng thái phòng sau khi thanh lý")
		}

		return nil
	})
}

func SyncRoomStatus(tx *gorm.DB, roomID uint) error {
	var activeContracts int64
	tx.Model(&models.Contract{}).Where("room_id = ? AND status = ?", roomID, "ACTIVE").Count(&activeContracts)

	var room models.Room
	if err := tx.First(&room, roomID).Error; err != nil {
		return err
	}

	if activeContracts == 0 {
		room.Status = "AVAILABLE"
	} else {
		room.Status = "OCCUPIED"
	}

	return tx.Save(&room).Error
}

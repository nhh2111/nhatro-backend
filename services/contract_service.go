package services

import (
	"doAnHTTT_go/config"
	"doAnHTTT_go/models"
	"doAnHTTT_go/utils"
	"errors"

	"gorm.io/gorm"
)

func GetAllContracts(page int, pageSize int, search string) (map[string]interface{}, error) {
	var contractList []models.Contract
	var totalRecords int64

	query := config.DB.Model(&models.Contract{}).Preload("Room").Preload("Tenant")

	if search != "" {
		searchKeyword := "%" + search + "%"
		query = query.
			Joins("JOIN rooms ON contracts.room_id = rooms.id").
			Joins("JOIN tenants ON contracts.tenant_id = tenants.id").
			Where("rooms.room_number LIKE ? OR tenants.full_name LIKE ?", searchKeyword, searchKeyword)
	}

	query.Count(&totalRecords)

	pageCount := utils.GetPageCount(totalRecords, pageSize)
	offset := utils.GetOffset(page, pageSize)

	result := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&contractList)
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

func CreateContract(contract *models.Contract) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		var room models.Room
		if err := tx.First(&room, contract.RoomID).Error; err != nil {
			return errors.New("không tìm thấy phòng")
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

func TerminateContract(contractID uint) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		var contract models.Contract
		if err := tx.First(&contract, contractID).Error; err != nil {
			return errors.New("không tìm thấy hợp đồng")
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

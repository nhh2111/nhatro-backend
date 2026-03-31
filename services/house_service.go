package services

import (
	"doAnHTTT_go/config"
	"doAnHTTT_go/models"
	"doAnHTTT_go/utils"
	"errors"
)

func GetAllHouses(page int, pageSize int, search string) (map[string]interface{}, error) {
	var houseList []models.House
	var totalRecords int64

	countQuery := config.DB.Model(&models.House{})
	if search != "" {
		searchKeyword := "%" + search + "%"
		countQuery = countQuery.Where("name LIKE ? OR address LIKE ? OR ward LIKE ? OR district LIKE ? OR city LIKE ?", searchKeyword, searchKeyword, searchKeyword, searchKeyword, searchKeyword)
	}
	countQuery.Count(&totalRecords)

	query := config.DB.Table("houses").
		Select(`houses.*, 
			COUNT(DISTINCT rooms.id) as total_rooms, 
			COALESCE(SUM(
				CASE 
					WHEN rooms.status = 'AVAILABLE' THEN 1 
					WHEN rooms.status = 'OCCUPIED' AND (
						SELECT COUNT(id) FROM contracts WHERE contracts.room_id = rooms.id AND contracts.status = 'ACTIVE'
					) < rooms.max_occupants THEN 1 
					ELSE 0 
				END
			), 0) as empty_rooms`).
		Joins("LEFT JOIN rooms ON rooms.house_id = houses.id").
		Group("houses.id")

	if search != "" {
		searchKeyword := "%" + search + "%"
		query = query.Where("houses.name LIKE ? OR houses.address LIKE ? OR houses.ward LIKE ? OR houses.district LIKE ? OR houses.city LIKE ?", searchKeyword, searchKeyword, searchKeyword, searchKeyword, searchKeyword)
	}

	pageCount := utils.GetPageCount(totalRecords, pageSize)
	offset := utils.GetOffset(page, pageSize)

	result := query.Offset(offset).Limit(pageSize).Order("houses.id DESC").Find(&houseList)
	if result.Error != nil {
		return nil, result.Error
	}

	return map[string]interface{}{
		"recordCount": totalRecords,
		"pageCount":   pageCount,
		"currentPage": page,
		"pageSize":    pageSize,
		"records":     houseList,
	}, nil
}

func CreateNewHouse(newHouse *models.House) error {
	result := config.DB.Create(newHouse)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func UpdateHouse(HouseID uint, updatedData map[string]interface{}) error {
	var house models.House

	errFind := config.DB.First(&house, HouseID).Error
	if errFind != nil {
		return errors.New("không tìm thấy dữ liệu nhà cần sửa")
	}

	errUpdate := config.DB.Model(&house).Updates(updatedData).Error
	if errUpdate != nil {
		return errors.New("lỗi khi cập nhật thông tin nhà")
	}

	return nil
}

func DeleteHouse(HouseID uint) error {
	var house models.House

	errFind := config.DB.First(&house, HouseID).Error
	if errFind != nil {
		return errors.New("không tìm thấy dữ liệu nhà cần xóa")
	}

	errDelete := config.DB.Delete(&house).Error
	if errDelete != nil {
		return errors.New("lỗi khi xóa nhà khỏi hệ thống")
	}

	return nil
}

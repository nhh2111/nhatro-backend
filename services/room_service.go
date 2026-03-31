package services

import (
	"doAnHTTT_go/config"
	"doAnHTTT_go/models"
	"doAnHTTT_go/utils"
	"errors"
)

// Thêm tham số houseId uint vào hàm
func GetAllRooms(page int, pageSize int, search string, houseId uint) (map[string]interface{}, error) {
	var roomList []models.Room
	var totalRecords int64

	query := config.DB.Model(&models.Room{})

	// Lọc theo từ khóa (Số phòng)
	if search != "" {
		searchKeyword := "%" + search + "%"
		query = query.Where("room_number LIKE ? OR description LIKE ?", searchKeyword, searchKeyword)
	}

	// Lọc theo ID nhà (Nếu có truyền lên)
	if houseId > 0 {
		query = query.Where("house_id = ?", houseId)
	}

	query.Count(&totalRecords)

	pageCount := utils.GetPageCount(totalRecords, pageSize)
	offset := utils.GetOffset(page, pageSize)

	result := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&roomList)
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

func CreateNewRoom(newRoom *models.Room) error {
	result := config.DB.Create(newRoom)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func UpdateRoom(roomID uint, updatedData map[string]interface{}) error {
	var room models.Room

	errFind := config.DB.First(&room, roomID).Error
	if errFind != nil {
		return errors.New("không tìm thấy dữ liệu phòng cần sửa")
	}

	errUpdate := config.DB.Model(&room).Updates(updatedData).Error
	if errUpdate != nil {
		return errors.New("lỗi khi cập nhật thông tin phòng")
	}

	return nil
}

func DeleteRoom(roomID uint) error {
	var room models.Room

	errFind := config.DB.First(&room, roomID).Error
	if errFind != nil {
		return errors.New("không tìm thấy dữ liệu phòng cần xóa")
	}

	if room.Status == "OCCUPIED" {
		return errors.New("không thể xóa phòng đang có khách thuê")
	}

	errDelete := config.DB.Delete(&room).Error
	if errDelete != nil {
		return errors.New("lỗi khi xóa phòng khỏi hệ thống")
	}

	return nil
}

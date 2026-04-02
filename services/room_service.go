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

	// BỘ LỌC ĐA KHÁCH HÀNG KHI ĐẾM (JOIN sang houses để lấy owner_id)
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

	// BỘ LỌC ĐA KHÁCH HÀNG KHI LẤY DỮ LIỆU
	query := config.DB.Table("rooms").
		Select("rooms.*, COALESCE((SELECT COUNT(id) FROM contracts WHERE contracts.room_id = rooms.id AND contracts.status = 'ACTIVE'), 0) AS current_occupants").
		Joins("JOIN houses ON houses.id = rooms.house_id").
		Where("houses.owner_id = ?", ownerID)

	if search != "" {
		searchKeyword := "%" + search + "%"
		query = query.Where("(rooms.room_number LIKE ? OR rooms.description LIKE ?)", searchKeyword, searchKeyword)
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
	// KIỂM TRA BẢO MẬT: Nhà (HouseID) mà người dùng định thêm phòng vào có thuộc về họ không?
	var houseCount int64
	config.DB.Model(&models.House{}).Where("id = ? AND owner_id = ?", newRoom.HouseID, ownerID).Count(&houseCount)
	if houseCount == 0 {
		return errors.New("khu trọ không tồn tại hoặc bạn không có quyền thêm phòng vào đây")
	}

	result := config.DB.Create(newRoom)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func UpdateRoom(ownerID uint, roomID uint, updatedData map[string]interface{}) error {
	var room models.Room

	// KIỂM TRA BẢO MẬT: Phải Join với houses để chắc chắn phòng này thuộc nhà của ownerID
	errFind := config.DB.Joins("JOIN houses ON houses.id = rooms.house_id").
		Where("rooms.id = ? AND houses.owner_id = ?", roomID, ownerID).
		First(&room).Error
	if errFind != nil {
		return errors.New("không tìm thấy dữ liệu phòng hoặc bạn không có quyền sửa")
	}

	// Nếu họ cố tình sửa đổi HouseID (Chuyển phòng sang nhà khác), cũng phải kiểm tra nhà mới
	if newHouseID, ok := updatedData["house_id"]; ok {
		var houseCount int64
		// Ép kiểu an toàn (float64 do JSON parse số thành float64)
		houseIdToFind := uint(newHouseID.(float64))
		config.DB.Model(&models.House{}).Where("id = ? AND owner_id = ?", houseIdToFind, ownerID).Count(&houseCount)
		if houseCount == 0 {
			return errors.New("khu trọ đích không hợp lệ")
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

	// KIỂM TRA BẢO MẬT: Join với houses để chắc chắn quyền xóa
	errFind := config.DB.Joins("JOIN houses ON houses.id = rooms.house_id").
		Where("rooms.id = ? AND houses.owner_id = ?", roomID, ownerID).
		First(&room).Error
	if errFind != nil {
		return errors.New("không tìm thấy dữ liệu phòng hoặc bạn không có quyền xóa")
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

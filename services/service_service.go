package services

import (
	"doAnHTTT_go/config"
	"doAnHTTT_go/models"
	"doAnHTTT_go/utils"
	"errors"
)

func GetAllService(ownerID uint, page int, pageSize int, search string) (map[string]interface{}, error) {
	var serviceList []models.Service
	var totalRecords int64

	query := config.DB.Model(&models.Service{}).Where("owner_id = ?", ownerID)

	if search != "" {
		searchKeyword := "%" + search + "%"
		query = query.Where("name LIKE ? ", searchKeyword)
	}

	query.Count(&totalRecords)
	pageCount := utils.GetPageCount(totalRecords, pageSize)
	offset := utils.GetOffset(page, pageSize)

	result := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&serviceList)
	if result.Error != nil {
		return nil, result.Error
	}
	return map[string]interface{}{
		"recordCount": totalRecords,
		"pageCount":   pageCount,
		"currentPage": page,
		"pageSize":    pageSize,
		"records":     serviceList,
	}, nil
}

func CreateNewService(newService *models.Service) error {
	return config.DB.Create(newService).Error
}

func UpdateService(ownerID uint, serviceID uint, updatedData map[string]interface{}) error {
	var service models.Service

	errFind := config.DB.Where("id = ? AND owner_id = ?", serviceID, ownerID).First(&service).Error
	if errFind != nil {
		return errors.New("không tìm thấy dữ liệu dịch vụ hoặc bạn không có quyền sửa")
	}

	errUpdate := config.DB.Model(&service).Updates(updatedData).Error
	if errUpdate != nil {
		return errors.New("lỗi khi cập nhật thông tin dịch vụ")
	}
	return nil
}

func DeleteService(ownerID uint, serviceID uint) error {
	var service models.Service

	errFind := config.DB.Where("id = ? AND owner_id = ?", serviceID, ownerID).First(&service).Error
	if errFind != nil {
		return errors.New("không tìm thấy dữ liệu dịch vụ hoặc bạn không có quyền xóa")
	}

	errDelete := config.DB.Delete(&service).Error
	if errDelete != nil {
		return errors.New("lỗi khi xóa dịch vụ khỏi hệ thống")
	}
	return nil
}

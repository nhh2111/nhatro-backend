package services

import (
	"doAnHTTT_go/config"
	"doAnHTTT_go/models"
	"doAnHTTT_go/utils"
	"errors"
)

func GetAllTenant(ownerID uint, page int, pageSize int, search string) (map[string]interface{}, error) {
	var tenantList []models.Tenant
	var totalRecords int64

	query := config.DB.Model(&models.Tenant{}).Where("owner_id = ?", ownerID)
	if search != "" {
		searchKeyword := "%" + search + "%"
		query = query.Where("(full_name LIKE ? OR phone LIKE ? OR cccd LIKE ? OR license_plates LIKE ?)", searchKeyword, searchKeyword, searchKeyword, searchKeyword)
	}

	query.Count(&totalRecords)

	pageCount := utils.GetPageCount(totalRecords, pageSize)
	offset := utils.GetOffset(page, pageSize)

	result := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&tenantList)
	if result.Error != nil {
		return nil, result.Error
	}

	return map[string]interface{}{
		"recordCount": totalRecords,
		"pageCount":   pageCount,
		"currentPage": page,
		"pageSize":    pageSize,
		"records":     tenantList,
	}, nil
}

func UpdateTenant(ownerID uint, tenantID uint, updatedData map[string]interface{}) error {
	var tenant models.Tenant

	errFind := config.DB.Where("id = ? AND owner_id = ?", tenantID, ownerID).First(&tenant).Error
	if errFind != nil {
		return errors.New("không tìm thấy dữ liệu khách hàng hoặc bạn không có quyền sửa")
	}

	errUpdate := config.DB.Model(&tenant).Updates(updatedData).Error
	if errUpdate != nil {
		return errors.New("lỗi khi cập nhật thông tin khách hàng")
	}

	return nil
}

func DeleteTenant(ownerID uint, tenantID uint) error {
	var tenant models.Tenant

	errFind := config.DB.Where("id = ? AND owner_id = ?", tenantID, ownerID).First(&tenant).Error
	if errFind != nil {
		return errors.New("không tìm thấy dữ liệu khách hàng hoặc bạn không có quyền xóa")
	}

	errDelete := config.DB.Delete(&tenant).Error
	if errDelete != nil {
		return errors.New("lỗi khi xóa khách hàng khỏi hệ thống")
	}

	return nil
}

func CreateNewTenant(newTenant *models.Tenant) error {
	if newTenant.CCCD != "" {
		var existingCCCD models.Tenant
		config.DB.Where("cccd = ? AND owner_id = ?", newTenant.CCCD, newTenant.OwnerID).First(&existingCCCD)
		if existingCCCD.ID != 0 {
			return errors.New("Khách hàng với số CCCD này đã tồn tại trong danh sách của bạn")
		}
	}

	if newTenant.Phone != "" {
		var existingPhone models.Tenant
		config.DB.Where("phone = ? AND owner_id = ?", newTenant.Phone, newTenant.OwnerID).First(&existingPhone)
		if existingPhone.ID != 0 {
			return errors.New("Số điện thoại này đã được đăng ký cho khách khác của bạn")
		}
	}

	result := config.DB.Create(newTenant)
	if result.Error != nil {
		return errors.New("Không thể lưu vào cơ sở dữ liệu: " + result.Error.Error())
	}

	return nil
}

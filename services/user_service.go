package services

import (
	"doAnHTTT_go/config"
	"doAnHTTT_go/models"
	"doAnHTTT_go/utils"
	"errors"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

func GetAllStaffs(employerID uint, page int, pageSize int, search string) (map[string]interface{}, error) {
	var userList []models.User
	var totalRecords int64

	query := config.DB.Model(&models.User{}).Where("role = ? AND employer_id = ?", "STAFF", employerID)
	if search != "" {
		searchKeyword := "%" + search + "%"
		query = query.Where("(full_name LIKE ? OR phone LIKE ? OR email LIKE ?)", searchKeyword, searchKeyword, searchKeyword)
	}

	query.Count(&totalRecords)

	pageCount := utils.GetPageCount(totalRecords, pageSize)
	offset := utils.GetOffset(page, pageSize)

	result := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&userList)
	if result.Error != nil {
		return nil, result.Error
	}

	return map[string]interface{}{
		"recordCount": totalRecords,
		"pageCount":   pageCount,
		"currentPage": page,
		"pageSize":    pageSize,
		"records":     userList,
	}, nil
}

func UpdateStaffs(employerID uint, userID uint, updatedData map[string]interface{}) error {
	var user models.User

	errFind := config.DB.Where("id = ? AND employer_id = ?", userID, employerID).First(&user).Error
	if errFind != nil {
		return errors.New("không tìm thấy dữ liệu nhân viên hoặc bạn không có quyền sửa")
	}

	errUpdate := config.DB.Model(&user).Updates(updatedData).Error
	if errUpdate != nil {
		return errors.New("lỗi khi cập nhật thông tin nhân viên")
	}

	return nil
}

func DeleteUser(employerID uint, userID uint) error {
	var user models.User

	errFind := config.DB.Where("id = ? AND employer_id = ?", userID, employerID).First(&user).Error
	if errFind != nil {
		return errors.New("không tìm thấy dữ liệu nhân viên hoặc bạn không có quyền xóa")
	}

	errDelete := config.DB.Delete(&user).Error
	if errDelete != nil {
		return errors.New("lỗi khi xóa nhân viên khỏi hệ thống")
	}

	return nil
}

func GetMyProfile(userID uint) (map[string]interface{}, error) {
	var user models.User

	if err := config.DB.First(&user, userID).Error; err != nil {
		return nil, errors.New("không tìm thấy thông tin tài khoản")
	}

	profileData := map[string]interface{}{
		"id":         user.ID,
		"full_name":  user.FullName,
		"email":      user.Email,
		"phone":      user.Phone,
		"cccd":       user.CCCD,
		"role":       user.Role,
		"avatar":     user.Avatar,
		"created_at": user.CreatedAt,
	}

	return profileData, nil
}

func UpdateMyProfile(userID uint, updateData map[string]interface{}) error {
	var user models.User

	if err := config.DB.First(&user, userID).Error; err != nil {
		return errors.New("không tìm thấy tài khoản để cập nhật")
	}

	delete(updateData, "id")
	delete(updateData, "role")
	delete(updateData, "email")
	delete(updateData, "password_hash")
	delete(updateData, "is_first_login")
	delete(updateData, "status")

	if err := config.DB.Model(&user).Updates(updateData).Error; err != nil {
		return errors.New("lỗi khi cập nhật thông tin cá nhân")
	}

	return nil
}

func isValidPassword(password string) bool {
	if len(password) < 8 {
		return false
	}
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)

	return hasUpper && hasLower && hasNumber
}

func ChangePassword(userID uint, oldPassword string, newPassword string) error {
	if !isValidPassword(newPassword) {
		return errors.New("Mật khẩu mới quá yếu. Yêu cầu tối thiểu 8 ký tự, bao gồm chữ hoa, ký tự đặc biệt và số.")
	}

	var user models.User

	if err := config.DB.First(&user, userID).Error; err != nil {
		return errors.New("không tìm thấy tài khoản")
	}

	errCompare := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword))
	if errCompare != nil {
		return errors.New("mật khẩu hiện tại không chính xác")
	}

	hashedPassword, errHash := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if errHash != nil {
		return errors.New("lỗi khi mã hóa mật khẩu mới")
	}

	if err := config.DB.Model(&user).Update("password_hash", string(hashedPassword)).Error; err != nil {
		return errors.New("lỗi khi lưu mật khẩu mới")
	}

	return nil
}

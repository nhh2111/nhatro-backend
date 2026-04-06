package services

import (
	"doAnHTTT_go/config"
	"doAnHTTT_go/models"
	"doAnHTTT_go/utils"
	"errors"
	"time"

	"gorm.io/gorm"
)

func RegisterOwner(email, password, fullName string) error {
	errValidate := utils.ValidateStrongPassword(password)
	if errValidate != nil {
		return errValidate
	}
	errVerifyEmail := utils.VerifyEmailDomain(email)
	if errVerifyEmail != nil {
		return errVerifyEmail
	}

	var existingUser models.User
	config.DB.Where("email = ?", email).First(&existingUser)

	if existingUser.ID != 0 {
		if existingUser.Status == true {
			return errors.New("email này đã được đăng ký và kích hoạt")
		}

		hashedPassword, errHash := utils.HashPassword(password)
		if errHash != nil {
			return errors.New("lỗi khi mã hóa mật khẩu")
		}

		config.DB.Model(&existingUser).Updates(map[string]interface{}{
			"password_hash": hashedPassword,
			"full_name":     fullName,
		})

		config.DB.Where("email = ?", email).Delete(&models.OTP{})

		otpCode, _ := utils.GenerateOTP()
		newOTP := models.OTP{
			Email:     email,
			Code:      otpCode,
			ExpiresAt: time.Now().Add(5 * time.Minute),
		}
		config.DB.Create(&newOTP)

		errEmail := utils.SendOTPEmail(email, otpCode, "Mã xác nhận đăng ký tài khoản", "Mã xác thực đăng ký tài khoản Hệ thống Quản lý Nhà Trọ của bạn là:")
		if errEmail != nil {
			return errors.New("tài khoản đã ghi đè nhưng không thể gửi email OTP")
		}

		return nil
	}

	hashedPassword, errHash := utils.HashPassword(password)
	if errHash != nil {
		return errors.New("lỗi khi mã hóa mật khẩu")
	}

	return config.DB.Transaction(func(tx *gorm.DB) error {
		newUser := models.User{
			Email:        email,
			PasswordHash: hashedPassword,
			FullName:     fullName,
			Role:         "OWNER",
			Status:       false,
			IsFirstLogin: false,
		}

		if errCreate := tx.Create(&newUser).Error; errCreate != nil {
			return errors.New("không thể tạo tài khoản")
		}

		otpCode, _ := utils.GenerateOTP()
		newOTP := models.OTP{
			Email:     email,
			Code:      otpCode,
			ExpiresAt: time.Now().Add(5 * time.Minute),
		}

		if errOTP := tx.Create(&newOTP).Error; errOTP != nil {
			return errors.New("lỗi hệ thống khi tạo mã xác nhận: " + errOTP.Error())
		}

		errEmail := utils.SendOTPEmail(email, otpCode, "Mã xác nhận đăng ký tài khoản", "Mã xác thực đăng ký tài khoản Hệ thống Quản lý Nhà Trọ của bạn là:")
		if errEmail != nil {
			return errEmail
		}

		return nil
	})
}

func VerifyRegistrationOTP(email, otpCode string) error {
	var validOTP models.OTP
	config.DB.Where("email = ? AND code = ? AND is_used = ? AND expires_at > ?", email, otpCode, false, time.Now()).First(&validOTP)
	if validOTP.ID == 0 {
		return errors.New("mã OTP không hợp lệ hoặc đã hết hạn")
	}

	errUpdateUser := config.DB.Model(&models.User{}).Where("email = ?", email).Update("status", true).Error
	if errUpdateUser != nil {
		return errors.New("không thể kích hoạt tài khoản")
	}

	if err := config.DB.Delete(&validOTP).Error; err != nil {
		return errors.New("không thể dọn dẹp mã OTP")
	}

	return nil
}

func CreateStaffAccount(employerID uint, email, fullName string) error {
	errVerifyEmail := utils.VerifyEmailDomain(email)
	if errVerifyEmail != nil {
		return errors.New("tên miền email không tồn tại hoặc không thể nhận thư")
	}

	var existingUser models.User
	config.DB.Where("email = ?", email).First(&existingUser)
	if existingUser.ID != 0 {
		return errors.New("email này đã tồn tại trong hệ thống")
	}

	defaultPassword, _ := utils.HashPassword("123456")

	newStaff := models.User{
		Email:        email,
		PasswordHash: defaultPassword,
		FullName:     fullName,
		Role:         "STAFF",
		Status:       true,
		IsFirstLogin: true,
		EmployerID:   int(employerID),
	}

	if errCreate := config.DB.Create(&newStaff).Error; errCreate != nil {
		return errors.New("lỗi khi tạo tài khoản nhân viên")
	}

	return nil
}

func RequestPasswordChangeOTP(email string) error {
	otpCode, _ := utils.GenerateOTP()
	newOTP := models.OTP{
		Email:     email,
		Code:      otpCode,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}
	config.DB.Create(&newOTP)

	errEmail := utils.SendOTPEmail(email, otpCode, "Mã xác nhận đổi mật khẩu", "Bạn đang yêu cầu thay đổi mật khẩu. Mã xác nhận của bạn là:")
	if errEmail != nil {
		return errors.New("không thể gửi email OTP")
	}
	return nil
}

func ConfirmNewPassword(email, otpCode, newPassword string) error {
	var validOTP models.OTP
	config.DB.Where("email = ? AND code = ? AND is_used = ? AND expires_at > ?", email, otpCode, false, time.Now()).First(&validOTP)
	if validOTP.ID == 0 {
		return errors.New("mã OTP không đúng hoặc đã hết hạn")
	}

	errValidate := utils.ValidateStrongPassword(newPassword)
	if errValidate != nil {
		return errValidate
	}

	hashedPassword, _ := utils.HashPassword(newPassword)

	errUpdate := config.DB.Model(&models.User{}).
		Where("email = ?", email).
		Updates(map[string]interface{}{"password_hash": hashedPassword, "is_first_login": false}).Error

	if errUpdate != nil {
		return errors.New("không thể cập nhật mật khẩu")
	}

	config.DB.Delete(&validOTP)

	return nil
}

func LoginUser(email string, password string) (string, string, bool, error) {
	var user models.User
	result := config.DB.Where("email = ?", email).First(&user)

	if result.Error != nil {
		return "", "", false, errors.New("email hoặc mật khẩu không chính xác")
	}

	if !user.Status {
		return "", "", false, errors.New("tài khoản chưa được kích hoạt, vui lòng kiểm tra email")
	}

	isPasswordValid := utils.CheckPasswordHash(password, user.PasswordHash)
	if !isPasswordValid {
		return "", "", false, errors.New("email hoặc mật khẩu không chính xác")
	}

	if user.IsFirstLogin {
		RequestPasswordChangeOTP(user.Email)
		return "", "", true, errors.New("yêu cầu đổi mật khẩu ở lần đăng nhập đầu tiên. OTP đã được gửi về email")
	}

	accessToken, refreshToken, errToken := utils.GenerateTokens(user.ID, user.Role)
	if errToken != nil {
		return "", "", false, errors.New("lỗi khi tạo phiên đăng nhập")
	}

	return accessToken, refreshToken, false, nil
}

func CleanupUnverifiedUsersAndOTP() {
	for {
		config.DB.Where("expires_at <= ? OR is_used = ?", time.Now(), true).Delete(&models.OTP{})

		expirationTime := time.Now().Add(-15 * time.Minute)
		config.DB.Where("status = ? AND created_at <= ?", false, expirationTime).Delete(&models.User{})

		time.Sleep(10 * time.Minute)
	}
}

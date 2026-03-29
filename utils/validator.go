package utils

import (
	"errors"
	"regexp"
)

func ValidateStrongPassword(password string) error {
	if len(password) < 8 {
		return errors.New("mật khẩu phải có ít nhất 8 ký tự")
	}

	hasUppercase := regexp.MustCompile(`[A-Z]`).MatchString(password)
	if !hasUppercase {
		return errors.New("mật khẩu phải chứa ít nhất 1 chữ cái in hoa")
	}

	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	if !hasNumber {
		return errors.New("mật khẩu phải chứa ít nhất 1 chữ số")
	}

	hasSpecialChar := regexp.MustCompile(`[!@#~$%^&*(),.?":{}|<>]`).MatchString(password)
	if !hasSpecialChar {
		return errors.New("mật khẩu phải chứa ít nhất 1 ký tự đặc biệt")
	}

	return nil
}

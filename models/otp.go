package models

import "time"

type OTP struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Email     string    `json:"email"`
	Code      string    `json:"code"`
	ExpiresAt time.Time `json:"expires_at"`
	IsUsed    bool      `json:"is_used" gorm:"default:false"`
}

func (OTP) TableName() string {
	return "otp"
}

package models

import "time"

type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	FullName     string    `json:"full_name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Phone        string    `json:"phone"`
	CCCD         string    `json:"cccd"`
	Role         string    `json:"role"`
	Status       bool      `json:"status"`
	IsFirstLogin bool      `json:"is_first_login" gorm:"default:false"`
	Avatar       string    `json:"avatar"`
	CreatedAt    time.Time `json:"created_at"`
	EmployerID   int       `json:"employer_id"`
}

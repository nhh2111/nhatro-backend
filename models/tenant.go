package models

import (
	"time"
)

type Tenant struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	FullName       string    `json:"full_name"`
	CCCD           string    `json:"cccd"`
	Phone          string    `json:"phone"`
	Dob            string    `json:"dob"`
	Gender         string    `json:"gender"`
	MotorbikeCount int       `json:"motorbike_count"`
	CarCount       int       `json:"car_count"`
	Address        string    `json:"address"`
	LicensePlates  string    `json:"license_plates"`
	ImageUrl       string    `json:"image_url"`
	CreatedAt      time.Time `json:"created_at"`
}

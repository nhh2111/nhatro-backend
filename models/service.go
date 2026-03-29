package models

type Service struct {
	ID          uint    `gorm:"primaryKey" json:"id"`
	Name        string  `json:"name"`
	ServiceType string  `json:"service_type"`
	UnitPrice   float64 `json:"unit_price"`
	Unit        string  `json:"unit"`
	Description string  `json:"description"`
}

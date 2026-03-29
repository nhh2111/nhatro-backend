package models

import "time"

type Invoice struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	ContractID  uint      `json:"contract_id"`
	MonthYear   string    `json:"month_year"`
	DueDate     time.Time `json:"due_date"`
	TotalAmount float64   `json:"total_amount"`
	PaidAmount  float64   `json:"paid_amount"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`

	Contract Contract      `gorm:"foreignKey:ContractID" json:"Contract"`
	Items    []InvoiceItem `gorm:"foreignKey:InvoiceID" json:"Items"`
}

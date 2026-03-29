package dto

import "time"

// CreateTransactionDTO dùng cho việc tạo Phiếu Thu/Chi thủ công
type CreateTransactionDTO struct {
	HouseID         uint      `json:"house_id" binding:"required"`
	RoomID          uint      `json:"room_id"`                 // Có thể rỗng nếu là chi phí chung của nhà
	Type            string    `json:"type" binding:"required"` // INCOME hoặc EXPENSE
	Category        string    `json:"category" binding:"required"`
	Amount          float64   `json:"amount" binding:"required"`
	TransactionDate time.Time `json:"transaction_date" binding:"required"`
	PayerPayeeName  string    `json:"payer_payee_name"`
	Description     string    `json:"description"`
}

// PayInvoiceDTO dùng khi khách đóng tiền phòng
type PayInvoiceDTO struct {
	InvoiceID   uint    `json:"invoice_id" binding:"required"`
	Amount      float64 `json:"amount" binding:"required"`
	PaymentDate string  `json:"payment_date"`
}

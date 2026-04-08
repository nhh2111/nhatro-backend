package dto

import "time"

type CreateTransactionDTO struct {
	HouseID         uint      `json:"house_id" binding:"required"`
	RoomID          uint      `json:"room_id"`
	Type            string    `json:"type" binding:"required"`
	Category        string    `json:"category" binding:"required"`
	Amount          float64   `json:"amount" binding:"required"`
	TransactionDate time.Time `json:"transaction_date" binding:"required"`
	PayerPayeeName  string    `json:"payer_payee_name"`
	Description     string    `json:"description"`
}

type PayInvoiceDTO struct {
	InvoiceID   uint    `json:"invoice_id" binding:"required"`
	Amount      float64 `json:"amount" binding:"required"`
	PaymentDate string  `json:"payment_date"`
}

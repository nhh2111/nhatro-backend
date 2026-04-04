package models

import "time"

type Transaction struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	HouseID         *uint     `json:"house_id"`
	RoomID          *uint     `json:"room_id"`
	Type            string    `json:"type"`
	Category        string    `json:"category"`
	Amount          float64   `json:"amount"`
	TransactionDate time.Time `json:"transaction_date"`
	PayerPayeeName  string    `json:"payer_payee_name"`
	Description     string    `json:"description"`

	House House `gorm:"foreignKey:HouseID" json:"House"`
	Room  Room  `gorm:"foreignKey:RoomID" json:"Room"`
}

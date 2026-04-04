package dto

import "time"

type CreateContractDTO struct {
	FullName      string    `json:"full_name"`
	CCCD          string    `json:"cccd"`
	Phone         string    `json:"phone"`
	RoomID        uint      `json:"room_id"`
	StartDate     time.Time `json:"start_date"`
	EndDate       time.Time `json:"end_date"`
	DepositAmount float64   `json:"deposit_amount"`
	Terms         string    `json:"terms"`
}

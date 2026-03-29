package dto

import "time"

type CreateMeterReadingDTO struct {
	RoomID      uint      `json:"room_id" binding:"required"`
	ServiceID   uint      `json:"service_id" binding:"required"`
	ReadingDate time.Time `json:"reading_date" binding:"required"`
	OldIndex    float64   `json:"old_index"`
	NewIndex    float64   `json:"new_index" binding:"required"`
}

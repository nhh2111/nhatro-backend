package models

type RoomService struct {
	RoomID    uint `gorm:"primaryKey" json:"room_id"`
	ServiceID uint `gorm:"primaryKey" json:"service_id"`
}

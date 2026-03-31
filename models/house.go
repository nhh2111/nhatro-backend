package models

type House struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	OwnerID  uint   `json:"owner_id"`
	Name     string `json:"name"`
	City     string `json:"city"`
	District string `json:"district"`
	Ward     string `json:"ward"`
	Address  string `json:"address"`

	TotalRooms int `gorm:"->" json:"total_rooms"`
	EmptyRooms int `gorm:"->" json:"empty_rooms"`
}

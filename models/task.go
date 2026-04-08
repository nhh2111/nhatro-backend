package models

import "time"

type Task struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	HouseID    uint       `json:"house_id"`
	RoomID     *uint      `json:"room_id"`
	Title      string     `json:"title"`
	Content    string     `json:"content"`
	ImageUrl   string     `json:"image_url"`
	Status     string     `json:"status"`
	CreatedAt  time.Time  `json:"created_at"`
	FinishedAt *time.Time `json:"finished_at"`
	OwnerID    uint       `json:"owner_id"`
	Assignee   string     `json:"assignee"`
	Cost       float64    `json:"cost"`

	House House `gorm:"foreignKey:HouseID" json:"House"`
	Room  Room  `gorm:"foreignKey:RoomID" json:"Room"`
}

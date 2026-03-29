package models

type MeterReading struct {
	ID           uint    `gorm:"primaryKey" json:"id"`
	RoomID       uint    `json:"room_id"`
	ServiceID    uint    `json:"service_id"`
	BillingMonth string  `json:"billing_month"`
	ReadingDate  string  `json:"reading_date"`
	OldIndex     float64 `json:"old_index"`
	NewIndex     float64 `json:"new_index"`
	UsageValue   float64 `json:"usage_value"`

	Room    Room    `gorm:"foreignKey:RoomID" json:"Room"`
	Service Service `gorm:"foreignKey:ServiceID" json:"Service"`
}

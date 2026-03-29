package models

type Room struct {
	ID                uint    `gorm:"primaryKey" json:"id"`
	HouseID           uint    `json:"house_id"`
	RoomNumber        string  `json:"room_number"`
	Floor             int     `json:"floor"`
	Width             float64 `json:"width"`
	Length            float64 `json:"length"`
	BasePrice         float64 `json:"base_price"`
	MaxOccupants      int     `json:"max_occupants"`
	GenderRestriction string  `json:"gender_restriction"`
	Status            string  `json:"status"`
	Description       string  `json:"description"`
	Images            string  `json:"images"`
	CurrentOccupants  int     `json:"current_occupants" gorm:"-"`
}

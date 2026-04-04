package models

type Contract struct {
	ID            uint    `gorm:"primaryKey" json:"id"`
	RoomID        uint    `json:"room_id"`
	TenantID      uint    `json:"tenant_id"`
	StartDate     string  `json:"start_date"`
	EndDate       string  `json:"end_date"`
	DepositAmount float64 `json:"deposit_amount"`
	Status        string  `json:"status"`
	Room          Room    `gorm:"foreignKey:RoomID" json:"Room"`
	Tenant        Tenant  `gorm:"foreignKey:TenantID" json:"Tenant"`
	Terms         string  `json:"terms"`
}

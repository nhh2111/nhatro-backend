package config

import (
	"doAnHTTT_go/models"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("Lỗi: Chưa cài đặt biến môi trường DB_DSN")
	}

	databaseConnection, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})

	if err != nil {
		log.Fatalf("Không thể kết nối đến cơ sở dữ liệu: %v", err)
	}

	DB = databaseConnection
	log.Println("Kết nối cơ sở dữ liệu thành công!")

	errMigrate := DB.AutoMigrate(
		&models.User{},
		&models.House{},
		&models.Room{},
		&models.Tenant{},
		&models.Contract{},
		&models.Service{},
		&models.MeterReading{},
		&models.Invoice{},
		&models.Transaction{},
		&models.Task{},
		&models.RoomService{},
		&models.OTP{},
		&models.InvoiceItem{},
	)
	if errMigrate != nil {
		log.Println("Cảnh báo: Lỗi khi AutoMigrate:", errMigrate)
	} else {
		log.Println("Tạo bảng thành công (AutoMigrate Done)!")
	}
}

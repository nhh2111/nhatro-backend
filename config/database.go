package config

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := "root:@tcp(127.0.0.1:3306)/rental_system?charset=utf8mb4&parseTime=True&loc=Local"

	databaseConnection, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Không thể kết nối đến cơ sở dữ liệu: %v", err)
	}

	DB = databaseConnection
	log.Println("Kết nối cơ sở dữ liệu thành công!")
}

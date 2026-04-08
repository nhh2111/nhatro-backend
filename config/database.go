package config

import (
	"doAnHTTT_go/models"
	"log"
	"os"
	"strconv"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("Lỗi: Chưa cài đặt biến môi trường DB_DSN")
	}

	databaseConnection, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   logger.Default.LogMode(logger.Error),
	})

	if err != nil {
		log.Fatalf("Không thể kết nối đến cơ sở dữ liệu: %v", err)
	}

	DB = databaseConnection
	log.Println("Kết nối cơ sở dữ liệu thành công!")

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("Không thể lấy sql.DB từ GORM: %v", err)
	}

	maxOpenConns := getEnvInt("DB_MAX_OPEN_CONNS", 25)
	maxIdleConns := getEnvInt("DB_MAX_IDLE_CONNS", 25)
	connMaxLifetime := getEnvDurationSeconds("DB_CONN_MAX_LIFETIME_SECONDS", 300)
	connMaxIdleTime := getEnvDurationSeconds("DB_CONN_MAX_IDLE_TIME_SECONDS", 60)

	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetConnMaxLifetime(connMaxLifetime)
	sqlDB.SetConnMaxIdleTime(connMaxIdleTime)

	if os.Getenv("AUTO_MIGRATE") == "true" {
		log.Println("Đang kiểm tra và cập nhật cấu trúc DB (AutoMigrate)...")
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
			log.Println("Cập nhật bảng thành công (AutoMigrate Done)!")
		}
	} else {
		log.Println("Bỏ qua AutoMigrate. Khởi động nhanh Server!")
	}
}

func getEnvInt(envName string, defaultValue int) int {
	raw := os.Getenv(envName)
	if raw == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(raw)
	if err != nil || value <= 0 {
		return defaultValue
	}
	return value
}

func getEnvDurationSeconds(envName string, defaultSeconds int) time.Duration {
	raw := os.Getenv(envName)
	if raw == "" {
		return time.Duration(defaultSeconds) * time.Second
	}

	seconds, err := strconv.Atoi(raw)
	if err != nil || seconds <= 0 {
		return time.Duration(defaultSeconds) * time.Second
	}
	return time.Duration(seconds) * time.Second
}

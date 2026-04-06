package services

import (
	"doAnHTTT_go/config"
	"doAnHTTT_go/dto"
	"doAnHTTT_go/models"
	"errors"
	"time"
)

func CreateNewMeterReading(ownerID uint, dtoInput dto.CreateMeterReadingDTO, userRole string) error {
	if dtoInput.NewIndex < dtoInput.OldIndex {
		return errors.New("chỉ số mới không được nhỏ hơn chỉ số cũ")
	}

	// KIỂM TRA BẢO MẬT: Phòng này có thuộc về nhà của Owner không?
	var roomCount int64
	config.DB.Table("rooms").Joins("JOIN houses ON rooms.house_id = houses.id").
		Where("rooms.id = ? AND houses.owner_id = ?", dtoInput.RoomID, ownerID).Count(&roomCount)
	if roomCount == 0 {
		return errors.New("phòng không hợp lệ hoặc bạn không có quyền thao tác")
	}

	if userRole == "STAFF" {
		currentYear, currentMonth, _ := time.Now().Date()
		inputYear, inputMonth, _ := dtoInput.ReadingDate.Date()

		isSameMonth := (currentYear == inputYear) && (currentMonth == inputMonth)
		if !isSameMonth {
			return errors.New("nhân viên chỉ được phép nhập chỉ số cho tháng hiện tại")
		}
	}

	newReading := models.MeterReading{
		RoomID:       dtoInput.RoomID,
		ServiceID:    dtoInput.ServiceID,
		BillingMonth: dtoInput.ReadingDate.Format("2006-01"),
		ReadingDate:  dtoInput.ReadingDate.Format("2006-01-02"),
		OldIndex:     dtoInput.OldIndex,
		NewIndex:     dtoInput.NewIndex,
		UsageValue:   dtoInput.NewIndex - dtoInput.OldIndex,
	}

	result := config.DB.Create(&newReading)
	if result.Error != nil {
		return errors.New("lỗi khi lưu chỉ số đồng hồ vào hệ thống")
	}
	return nil
}

func GetMeterReadingsByMonth(ownerID uint, month string) ([]models.MeterReading, error) {
	var readings []models.MeterReading

	// ĐÃ SỬA: Thay Preload("Room") thành Preload("Room.House") để lấy tên Khu trọ
	err := config.DB.Preload("Room.House").Preload("Service").
		Joins("JOIN rooms ON meter_readings.room_id = rooms.id").
		Joins("JOIN houses ON rooms.house_id = houses.id").
		Where("meter_readings.billing_month = ? AND houses.owner_id = ?", month, ownerID).
		Order("meter_readings.reading_date DESC").
		Find(&readings).Error

	if err != nil {
		return nil, errors.New("không thể lấy lịch sử ghi chỉ số")
	}
	return readings, nil
}

func UpdateMeterReading(ownerID uint, id uint, dtoInput dto.CreateMeterReadingDTO) error {
	var reading models.MeterReading

	err := config.DB.Joins("JOIN rooms ON meter_readings.room_id = rooms.id").
		Joins("JOIN houses ON rooms.house_id = houses.id").
		Where("meter_readings.id = ? AND houses.owner_id = ?", id, ownerID).
		First(&reading).Error

	if err != nil {
		return errors.New("không tìm thấy dữ liệu chỉ số này hoặc bạn không có quyền sửa")
	}

	if dtoInput.NewIndex < dtoInput.OldIndex {
		return errors.New("chỉ số mới không được nhỏ hơn chỉ số cũ")
	}

	reading.ReadingDate = dtoInput.ReadingDate.Format("2006-01-02")
	reading.BillingMonth = dtoInput.ReadingDate.Format("2006-01")
	reading.OldIndex = dtoInput.OldIndex
	reading.NewIndex = dtoInput.NewIndex
	reading.UsageValue = dtoInput.NewIndex - dtoInput.OldIndex

	if err := config.DB.Save(&reading).Error; err != nil {
		return errors.New("lỗi khi cập nhật chỉ số")
	}
	return nil
}

func DeleteMeterReading(ownerID uint, id uint) error {
	var reading models.MeterReading

	err := config.DB.Joins("JOIN rooms ON meter_readings.room_id = rooms.id").
		Joins("JOIN houses ON rooms.house_id = houses.id").
		Where("meter_readings.id = ? AND houses.owner_id = ?", id, ownerID).
		First(&reading).Error

	if err != nil {
		return errors.New("không tìm thấy dữ liệu chỉ số này hoặc bạn không có quyền xóa")
	}

	if err := config.DB.Delete(&reading).Error; err != nil {
		return errors.New("lỗi khi xóa chỉ số")
	}
	return nil
}

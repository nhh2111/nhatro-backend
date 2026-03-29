package services

import (
	"doAnHTTT_go/config"
	"doAnHTTT_go/dto"
	"doAnHTTT_go/models"
	"errors"
	"time"
)

func CreateNewMeterReading(dtoInput dto.CreateMeterReadingDTO, userRole string) error {
	if dtoInput.NewIndex < dtoInput.OldIndex {
		return errors.New("chỉ số mới không được nhỏ hơn chỉ số cũ")
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

		OldIndex:   dtoInput.OldIndex,
		NewIndex:   dtoInput.NewIndex,
		UsageValue: dtoInput.NewIndex - dtoInput.OldIndex,
	}

	result := config.DB.Create(&newReading)
	if result.Error != nil {
		return errors.New("lỗi khi lưu chỉ số đồng hồ vào hệ thống")
	}

	return nil
}
func GetMeterReadingsByMonth(month string) ([]models.MeterReading, error) {
	var readings []models.MeterReading

	err := config.DB.Preload("Room").Preload("Service").
		Where("billing_month = ?", month).
		Order("reading_date DESC").
		Find(&readings).Error

	if err != nil {
		return nil, errors.New("không thể lấy lịch sử ghi chỉ số")
	}

	return readings, nil
}

func UpdateMeterReading(id uint, dtoInput dto.CreateMeterReadingDTO) error {
	var reading models.MeterReading
	if err := config.DB.First(&reading, id).Error; err != nil {
		return errors.New("không tìm thấy dữ liệu chỉ số này")
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

func DeleteMeterReading(id uint) error {
	var reading models.MeterReading
	if err := config.DB.First(&reading, id).Error; err != nil {
		return errors.New("không tìm thấy dữ liệu chỉ số này")
	}

	if err := config.DB.Delete(&reading).Error; err != nil {
		return errors.New("lỗi khi xóa chỉ số")
	}
	return nil
}

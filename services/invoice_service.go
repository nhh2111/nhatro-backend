package services

import (
	"doAnHTTT_go/config"
	"doAnHTTT_go/dto"
	"doAnHTTT_go/models"
	"doAnHTTT_go/utils"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

func GetAllInvoices(ownerID uint, page int, pageSize int, monthYear string) (map[string]interface{}, error) {
	var invoiceList []models.Invoice
	var totalRecords int64

	query := config.DB.Model(&models.Invoice{}).
		Preload("Items").Preload("Contract.Room").Preload("Contract.Tenant").
		Joins("JOIN contracts ON invoices.contract_id = contracts.id").
		Joins("JOIN rooms ON contracts.room_id = rooms.id").
		Joins("JOIN houses ON rooms.house_id = houses.id").
		Where("houses.owner_id = ?", ownerID)

	if monthYear != "" {
		query = query.Where("invoices.month_year = ?", monthYear)
	}

	query.Count(&totalRecords)
	pageCount := utils.GetPageCount(totalRecords, pageSize)
	offset := utils.GetOffset(page, pageSize)

	result := query.Offset(offset).Limit(pageSize).Order("invoices.created_at DESC").Find(&invoiceList)
	if result.Error != nil {
		return nil, result.Error
	}

	return map[string]interface{}{
		"recordCount": totalRecords,
		"pageCount":   pageCount,
		"currentPage": page,
		"pageSize":    pageSize,
		"records":     invoiceList,
	}, nil
}

// Hàm GenerateInvoice giữ nguyên nội tại vì nó được gọi bởi AutoGenerateInvoices (đã được bọc bảo mật)
func GenerateInvoice(roomID uint, monthYear string) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		var room models.Room
		if err := tx.First(&room, roomID).Error; err != nil {
			return errors.New("không tìm thấy phòng")
		}

		var repContract models.Contract
		if err := tx.Where("room_id = ? AND status = 'ACTIVE'", roomID).Order("start_date ASC").First(&repContract).Error; err != nil {
			return errors.New("phòng không có hợp đồng nào đang hoạt động")
		}

		var existing models.Invoice
		tx.Where("contract_id = ? AND month_year = ?", repContract.ID, monthYear).First(&existing)
		if existing.ID != 0 {
			return nil
		}

		invoice := models.Invoice{
			ContractID:  repContract.ID,
			MonthYear:   monthYear,
			DueDate:     time.Now().AddDate(0, 0, 7),
			Status:      "UNPAID",
			TotalAmount: 0,
			PaidAmount:  0,
		}
		if err := tx.Create(&invoice).Error; err != nil {
			return err
		}

		var total float64 = 0

		rentItem := models.InvoiceItem{
			InvoiceID:   invoice.ID,
			Description: fmt.Sprintf("Tiền thuê phòng %s", room.RoomNumber),
			Quantity:    1,
			UnitPrice:   room.BasePrice,
			Amount:      room.BasePrice,
		}
		tx.Create(&rentItem)
		total += rentItem.Amount

		var roomServices []models.Service
		tx.Table("services").
			Joins("JOIN room_services ON services.id = room_services.service_id").
			Where("room_services.room_id = ? AND services.service_type IN ?", roomID, []string{"FIXED", "PER_PERSON", "PER_MOTORBIKE", "PER_CAR"}).
			Find(&roomServices)

		for _, s := range roomServices {
			quantity := 1.0
			if s.ServiceType == "PER_PERSON" {
				var personCount int64
				tx.Table("contracts").Where("room_id = ? AND status = 'ACTIVE'", roomID).Count(&personCount)
				quantity = float64(personCount)
			} else if s.ServiceType == "PER_MOTORBIKE" {
				var totalMotos int64
				tx.Table("tenants").Joins("JOIN contracts ON tenants.id = contracts.tenant_id").
					Where("contracts.room_id = ? AND contracts.status = 'ACTIVE'", roomID).Select("COALESCE(SUM(tenants.motorbike_count), 0)").Scan(&totalMotos)
				if totalMotos == 0 {
					continue
				}
				quantity = float64(totalMotos)
			} else if s.ServiceType == "PER_CAR" {
				var totalCars int64
				tx.Table("tenants").Joins("JOIN contracts ON tenants.id = contracts.tenant_id").
					Where("contracts.room_id = ? AND contracts.status = 'ACTIVE'", roomID).Select("COALESCE(SUM(tenants.car_count), 0)").Scan(&totalCars)
				if totalCars == 0 {
					continue
				}
				quantity = float64(totalCars)
			}

			item := models.InvoiceItem{
				InvoiceID:   invoice.ID,
				ServiceID:   &s.ID,
				Description: s.Name,
				Quantity:    quantity,
				UnitPrice:   s.UnitPrice,
				Amount:      quantity * s.UnitPrice,
			}
			tx.Create(&item)
			total += item.Amount
		}

		var readings []models.MeterReading
		tx.Preload("Service").Where("room_id = ? AND billing_month = ?", roomID, monthYear).Find(&readings)

		for _, r := range readings {
			item := models.InvoiceItem{
				InvoiceID:   invoice.ID,
				ServiceID:   &r.ServiceID,
				Description: fmt.Sprintf("%s (Số mới: %.1f - Số cũ: %.1f)", r.Service.Name, r.NewIndex, r.OldIndex),
				Quantity:    r.UsageValue,
				UnitPrice:   r.Service.UnitPrice,
				Amount:      r.UsageValue * r.Service.UnitPrice,
			}
			tx.Create(&item)
			total += item.Amount
		}
		return tx.Model(&invoice).Update("total_amount", total).Error
	})
}

func AutoGenerateInvoices(ownerID uint, monthYear string) error {
	var occupiedRooms []models.Room

	// CHỈ LẤY CÁC PHÒNG THUỘC SỞ HỮU CỦA CHỦ TRỌ NÀY
	config.DB.Joins("JOIN houses ON rooms.house_id = houses.id").
		Where("rooms.status = ? AND houses.owner_id = ?", "OCCUPIED", ownerID).
		Find(&occupiedRooms)

	if len(occupiedRooms) == 0 {
		return errors.New("không có phòng nào đang cho thuê để chốt sổ")
	}

	for _, room := range occupiedRooms {
		_ = GenerateInvoice(room.ID, monthYear)
	}
	return nil
}

func DeleteInvoice(ownerID uint, invoiceID uint, userRole string) error {
	var invoice models.Invoice

	err := config.DB.Joins("JOIN contracts ON invoices.contract_id = contracts.id").
		Joins("JOIN rooms ON contracts.room_id = rooms.id").
		Joins("JOIN houses ON rooms.house_id = houses.id").
		Where("invoices.id = ? AND houses.owner_id = ?", invoiceID, ownerID).
		First(&invoice).Error

	if err != nil {
		return errors.New("không tìm thấy hóa đơn hoặc bạn không có quyền xóa")
	}

	if userRole == "STAFF" && (invoice.Status == "PAID" || invoice.Status == "PARTIAL") {
		return errors.New("nhân viên không được xóa hóa đơn đã thu tiền")
	}

	config.DB.Where("invoice_id = ?", invoiceID).Delete(&models.InvoiceItem{})
	return config.DB.Delete(&invoice).Error
}

func PayInvoice(ownerID uint, dtoInput dto.PayInvoiceDTO) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		var invoice models.Invoice

		err := tx.Preload("Contract.Room").
			Joins("JOIN contracts ON invoices.contract_id = contracts.id").
			Joins("JOIN rooms ON contracts.room_id = rooms.id").
			Joins("JOIN houses ON rooms.house_id = houses.id").
			Where("invoices.id = ? AND houses.owner_id = ?", dtoInput.InvoiceID, ownerID).
			First(&invoice).Error

		if err != nil {
			return errors.New("không tìm thấy hóa đơn hoặc bạn không có quyền thao tác")
		}

		if invoice.Status == "PAID" {
			return errors.New("hóa đơn này đã được thanh toán đủ")
		}

		newPaidAmount := invoice.PaidAmount + dtoInput.Amount
		newStatus := "PARTIAL"
		if newPaidAmount >= invoice.TotalAmount {
			newStatus = "PAID"
		}

		if err := tx.Model(&invoice).Updates(map[string]interface{}{
			"paid_amount": newPaidAmount,
			"status":      newStatus,
		}).Error; err != nil {
			return errors.New("lỗi cập nhật hóa đơn")
		}

		houseID := invoice.Contract.Room.HouseID
		roomID := invoice.Contract.RoomID

		incomeTx := models.Transaction{
			HouseID:         &houseID,
			RoomID:          &roomID,
			Type:            "INCOME",
			Category:        "Tiền phòng",
			Amount:          dtoInput.Amount,
			TransactionDate: time.Now(),
			Description:     fmt.Sprintf("Thu tiền hóa đơn #%d - P.%s", invoice.ID, invoice.Contract.Room.RoomNumber),
		}

		if errTx := tx.Create(&incomeTx).Error; errTx != nil {
			return errors.New("lỗi tạo phiếu thu: " + errTx.Error())
		}
		return nil
	})
}

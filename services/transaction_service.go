package services

import (
	"doAnHTTT_go/config"
	"doAnHTTT_go/dto"
	"doAnHTTT_go/models"
	"doAnHTTT_go/utils"
	"errors"
)

func CreateNewTransaction(dtoInput dto.CreateTransactionDTO) error {
	isValidType := (dtoInput.Type == "INCOME" || dtoInput.Type == "EXPENSE")
	if !isValidType {
		return errors.New("loại giao dịch chỉ được phép là INCOME (Thu) hoặc EXPENSE (Chi)")
	}

	newTransaction := models.Transaction{
		HouseID:         &dtoInput.HouseID,
		RoomID:          &dtoInput.RoomID,
		Type:            dtoInput.Type,
		Category:        dtoInput.Category,
		Amount:          dtoInput.Amount,
		TransactionDate: dtoInput.TransactionDate,
		PayerPayeeName:  dtoInput.PayerPayeeName,
		Description:     dtoInput.Description,
	}

	result := config.DB.Create(&newTransaction)
	if result.Error != nil {
		return errors.New("lỗi khi lưu phiếu thu/chi vào cơ sở dữ liệu")
	}

	return nil
}
func GetAllTransactions(page int, pageSize int, search string) (map[string]interface{}, error) {
	var transactionList []models.Transaction
	var totalRecords int64

	query := config.DB.Model(&models.Transaction{}).Preload("House").Preload("Room")

	if search != "" {
		searchKeyword := "%" + search + "%"
		query = query.Where("type LIKE ? OR category LIKE ?", searchKeyword, searchKeyword)
	}
	query.Count(&totalRecords)

	pageCount := utils.GetPageCount(totalRecords, pageSize)
	offset := utils.GetOffset(page, pageSize)

	result := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&transactionList)
	if result.Error != nil {
		return nil, result.Error
	}
	return map[string]interface{}{
		"recordCount": totalRecords,
		"pageCount":   pageCount,
		"currentPage": page,
		"pageSize":    pageSize,
		"records":     transactionList,
	}, nil
}

func UpdateTransaction(transactionID uint, updatedData map[string]interface{}) error {
	var transaction models.Transaction
	errFind := config.DB.First(&transaction, transactionID).Error
	if errFind != nil {
		return errors.New("không tìm thấy phiếu thu/chi cần sửa")
	}

	errUpdate := config.DB.Model(&transaction).Updates(updatedData).Error
	if errUpdate != nil {
		return errors.New("lỗi khi cập nhật phiếu thu/chi")
	}
	return nil
}

func DeleteTransaction(transactionID uint) error {
	var transaction models.Transaction
	errFind := config.DB.First(&transaction, transactionID).Error
	if errFind != nil {
		return errors.New("không tìm thấy phiếu thu/chi cần xóa")
	}

	errDelete := config.DB.Delete(&transaction).Error
	if errDelete != nil {
		return errors.New("lỗi khi xóa phiếu thu/chi")
	}
	return nil
}

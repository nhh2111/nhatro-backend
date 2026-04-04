package services

import (
	"doAnHTTT_go/config"
	"doAnHTTT_go/dto"
	"doAnHTTT_go/models"
	"doAnHTTT_go/utils"
	"errors"
)

func CreateNewTransaction(ownerID uint, dtoInput dto.CreateTransactionDTO) error {
	isValidType := (dtoInput.Type == "INCOME" || dtoInput.Type == "EXPENSE")
	if !isValidType {
		return errors.New("loại giao dịch chỉ được phép là INCOME (Thu) hoặc EXPENSE (Chi)")
	}

	// KIỂM TRA BẢO MẬT: Nhà này có thuộc về ownerID không?
	var house models.House
	if err := config.DB.Where("id = ? AND owner_id = ?", dtoInput.HouseID, ownerID).First(&house).Error; err != nil {
		return errors.New("khu trọ không hợp lệ hoặc bạn không có quyền thêm giao dịch vào đây")
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

func GetAllTransactions(ownerID uint, page int, pageSize int, search string, monthYear string) (map[string]interface{}, error) {
	var transactionList []models.Transaction
	var totalRecords int64

	// JOIN SANG NHÀ ĐỂ LỌC THEO CHỦ
	query := config.DB.Model(&models.Transaction{}).
		Preload("House").Preload("Room").
		Joins("JOIN houses ON transactions.house_id = houses.id").
		Where("houses.owner_id = ?", ownerID)

	if search != "" {
		searchKeyword := "%" + search + "%"
		query = query.Where("(transactions.type LIKE ? OR transactions.category LIKE ? OR transactions.description LIKE ?)", searchKeyword, searchKeyword, searchKeyword)
	}

	// BỘ LỌC THEO THÁNG MỚI BỔ SUNG
	if monthYear != "" {
		// Sử dụng DATE_FORMAT để ép kiểu ngày giờ trong DB về dạng YYYY-MM chuẩn xác 100%
		query = query.Where("DATE_FORMAT(transactions.transaction_date, '%Y-%m') = ?", monthYear)
	}
	query.Count(&totalRecords)

	pageCount := utils.GetPageCount(totalRecords, pageSize)
	offset := utils.GetOffset(page, pageSize)

	result := query.Offset(offset).Limit(pageSize).Order("transactions.transaction_date DESC, transactions.id DESC").Find(&transactionList)
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

func UpdateTransaction(ownerID uint, transactionID uint, updatedData map[string]interface{}) error {
	var transaction models.Transaction

	// KIỂM TRA BẢO MẬT
	errFind := config.DB.Joins("JOIN houses ON transactions.house_id = houses.id").
		Where("transactions.id = ? AND houses.owner_id = ?", transactionID, ownerID).
		First(&transaction).Error

	if errFind != nil {
		return errors.New("không tìm thấy phiếu thu/chi cần sửa hoặc bạn không có quyền")
	}

	errUpdate := config.DB.Model(&transaction).Updates(updatedData).Error
	if errUpdate != nil {
		return errors.New("lỗi khi cập nhật phiếu thu/chi")
	}
	return nil
}

func DeleteTransaction(ownerID uint, transactionID uint) error {
	var transaction models.Transaction

	errFind := config.DB.Joins("JOIN houses ON transactions.house_id = houses.id").
		Where("transactions.id = ? AND houses.owner_id = ?", transactionID, ownerID).
		First(&transaction).Error

	if errFind != nil {
		return errors.New("không tìm thấy phiếu thu/chi cần xóa hoặc bạn không có quyền")
	}

	errDelete := config.DB.Delete(&transaction).Error
	if errDelete != nil {
		return errors.New("lỗi khi xóa phiếu thu/chi")
	}
	return nil
}

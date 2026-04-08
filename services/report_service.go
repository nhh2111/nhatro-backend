package services

import (
	"doAnHTTT_go/config"
	"doAnHTTT_go/dto"
	"errors"
)

func GetProfitLossReport(ownerID uint, month string, year string) (dto.ProfitLossReportDTO, error) {
	var totalIncome float64
	var totalExpense float64

	errIncome := config.DB.Table("transactions").
		Joins("JOIN houses ON transactions.house_id = houses.id").
		Where("houses.owner_id = ? AND transactions.type = ? AND MONTH(transactions.transaction_date) = ? AND YEAR(transactions.transaction_date) = ?", ownerID, "INCOME", month, year).
		Select("COALESCE(SUM(transactions.amount), 0)").Scan(&totalIncome).Error

	if errIncome != nil {
		return dto.ProfitLossReportDTO{}, errors.New("lỗi khi thống kê tổng thu")
	}

	errExpense := config.DB.Table("transactions").
		Joins("JOIN houses ON transactions.house_id = houses.id").
		Where("houses.owner_id = ? AND transactions.type = ? AND MONTH(transactions.transaction_date) = ? AND YEAR(transactions.transaction_date) = ?", ownerID, "EXPENSE", month, year).
		Select("COALESCE(SUM(transactions.amount), 0)").Scan(&totalExpense).Error

	if errExpense != nil {
		return dto.ProfitLossReportDTO{}, errors.New("lỗi khi thống kê tổng chi")
	}

	report := dto.ProfitLossReportDTO{
		MonthYear:    month + "/" + year,
		TotalIncome:  totalIncome,
		TotalExpense: totalExpense,
		NetProfit:    totalIncome - totalExpense,
	}

	return report, nil
}

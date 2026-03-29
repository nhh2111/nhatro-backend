package services

import (
	"doAnHTTT_go/config"
	"doAnHTTT_go/dto"
	"errors"
)

func GetProfitLossReport(month string, year string) (dto.ProfitLossReportDTO, error) {
	var totalIncome float64
	var totalExpense float64

	errIncome := config.DB.Table("transactions").
		Where("type = ? AND MONTH(transaction_date) = ? AND YEAR(transaction_date) = ?", "INCOME", month, year).
		Select("COALESCE(SUM(amount), 0)").Scan(&totalIncome).Error

	if errIncome != nil {
		return dto.ProfitLossReportDTO{}, errors.New("lỗi khi thống kê tổng thu")
	}

	errExpense := config.DB.Table("transactions").
		Where("type = ? AND MONTH(transaction_date) = ? AND YEAR(transaction_date) = ?", "EXPENSE", month, year).
		Select("COALESCE(SUM(amount), 0)").Scan(&totalExpense).Error

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

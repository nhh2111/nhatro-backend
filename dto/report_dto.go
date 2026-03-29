package dto

type ProfitLossReportDTO struct {
	MonthYear    string  `json:"month_year"`
	TotalIncome  float64 `json:"total_income"`
	TotalExpense float64 `json:"total_expense"`
	NetProfit    float64 `json:"net_profit"`
}

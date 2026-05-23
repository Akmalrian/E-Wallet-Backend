package dto

type DashboardResponse struct {
	Balance      float64 `json:"balance"`
	TotalIncome  float64 `json:"total_income"`
	TotalExpense float64 `json:"total_expense"`
}

package dto

// DashboardResponse — response dashboard info
type DashboardResponse struct {
	Balance      float64 `json:"balance"`
	TotalIncome  float64 `json:"total_income"`
	TotalExpense float64 `json:"total_expense"`
}

// GraphPoint — satu titik data di graph
type GraphPoint struct {
	Label   string  `json:"label"`   // ← label sumbu X (tanggal/bulan)
	Income  float64 `json:"income"`  // ← total topup
	Expense float64 `json:"expense"` // ← total transfer
}

// GraphResponse — response data graph
type GraphResponse struct {
	Points []GraphPoint `json:"points"`
	Filter GraphFilter  `json:"filter"`
}

// GraphFilter — info filter yang dipakai
type GraphFilter struct {
	Type      string `json:"type"`       // "income", "expense", "both"
	StartDate string `json:"start_date"` // "2024-01-01"
	EndDate   string `json:"end_date"`   // "2024-01-31"
}

// ── Swagger ──────────────────────────────────
type SwaggerDashboardResponse struct {
	Message string            `json:"message"`
	Success bool              `json:"success"`
	Data    DashboardResponse `json:"data"`
}

type SwaggerGraphResponse struct {
	Message string        `json:"message"`
	Success bool          `json:"success"`
	Data    GraphResponse `json:"data"`
}

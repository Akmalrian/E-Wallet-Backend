package dto

// ── Auth Swagger Response ─────────────────────────────

// SwaggerLoginResponse — response login untuk swagger docs
type SwaggerLoginResponse struct {
	Message string        `json:"message"`
	Success bool          `json:"success"`
	Data    LoginResponse `json:"data"`
}

// ── User Swagger Response ─────────────────────────────

// SwaggerProfileResponse — response get profile untuk swagger docs
type SwaggerProfileResponse struct {
	Message string             `json:"message"`
	Success bool               `json:"success"`
	Data    GetProfileResponse `json:"data"`
}

// SwaggerReceiversResponse — response find receivers untuk swagger docs
type SwaggerReceiversResponse struct {
	Message string               `json:"message"`
	Success bool                 `json:"success"`
	Data    ReceiverListResponse `json:"data"`
}

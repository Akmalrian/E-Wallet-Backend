package pkg

// Response — generic response untuk semua endpoint
type Response[T any] struct {
	Message string `json:"message"`
	Data    T      `json:"data"`
	Success bool   `json:"success"`
}

// ErrorResponse — response khusus error
type ErrorResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

// BaseResponse — response tanpa data
type BaseResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

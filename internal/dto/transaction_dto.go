package dto

import (
	"time"
)

// TopupBody — request body untuk topup
type TopupBody struct {
	PaymentMethodId int     `json:"payment_method_id" binding:"required"`
	OrderAmount     float64 `json:"order_amount"      binding:"required,gt=0"`
}

// TransferBody — request body untuk transfer
type TransferBody struct {
	ReceiverWalletId int     `json:"receiver_wallet_id" binding:"required"`
	Amount           float64 `json:"amount"             binding:"required,gt=0"`
	Pin              string  `json:"pin"                binding:"required,len=6"`
	Notes            string  `json:"notes"`
}

// TransactionResponse — response data transaksi
type TransactionResponse struct {
	Id        int       `json:"id"`
	Type      string    `json:"type"`
	Amount    float64   `json:"amount"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// TopupResponse — response detail topup
type TopupResponse struct {
	TransactionId int       `json:"transaction_id"`
	PaymentMethod string    `json:"payment_method"`
	OrderAmount   float64   `json:"order_amount"`
	TaxAmount     float64   `json:"tax_amount"`
	TotalAmount   float64   `json:"total_amount"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
}

// TransferResponse — response detail transfer
type TransferResponse struct {
	TransactionId    int       `json:"transaction_id"`
	ReceiverWalletId int       `json:"receiver_wallet_id"`
	Amount           float64   `json:"amount"`
	Notes            string    `json:"notes"`
	Status           string    `json:"status"`
	CreatedAt        time.Time `json:"created_at"`
}

// HistoryResponse — response list history transaksi
type HistoryResponse struct {
	Id           int           `json:"id"`
	Type         string        `json:"type"`
	Amount       float64       `json:"amount"`
	Notes        *string       `json:"notes"`
	Status       string        `json:"status"`
	SenderInfo   *SenderInfo   `json:"sender_info"`
	ReceiverInfo *ReceiverInfo `json:"receiver_info"`
	CreatedAt    time.Time     `json:"created_at"`
}

// SenderInfo — info pengirim untuk history penerima
type SenderInfo struct {
	WalletId    int     `json:"wallet_id"`
	Fullname    *string `json:"fullname"`
	PhoneNumber *string `json:"phone_number"`
	PhotoPath   *string `json:"photo_path"`
}

// ReceiverInfo — info penerima untuk history pengirim
type ReceiverInfo struct {
	WalletId    int     `json:"wallet_id"`
	Fullname    *string `json:"fullname"`
	PhoneNumber *string `json:"phone_number"`
	PhotoPath   *string `json:"photo_path"`
}

// HistoryListResponse — response list + pagination
type HistoryListResponse struct {
	Transactions []HistoryResponse `json:"transactions"`
	Meta         PaginationMeta    `json:"meta"`
}

type SwaggerTopupResponse struct {
	Message string        `json:"message"`
	Success bool          `json:"success"`
	Data    TopupResponse `json:"data"`
}

type SwaggerTransferResponse struct {
	Message string           `json:"message"`
	Success bool             `json:"success"`
	Data    TransferResponse `json:"data"`
}

type SwaggerHistoryResponse struct {
	Message string              `json:"message"`
	Success bool                `json:"success"`
	Data    HistoryListResponse `json:"data"`
}

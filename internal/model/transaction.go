package model

import "time"

type transfer struct {
	Id        int
	UserId    int
	Type      string
	Amount    float64
	Status    string
	CreatedAt time.Time
	updatedAt *time.Time
}

type TopupDetail struct {
	Id              int
	TrasnsactionId  int
	WalletId        int
	PaymentMethodId int
	OrderAmount     float64
	TaxAmount       float64
	DeliveryFee     float64
	TotalAmount     float64
	CreatedAt       time.Time
}

type TransferDetail struct {
	Id               int
	TrasnsactionId   int
	SenderWalletId   int
	ReceiverWalletId int
	Amount           float64
	Notes            *string
	CreatedAt        time.Time
}

package model

// PaymentReq represents data about card request
type PaymentReq struct {
	QRCode string  `json:"qr_code,omitempty" binding:"required" validate:"required"`
	Amount float64 `json:"amount,omitempty" binding:"required" validate:"required"`
}

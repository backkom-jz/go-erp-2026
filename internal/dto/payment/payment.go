package payment

type CallbackRequest struct {
	OrderNo   string `json:"order_no" binding:"required"`
	PaymentNo string `json:"payment_no" binding:"required"`
	Channel   string `json:"channel" binding:"required"`
	Status    string `json:"status" binding:"required"`
}

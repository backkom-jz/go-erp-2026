package order

type CreateOrderItem struct {
	SKUID      uint  `json:"sku_id" binding:"required"`
	Qty        int64 `json:"qty" binding:"required,gt=0"`
	PriceCents int64 `json:"price_cents" binding:"required,gt=0"`
}

type CreateOrderRequest struct {
	UserID   uint              `json:"user_id" binding:"required"`
	TenantID string            `json:"tenant_id" binding:"required"`
	Items    []CreateOrderItem `json:"items" binding:"required,min=1"`
}

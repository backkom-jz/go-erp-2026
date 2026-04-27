package product

type CreateSPURequest struct {
	Name       string `json:"name" binding:"required"`
	CategoryID uint   `json:"category_id"`
	Brand      string `json:"brand"`
}

type CreateSKURequest struct {
	SPUID      uint   `json:"spu_id" binding:"required"`
	Code       string `json:"code" binding:"required"`
	Name       string `json:"name" binding:"required"`
	PriceCents int64  `json:"price_cents" binding:"required"`
}

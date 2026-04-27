package inventory

type DeductRequest struct {
	SKUID      uint   `json:"sku_id" binding:"required"`
	Qty        int64  `json:"qty" binding:"required,gt=0"`
	BusinessNo string `json:"business_no" binding:"required"`
}

package request

type Product struct {
	SKU      string  `json:"sku" binding:"required"`
	Name     string  `json:"name" binding:"required"`
	Quantity uint64  `json:"quantity" binding:"required"`
	Price    float64 `json:"price" binding:"required"`
	Unit     string  `json:"unit" binding:"required"`
	Status   int     `json:"status" binding:"required"`
}

package dto

type CreateProductRequest struct {
	Name       string  `json:"name" binding:"required,min=2"`
	Price      float64 `json:"price" binding:"required,gt=0"`
	Stock      int     `json:"stock" binding:"gte=0"`
	CategoryID uint    `json:"category_id" binding:"required"`
	Image      string  `json:"image"`
}

type UpdateProductRequest struct {
	Name       string  `json:"name" binding:"required,min=2"`
	Price      float64 `json:"price" binding:"required,gt=0"`
	Stock      int     `json:"stock" binding:"gte=0"`
	CategoryID uint    `json:"category_id" binding:"required"`
	Image      string  `json:"image"`
}

type UpdateStockRequest struct {
	Stock int `json:"stock" binding:"gte=0"`
}

type ProductResponse struct {
	ID         uint    `json:"id"`
	Name       string  `json:"name"`
	Price      float64 `json:"price"`
	Stock      int     `json:"stock"`
	CategoryID uint    `json:"category_id"`
	Image      string  `json:"image,omitempty"`
}

type ProductListResponse struct {
	ID           uint    `json:"id"`
	Name         string  `json:"name"`
	Price        float64 `json:"price"`
	Stock        int     `json:"stock"`
	CategoryID   uint    `json:"category_id"`
	CategoryName string  `json:"category"`
	Image        string  `json:"image,omitempty"`
}

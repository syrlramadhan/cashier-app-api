package dto

import "time"

type TransactionItemRequest struct {
	ProductID uint `json:"product_id" binding:"required"`
	Quantity  int  `json:"quantity" binding:"required,gt=0"`
}

type CreateTransactionRequest struct {
	UserID        uint                     `json:"-"` // Set by controller from auth
	Items         []TransactionItemRequest `json:"items" binding:"required,min=1"`
	PaymentMethod string                   `json:"payment_method" binding:"required,oneof=cash card qris"`
}

type TransactionItemResponse struct {
	ID          uint    `json:"id"`
	ProductID   uint    `json:"product_id"`
	ProductName string  `json:"product_name"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
	Subtotal    float64 `json:"subtotal"`
}

type TransactionResponse struct {
	ID              uint                      `json:"id"`
	TransactionCode string                    `json:"transaction_code"`
	CashierName     string                    `json:"cashier_name"`
	Subtotal        float64                   `json:"subtotal"`
	Tax             float64                   `json:"tax"`
	Total           float64                   `json:"total"`
	PaymentMethod   string                    `json:"payment_method"`
	Status          string                    `json:"status"`
	Items           []TransactionItemResponse `json:"items"`
	CreatedAt       time.Time                 `json:"created_at"`
}

type TransactionListResponse struct {
	ID              uint      `json:"id"`
	TransactionCode string    `json:"transaction_code"`
	CashierName     string    `json:"cashier_name"`
	Total           float64   `json:"total"`
	PaymentMethod   string    `json:"payment_method"`
	ItemCount       int       `json:"item_count"`
	CreatedAt       time.Time `json:"created_at"`
}

type TransactionFilter struct {
	StartDate     string `form:"start_date"`
	EndDate       string `form:"end_date"`
	PaymentMethod string `form:"payment_method"`
	Page          int    `form:"page"`
	Limit         int    `form:"limit"`
}

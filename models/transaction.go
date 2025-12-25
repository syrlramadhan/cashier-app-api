package models

import (
	"time"

	"gorm.io/gorm"
)

type Transaction struct {
	ID              uint              `gorm:"primaryKey" json:"id"`
	TransactionCode string            `gorm:"size:50;uniqueIndex;not null" json:"transaction_code"`
	UserID          uint              `gorm:"not null" json:"user_id"`
	User            User              `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Subtotal        float64           `gorm:"not null" json:"subtotal"`
	TaxRate         float64           `gorm:"not null;default:0.11" json:"tax_rate"`
	Tax             float64           `gorm:"not null" json:"tax"`
	Total           float64           `gorm:"not null" json:"total"`
	PaymentMethod   string            `gorm:"size:20;not null" json:"payment_method"` // cash, card, qris
	Status          string            `gorm:"size:20;default:'completed'" json:"status"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
	DeletedAt       gorm.DeletedAt    `gorm:"index" json:"-"`
	Items           []TransactionItem `gorm:"foreignKey:TransactionID" json:"items,omitempty"`
}

func (Transaction) TableName() string {
	return "transactions"
}

type TransactionItem struct {
	ID            uint        `gorm:"primaryKey" json:"id"`
	TransactionID uint        `gorm:"not null" json:"transaction_id"`
	Transaction   Transaction `gorm:"foreignKey:TransactionID" json:"-"`
	ProductID     uint        `gorm:"not null" json:"product_id"`
	Product       Product     `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	ProductName   string      `gorm:"size:150;not null" json:"product_name"`
	Price         float64     `gorm:"not null" json:"price"`
	Quantity      int         `gorm:"not null" json:"quantity"`
	Subtotal      float64     `gorm:"not null" json:"subtotal"`
	CreatedAt     time.Time   `json:"created_at"`
}

func (TransactionItem) TableName() string {
	return "transaction_items"
}

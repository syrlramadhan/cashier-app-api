package models

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	Name       string         `gorm:"size:150;not null" json:"name"`
	Price      float64        `gorm:"not null" json:"price"`
	Stock      int            `gorm:"not null;default:0" json:"stock"`
	CategoryID uint           `gorm:"not null" json:"category_id"`
	Category   Category       `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Image      string         `gorm:"size:255" json:"image,omitempty"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Product) TableName() string {
	return "products"
}

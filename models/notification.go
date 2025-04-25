package models

import "time"

type Notification struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	ProductID int       `json:"product_id"`
	Product   Product   `json:"product" gorm:"foreignKey:ProductID"`
	Message   string    `json:"message"`
	Type      string    `json:"type"` // e.g., "price_drop", "back_in_stock"
	Sent      bool      `json:"sent" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

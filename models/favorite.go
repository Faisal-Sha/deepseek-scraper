package models

import "time"

type Favorite struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    string    `json:"user_id"`
	ProductID int       `json:"product_id"`
	Product   Product   `json:"product" gorm:"foreignKey:ProductID"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

package models

type Variant struct {
	ID        uint   `gorm:"primaryKey"`
	ProductID uint   `gorm:"index"`
	SKU       string `gorm:"uniqueIndex"`
	Name      string
	Price     float64
	Stock     int
	Available bool
}

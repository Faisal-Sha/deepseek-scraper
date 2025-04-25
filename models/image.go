package models

type Image struct {
    ID  uint   `gorm:"primarykey"`
    URL string `gorm:"not null"`
}
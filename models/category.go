package models

type Category struct {
    ID         string     `json:"id" gorm:"primaryKey"`
    Name       string     `json:"name"`
    URL        string     `json:"url"`
    ParentID   *string    `json:"parent_id"`
    Parent     *Category  `json:"-" gorm:"foreignKey:ParentID"`
    Children   []Category `json:"children,omitempty" gorm:"-"`
    IsLeaf     bool       `json:"is_leaf"`
    ProductCount int      `json:"product_count"`
}
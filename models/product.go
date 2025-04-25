// models/product.go
package models

import (
	"strings"
	"time"
)

// CustomTime is a custom time type that can handle time strings without timezone
type CustomTime time.Time

// UnmarshalJSON implements json.Unmarshaler interface
func (ct *CustomTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`) // Remove quotes
	if s == "null" || s == "" {
		*ct = CustomTime(time.Time{})
		return nil
	}

	// Try parsing with timezone first
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		// If that fails, try without timezone and set to UTC
		t, err = time.Parse("2006-01-02T15:04:05", s)
		if err != nil {
			return err
		}
	}

	*ct = CustomTime(t.UTC())
	return nil
}

// Time returns the time.Time representation
func (ct CustomTime) Time() time.Time {
	return time.Time(ct)
}

type Product struct {
    ID          int       `json:"id" gorm:"primaryKey"`
    Name        string    `json:"name"`
    URL         string    `json:"url"`
    Brand       string    `json:"brand"`
    BrandID     int       `json:"brandId"`
    MerchantID  int       `json:"merchantId"`
    CategoryID  int       `json:"categoryId"`
    ImageURL    string    `json:"image"`
    Rating      Rating    `json:"ratingScore" gorm:"embedded"`
    Price       Price     `json:"price" gorm:"embedded"`
    Promotions  []Promotion `json:"promotions" gorm:"serializer:json"`
    SocialProof []SocialProof `json:"socialProof" gorm:"serializer:json"`
    IsActive    bool      `json:"isActive" gorm:"default:true"`
    CreatedAt   time.Time `json:"createdAt" gorm:"autoCreateTime"`
    UpdatedAt   time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}

type Rating struct {
    AverageRating float64 `json:"averageRating"`
    TotalCount    int     `json:"totalCount"`
}

type Price struct {
    SellingPrice     float64 `json:"sellingPrice"`
    DiscountedPrice  float64 `json:"discountedPrice"`
    OriginalPrice    float64 `json:"originalPrice"`
    Currency         string  `json:"currency"`
}

type Promotion struct {
    ID              int        `json:"id"`
    Name            string     `json:"name"`
    DiscountType    int        `json:"discountType"`
    PromotionEndDate CustomTime `json:"promotionEndDate"`
}

type SocialProof struct {
    Key   string `json:"key"`
    Value string `json:"value"`
}

type PriceHistory struct {
    ID        uint      `gorm:"primaryKey"`
    ProductID string    `json:"productId"`
    Price     float64   `json:"price"`
    RecordedAt time.Time `json:"recordedAt" gorm:"autoCreateTime"`
}
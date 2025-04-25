package scraper

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"time"
	"trendyol-scraper/models"
)

type MockProcessor struct {
	filePath string
}

func NewMockProcessor(filePath string) *MockProcessor {
	return &MockProcessor{filePath: filePath}
}

func (mp *MockProcessor) ProcessMockData() ([]models.Product, error) {
	data, err := ioutil.ReadFile(mp.filePath)
	if err != nil {
		return nil, err
	}

	var mockData struct {
		Data struct {
			Contents []models.Product `json:"contents"`
		} `json:"data"`
	}

	if err := json.Unmarshal(data, &mockData); err != nil {
		return nil, err
	}

	// Process products to match our database model
	var products []models.Product
	for _, p := range mockData.Data.Contents {
		// Ensure CreatedAt and UpdatedAt have timezone information
		if p.CreatedAt.IsZero() {
			p.CreatedAt = time.Now()
		} else {
			// Add UTC timezone if not present
			p.CreatedAt = time.Date(
				p.CreatedAt.Year(),
				p.CreatedAt.Month(),
				p.CreatedAt.Day(),
				p.CreatedAt.Hour(),
				p.CreatedAt.Minute(),
				p.CreatedAt.Second(),
				p.CreatedAt.Nanosecond(),
				time.UTC,
			)
		}

		if p.UpdatedAt.IsZero() {
			p.UpdatedAt = p.CreatedAt
		} else {
			// Add UTC timezone if not present
			p.UpdatedAt = time.Date(
				p.UpdatedAt.Year(),
				p.UpdatedAt.Month(),
				p.UpdatedAt.Day(),
				p.UpdatedAt.Hour(),
				p.UpdatedAt.Minute(),
				p.UpdatedAt.Second(),
				p.UpdatedAt.Nanosecond(),
				time.UTC,
			)
		}

		// Note: No need to process promotions as CustomTime already handles timezone

		products = append(products, p)
	}

	log.Printf("Processed %d products from mock data", len(products))
	return products, nil
}
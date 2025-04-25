package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
	"trendyol-scraper/config"
	"trendyol-scraper/models"
)

type JSONStorage struct {
	outputPath string
}

func NewJSONStorage(cfg *config.Config) *JSONStorage {
	return &JSONStorage{outputPath: cfg.Scraper.JSONOutputPath}
}

func (js *JSONStorage) SaveCategories(categories []models.Category) error {
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := filepath.Join(js.outputPath, fmt.Sprintf("categories_%s.json", timestamp))

	if err := os.MkdirAll(js.outputPath, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create JSON file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(categories); err != nil {
		return fmt.Errorf("failed to encode categories to JSON: %w", err)
	}

	return nil
}

func (js *JSONStorage) SaveProducts(products []models.Product) error {
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := filepath.Join(js.outputPath, fmt.Sprintf("products_%s.json", timestamp))

	if err := os.MkdirAll(js.outputPath, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create JSON file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(products); err != nil {
		return fmt.Errorf("failed to encode products to JSON: %w", err)
	}

	return nil
}

func (js *JSONStorage) SaveVariants(variants []models.Variant) error {
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := filepath.Join(js.outputPath, fmt.Sprintf("variants_%s.json", timestamp))

	if err := os.MkdirAll(js.outputPath, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create JSON file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(variants); err != nil {
		return fmt.Errorf("failed to encode variants to JSON: %w", err)
	}

	return nil
}

func (js *JSONStorage) SaveImages(images []string) error {
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := filepath.Join(js.outputPath, fmt.Sprintf("images_%s.json", timestamp))

	if err := os.MkdirAll(js.outputPath, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create JSON file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(images); err != nil {
		return fmt.Errorf("failed to encode images to JSON: %w", err)
	}

	return nil
}

func (js *JSONStorage) GetProduct(id int) (*models.Product, error) {
	return nil, fmt.Errorf("GetProduct not implemented for JSON storage")
}

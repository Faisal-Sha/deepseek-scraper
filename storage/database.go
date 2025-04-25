package storage

import (
	"fmt"
	"trendyol-scraper/models"
	"gorm.io/gorm"
)

// StorageHandler defines the interface for storage operations
type StorageHandler interface {
	SaveCategories(categories []models.Category) error
	SaveProducts(products []models.Product) error
	SaveVariants(variants []models.Variant) error
	SaveImages(images []string) error
	GetProduct(id int) (*models.Product, error)
}

type DatabaseStorage struct {
	db *gorm.DB
}

func NewDatabaseStorage(db *gorm.DB) (*DatabaseStorage, error) {

	// Auto migrate models
	if err := db.AutoMigrate(&models.Category{}, &models.Product{}, &models.Variant{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return &DatabaseStorage{db: db}, nil
}

func (ds *DatabaseStorage) SaveCategories(categories []models.Category) error {
	return ds.db.Transaction(func(tx *gorm.DB) error {
		for _, cat := range categories {
			if err := tx.Save(&cat).Error; err != nil {
				return err
			}
			if len(cat.Children) > 0 {
				if err := ds.SaveCategories(cat.Children); err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func (s *DatabaseStorage) GetProduct(id int) (*models.Product, error) {
	var product models.Product
	if err := s.db.First(&product, id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (ds *DatabaseStorage) SaveProducts(products []models.Product) error {
	return ds.db.Transaction(func(tx *gorm.DB) error {
		for _, p := range products {
			if err := tx.Save(&p).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (ds *DatabaseStorage) SaveVariants(variants []models.Variant) error {
	return ds.db.Transaction(func(tx *gorm.DB) error {
		for _, v := range variants {
			if err := tx.Save(&v).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (ds *DatabaseStorage) SaveImages(images []string) error {
	return ds.db.Transaction(func(tx *gorm.DB) error {
		for _, img := range images {
			if err := tx.Create(&models.Image{URL: img}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

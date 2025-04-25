package main

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"strconv"
	"time"
	"trendyol-scraper/models"
	"trendyol-scraper/storage"

	"github.com/IBM/sarama"
	"gorm.io/gorm"
)

type ProductAnalysisService struct {
	db             *gorm.DB
	storageHandler storage.StorageHandler
	kafkaProducer  sarama.SyncProducer
}

func (s *ProductAnalysisService) ProcessProducts(ctx context.Context, products []models.Product) error {
	for _, product := range products {
		// Check if product exists
		var existingProduct models.Product
		result := s.db.Where("id = ?", product.ID).First(&existingProduct)

		if result.Error == gorm.ErrRecordNotFound {
			// New product
			if err := s.storageHandler.SaveProducts([]models.Product{product}); err != nil {
				log.Printf("Failed to find existing product %d: %v", product.ID, err)
				continue
			}
			log.Printf("New product inserted: %s", product.Name)
		} else if result.Error == nil {
			// Existing product - check for price changes
			if existingProduct.Price.DiscountedPrice != product.Price.DiscountedPrice {
				// Price changed - log and check for drop
				priceHistory := models.PriceHistory{
					ProductID: strconv.Itoa(product.ID),
					Price:     product.Price.DiscountedPrice,
				}
				
				if err := s.db.Create(&priceHistory).Error; err != nil {
					log.Printf("Failed to log price history for product %d: %v", product.ID, err)
				}

				// Check if price dropped
				if product.Price.DiscountedPrice < existingProduct.Price.DiscountedPrice {
					// Get users who favorited this product
					var favoriteUsers []models.Favorite
					if err := s.db.Where("product_id = ?", product.ID).Find(&favoriteUsers).Error; err != nil {
						log.Printf("Failed to get favorite users for product %d: %v", product.ID, err)
						continue
					}

					if len(favoriteUsers) > 0 {
						// Send price drop notification
						s.sendPriceDropNotification(favoriteUsers, product, existingProduct.Price.DiscountedPrice)
					}
				}
			}

			// Update product
			if err := s.db.Save(&product).Error; err != nil {
				log.Printf("Failed to scrape product %d: %v", product.ID, err)
			}
		} else {
			log.Printf("Error checking product existence: %v", result.Error)
		}
	}
	return nil
}

func (s *ProductAnalysisService) sendPriceDropNotification(users []models.Favorite, product models.Product, oldPrice float64) {
	// Prepare Kafka message
	message := PriceDropMessage{
		ProductID:   product.ID,
		ProductName: product.Name,
		OldPrice:    oldPrice,
		NewPrice:    product.Price.DiscountedPrice,
		Currency:    product.Price.Currency,
		ImageURL:    product.ImageURL,
		UserIDs:     make([]string, len(users)),
	}

	for i, fav := range users {
		message.UserIDs[i] = fav.UserID
	}

	// Send to Kafka
	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("Failed to marshal price drop message: %v", err)
		return
	}

	msg := &sarama.ProducerMessage{
		Topic: "price-drops",
		Value: sarama.StringEncoder(messageBytes),
	}

	if _, _, err := s.kafkaProducer.SendMessage(msg); err != nil {
		log.Printf("Failed to send price drop notification: %v", err)
	}
}

func (s *ProductAnalysisService) PrioritizeFavoritedProducts() {
	ticker := time.NewTicker(1 * time.Minute) // Check every minute
	for range ticker.C {
		// Get all favorited products
		var favorites []models.Favorite
		if err := s.db.Find(&favorites).Error; err != nil {
			log.Printf("Failed to get favorited products: %v", err)
			continue
		}

		// Create a map of product IDs to prioritize
		productIDs := make(map[int]bool)
		for _, fav := range favorites {
			productIDs[fav.ProductID] = true
		}

		// Process these products with higher priority
		for productID := range productIDs {
			// In a real implementation, we would fetch fresh data for this product
			// For mock data, we'll just update the existing record
			var product models.Product
			if err := s.db.Where("id = ?", productID).First(&product).Error; err != nil {
				log.Printf("Failed to find product %d: %v", productID, err)
				continue
			}

			// Simulate a potential price change
			// In real implementation, this would come from fresh scraping
			if rand.Intn(10) == 0 { // 10% chance of price change for demo
				oldPrice := product.Price.DiscountedPrice
				product.Price.DiscountedPrice = oldPrice * (0.9 + rand.Float64()*0.2) // Random price change ±10%
				
				if err := s.ProcessProducts(context.Background(), []models.Product{product}); err != nil {
					log.Printf("Failed to process prioritized product %d: %v", productID, err)
				}
			}
		}
	}
}
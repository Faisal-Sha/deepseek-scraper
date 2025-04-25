package main

import (
	"encoding/json"
	"fmt"
	"log"
	"trendyol-scraper/models"

	"github.com/IBM/sarama"
	"gorm.io/gorm"
)

type NotificationService struct {
	db *gorm.DB
}

type PriceDropMessage struct {
	ProductID   int      `json:"productId"`
	ProductName string   `json:"productName"`
	OldPrice    float64  `json:"oldPrice"`
	NewPrice    float64  `json:"newPrice"`
	Currency    string   `json:"currency"`
	ImageURL    string   `json:"imageUrl"`
	UserIDs     []string `json:"userIds"`
}

func (ns *NotificationService) StartConsumer() {
	config := sarama.NewConfig()
	consumer, err := sarama.NewConsumer([]string{"localhost:9092"}, config)
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}

	partitionConsumer, err := consumer.ConsumePartition("price-drops", 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Failed to create partition consumer: %v", err)
	}

	for {
		select {
		case msg := <-partitionConsumer.Messages():
			var priceDrop PriceDropMessage
			if err := json.Unmarshal(msg.Value, &priceDrop); err != nil {
				log.Printf("Failed to parse price drop message: %v", err)
				continue
			}

			if err := ns.sendNotifications(priceDrop); err != nil {
				log.Printf("Failed to send notifications: %v", err)
			}
		case err := <-partitionConsumer.Errors():
			log.Printf("Consumer error: %v", err)
		}
	}
}

func (ns *NotificationService) sendNotifications(msg PriceDropMessage) error {
	for _, userID := range msg.UserIDs {
		// In a real implementation, we would:
		// 1. Get user's notification preferences
		// 2. Send push notification/email/SMS based on preferences
		// 3. Log the notification

		log.Printf("Sending price drop notification to user %s for product %s (%s%.2f -> %s%.2f)",
			userID, msg.ProductName, msg.Currency, msg.OldPrice, msg.Currency, msg.NewPrice)

		// Log notification in database
		notification := models.Notification{
			ProductID: msg.ProductID, // Already an int
			Message:   fmt.Sprintf("Price dropped from %s%.2f to %s%.2f", msg.Currency, msg.OldPrice, msg.Currency, msg.NewPrice),
			Type:      "price_drop",
			Sent:      true,
		}

		if err := ns.db.Create(&notification).Error; err != nil {
			return err
		}
	}
	return nil
}
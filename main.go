package main

import (
	"context"
	"fmt"
	"log"
	"trendyol-scraper/config"
	"trendyol-scraper/models"
	"trendyol-scraper/scraper"
	"trendyol-scraper/storage"

	"github.com/IBM/sarama"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	dsn := buildDSN(cfg)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}

	// Drop existing tables
	if err := db.Migrator().DropTable(
		&models.Favorite{},
		&models.Notification{},
		&models.PriceHistory{},
		&models.Product{},
	); err != nil {
		log.Printf("Warning: Failed to drop tables: %v", err)
	}

	// Auto migrate models
	if err := db.AutoMigrate(
		&models.Product{},
		&models.PriceHistory{},
		&models.Notification{},
		&models.Favorite{},
	); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Initialize storage handler
	var storageHandler storage.StorageHandler
	if cfg.Scraper.OutputFormat == "db" {
		var err error
		storageHandler, err = storage.NewDatabaseStorage(db)
		if err != nil {
			log.Fatalf("Failed to initialize database storage: %v", err)
		}
	} else {
		storageHandler = storage.NewJSONStorage(cfg)
	}

	// Initialize Kafka producer (for price drop notifications)
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Producer.RequiredAcks = sarama.WaitForAll
	kafkaConfig.Producer.Return.Successes = true
	kafkaConfig.Net.MaxOpenRequests = 1

	kafkaProducer, err := sarama.NewSyncProducer([]string{"localhost:9092"}, kafkaConfig)
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
	}
	defer kafkaProducer.Close()

	// Initialize services
	productAnalysisSvc := &ProductAnalysisService{
		db:             db,
		storageHandler: storageHandler,
		kafkaProducer:  kafkaProducer,
	}

	// Process mock data
	mockProcessor := scraper.NewMockProcessor("data.json")
	products, err := mockProcessor.ProcessMockData()
	if err != nil {
		log.Fatalf("Failed to process mock data: %v", err)
	}

	// Analyze products
	if err := productAnalysisSvc.ProcessProducts(context.Background(), products); err != nil {
		log.Fatalf("Failed to process products: %v", err)
	}

	// Start notification service (in a separate goroutine)
	notificationSvc := &NotificationService{db: db}
	go notificationSvc.StartConsumer()

	// Keep main running
	select {}
}

func buildDSN(cfg *config.Config) string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.Database.Host, cfg.Database.User, cfg.Database.Password, cfg.Database.Name, cfg.Database.Port)
}

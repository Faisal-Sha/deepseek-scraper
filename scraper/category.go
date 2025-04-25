package scraper

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/chromedp/chromedp"
	"trendyol-scraper/models"
	"trendyol-scraper/config"
)

type CategoryScraper struct {
	config *config.Config
}

func NewCategoryScraper(cfg *config.Config) *CategoryScraper {
	return &CategoryScraper{config: cfg}
}

func (cs *CategoryScraper) ScrapeCategories() ([]models.Category, error) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.UserAgent(cs.config.Scraper.UserAgent),
	)

	ctx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel = chromedp.NewContext(ctx, chromedp.WithLogf(log.Printf))
	defer cancel()

	// Add timeout to context
	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	var categories []models.Category
	baseURL := cs.config.Scraper.BaseURL

	// Navigate to main page and extract top-level categories
	log.Printf("Navigating to %s", baseURL)
	err := chromedp.Run(ctx,
		chromedp.Navigate(baseURL),
		chromedp.WaitVisible(`nav`),
		chromedp.Sleep(2*time.Second),
		chromedp.Evaluate(`
			Array.from(document.querySelectorAll('nav a')).map(el => ({
				name: el.innerText.trim(),
				url: el.href
			})).filter(cat => cat.name && cat.url)
		`, &categories),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to scrape top-level categories: %w", err)
	}

	log.Printf("Found %d top-level categories", len(categories))

	// Scrape subcategories for each top-level category
	for i := range categories {
		if err := cs.scrapeSubcategories(ctx, &categories[i], 1); err != nil {
			log.Printf("Warning: failed to scrape subcategories for %s: %v", categories[i].Name, err)
			continue
		}
		time.Sleep(time.Duration(cs.config.Scraper.DelaySeconds) * time.Second)
	}

	return categories, nil
}

func (cs *CategoryScraper) scrapeSubcategories(ctx context.Context, parent *models.Category, depth int) error {
	if depth > cs.config.Scraper.MaxDepth {
		parent.IsLeaf = true
		return nil
	}

	log.Printf("Scraping subcategories for %s at depth %d", parent.Name, depth)
	time.Sleep(time.Duration(cs.config.Scraper.DelaySeconds) * time.Second)

	// Create a new context with timeout for this subcategory
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	var subcategories []models.Category
	err := chromedp.Run(ctx,
		chromedp.Navigate(parent.URL),
		chromedp.WaitVisible(`.sub-category-header`),
		chromedp.Evaluate(`
			Array.from(document.querySelectorAll('.sub-category-header')).map(el => ({
				name: el.innerText.trim(),
				url: el.href
			}))
		`, &subcategories),
	)
	if err != nil {
		return fmt.Errorf("failed to scrape subcategories: %w", err)
	}

	if len(subcategories) == 0 {
		parent.IsLeaf = true
		return nil
	}

	// Set parent reference
	for i := range subcategories {
		subcategories[i].ParentID = &parent.ID
	}

	// Recursively scrape deeper levels
	for i := range subcategories {
		if err := cs.scrapeSubcategories(ctx, &subcategories[i], depth+1); err != nil {
			return err
		}
	}

	parent.Children = subcategories
	return nil
}
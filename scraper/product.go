package scraper

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/chromedp/chromedp"
	"trendyol-scraper/models"
	"trendyol-scraper/config"
)

type ProductScraper struct {
	config *config.Config
}

func NewProductScraper(cfg *config.Config) *ProductScraper {
	return &ProductScraper{config: cfg}
}

func (ps *ProductScraper) ScrapeProductsFromCategory(categoryURL string) ([]models.Product, error) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var products []models.Product
	page := 1

	for {
		url := fmt.Sprintf("%s?pi=%d", categoryURL, page)
		var productLinks []string

		// Extract product links from page
		err := chromedp.Run(ctx,
			chromedp.Navigate(url),
			chromedp.WaitVisible(`.product-card`),
			chromedp.Evaluate(`
				Array.from(document.querySelectorAll('.product-card a')).map(a => a.href)
			`, &productLinks),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scrape product links: %w", err)
		}

		if len(productLinks) == 0 {
			break // No more products
		}

		// Scrape each product
		for _, link := range productLinks {
			time.Sleep(time.Duration(ps.config.Scraper.DelaySeconds) * time.Second)
			
			product, err := ps.scrapeProductPage(ctx, link)
			if err != nil {
				log.Printf("Failed to scrape product %s: %v", link, err)
				continue
			}

			products = append(products, *product)
		}

		page++
	}

	return products, nil
}

func (ps *ProductScraper) scrapeProductPage(ctx context.Context, url string) (*models.Product, error) {
	var product models.Product
	var rawData map[string]interface{}

	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitVisible(`.product-detail`),
		chromedp.Evaluate(`
			(() => {
				const script = document.querySelector('script[type="application/ld+json"]');
				return script ? JSON.parse(script.innerText) : {};
			})()
		`, &rawData),
		chromedp.Evaluate(`
			(() => {
				const product = {
					name: document.querySelector('.pr-new-br span')?.innerText.trim(),
					brand: document.querySelector('.merchant-text')?.innerText.trim(),
					price: parseFloat(document.querySelector('.prc-dsc')?.innerText.replace('TL', '').trim().replace('.', '').replace(',', '.')),
					originalPrice: parseFloat(document.querySelector('.prc-org')?.innerText.replace('TL', '').trim().replace('.', '').replace(',', '.')) || null,
					rating: parseFloat(document.querySelector('.rating-line')?.style.width || '0') / 20,
					images: Array.from(document.querySelectorAll('.gallery-modal-content img')).map(img => img.src),
					description: document.querySelector('.detail-attr-container')?.innerText.trim()
				};
				
				// Extract variants if available
				const variants = [];
				const variantElements = document.querySelectorAll('.variant-selector-item');
				variantElements.forEach(el => {
					variants.push({
						name: el.innerText.trim(),
						price: product.price, // Default to main price
						stock: 1 // Default available
					});
				});
				
				product.variants = variants;
				return product;
			})()
		`, &product),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to scrape product page: %w", err)
	}

	// Extract ID from URL
	re := regexp.MustCompile(`-p-(\d+)`)
	matches := re.FindStringSubmatch(url)
	if len(matches) > 1 {
		id, err := strconv.Atoi(matches[1])
		if err == nil {
			product.ID = id
		}
	}

	product.URL = url

	// Process structured data if available
	if rawData != nil {
		if offers, ok := rawData["offers"].(map[string]interface{}); ok {
			if price, err := strconv.ParseFloat(offers["price"].(string), 64); err == nil {
				product.Price = models.Price{
					SellingPrice: price,
					DiscountedPrice: price, // You might want to adjust this based on your business logic
				}
			}
		}
	}

	return &product, nil
}
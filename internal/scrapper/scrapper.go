package scrapper

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/Niiaks/price-tracker/internal/db"
	"github.com/Niiaks/price-tracker/pkg"
	"github.com/PuerkitoBio/goquery"
	"github.com/jmoiron/sqlx"
)

type Scrapper interface {
	Scrape(ctx context.Context, p pkg.Product) error
}
type ScraperService struct {
	*sqlx.DB
}

func NewScrapper(db *sqlx.DB) Scrapper {
	return &ScraperService{DB: db}
}

func (s *ScraperService) Scrape(ctx context.Context, p pkg.Product) error {

	fmt.Println("Starting to scrape product with id", p.ID)
	product, err := db.InsertProduct(ctx, s.DB, p)
	if err != nil {
		return err
	}

	// set browser headers to avoid being detected as bot
	req, _ := http.NewRequest("GET", p.Url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (Chrome/120.0.0.0 Safari/537.36)")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Connection", "keep-alive")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)
	if resp.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return err
	}

	price := doc.Find("span[data-price='true']").First().Text()
	parsedPrice, err := parseFloat(price)
	if err != nil {
		return err
	}

	priceHistory := pkg.PriceHistory{
		ProductID: product.ID,
		Price:     parsedPrice,
	}
	err = db.InsertPriceHistory(ctx, s.DB, priceHistory)
	if err != nil {
		return err
	}
	fmt.Println("Successfully scraped product", product)
	return nil
}

func parseFloat(raw string) (float64, error) {
	cleaned := strings.ReplaceAll(raw, "GH₵", "")
	cleaned = strings.ReplaceAll(cleaned, " ", "")
	cleaned = strings.ReplaceAll(cleaned, ",", "")
	cleaned = strings.TrimSpace(cleaned)

	return strconv.ParseFloat(cleaned, 64)
}

package scrapper

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Niiaks/price-tracker/pkg"
	"github.com/PuerkitoBio/goquery"
)

func Scrape(p pkg.Product) {
	// set browser headers to avoid being detected as bot
	req, _ := http.NewRequest("GET", p.Url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (Chrome/120.0.0.0 Safari/537.36)")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Connection", "keep-alive")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error getting url", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	price := doc.Find("span[data-price='true']").First().Text()
	fmt.Println("the price is", price)
}

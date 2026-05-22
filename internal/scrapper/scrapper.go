package scrapper

import (
	"fmt"
	"log"
	"net/http"
	"strings"

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

	doc.Find("a.core").Each(func(i int, s *goquery.Selection) {
		name := s.Find("h3.name").Text()
		price := s.Find("div.prc").Text()
		oldPrice := s.Find("div.old").Text()
		discount := s.Find("div.bdg._dsct").Text()
		link, _ := s.Attr("href")

		fmt.Printf("--- Product %d ---\n", i+1)
		fmt.Printf("Name:     %s\n", strings.TrimSpace(name))
		fmt.Printf("Price:    %s\n", strings.TrimSpace(price))
		fmt.Printf("Was:      %s\n", strings.TrimSpace(oldPrice))
		fmt.Printf("Discount: %s\n", strings.TrimSpace(discount))
		fmt.Printf("Link:     https://www.jumia.com.gh%s\n\n", link)
	})
}

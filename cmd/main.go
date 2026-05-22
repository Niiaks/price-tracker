package main

import (
	"fmt"

	"github.com/Niiaks/price-tracker/internal/scrapper"
	"github.com/Niiaks/price-tracker/pkg"
)

func main() {
	fmt.Println("price tracker")

	p := pkg.Product{
		Name: "iphone",
		Url:  "https://www.jumia.com.gh/mobile-phones/all-products/apple/#catalog-listing",
	}

	scrapper.Scrape(p)
}

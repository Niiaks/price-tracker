package pkg

import "time"

type Product struct {
	ID        string
	Name      string
	Threshold float64 // alerting when price drops below
	Url       string  // webpage to scrape
	CreatedAt time.Time
}

package pkg

import "time"

type Product struct {
	ID        int       `db:"id"`
	Name      string    `db:"name"`
	Threshold float64   `db:"threshold"` // alerting when price drops below
	Url       string    `db:"url"`       // webpage to scrape
	CreatedAt time.Time `db:"created_at"`
}

type PriceHistory struct {
	ID         int       `db:"id"`
	ProductID  int       `db:"product_id"`
	Price      float64   `db:"price"`
	ScrappedAt time.Time `db:"scrapped_at"`
}

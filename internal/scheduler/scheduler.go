package scheduler

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/Niiaks/price-tracker/internal/db"
	"github.com/Niiaks/price-tracker/internal/scrapper"
	"github.com/Niiaks/price-tracker/pkg"
	"github.com/jmoiron/sqlx"
)

type Scheduler struct {
	*sqlx.DB
	scrapper.Scrapper
}

func NewScheduler(db *sqlx.DB, s scrapper.Scrapper) *Scheduler {
	return &Scheduler{db, s}
}

func (s *Scheduler) Start(ctx context.Context) error {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			products, err := db.GetProducts(ctx, s.DB)
			if err != nil {
				continue
			}
			var wg sync.WaitGroup
			for _, product := range products {
				wg.Add(1)
				go func(p pkg.Product) {
					defer wg.Done()
					if err := s.Scrape(ctx, p); err != nil {
						log.Println(err)
					}
				}(product)
			}
			wg.Wait()
		}
	}
}

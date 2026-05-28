package main

import (
	"context"
	"flag"
	"log"
	"os"

	db2 "github.com/Niiaks/price-tracker/internal/db"
	"github.com/Niiaks/price-tracker/internal/scheduler"
	"github.com/Niiaks/price-tracker/internal/scrapper"
	"github.com/joho/godotenv"
)

var (
	dsn       = flag.String("dsn", "", "dsn to connect to postgresql")
	name      = flag.String("name", "", "name of product to scrape")
	threshold = flag.Float64("threshold", 0.0, "threshold of product to scrape")
	url       = flag.String("url", "", "url of product to scrape")
)

func main() {
	flag.Parse()
	if len(flag.Args()) > 0 {
		flag.Usage()
		os.Exit(1)
	}

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	if *dsn == "" {
		*dsn = os.Getenv("DATABASE_URL")
	}
	if *name == "" {
		log.Fatal("--name is required")
	}
	if *url == "" {
		log.Fatal("--url is required")
	}

	db, err := db2.ConnectDB(*dsn)
	if err != nil {
		log.Fatal(err)
	}

	s := scrapper.NewScrapper(db)
	schedule := scheduler.NewScheduler(db, s)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = schedule.Start(ctx)
	if err != nil {
		return
	}
}

package main

import (
	"flag"
	"log"
	"os"

	db2 "github.com/Niiaks/price-tracker/internal/db"
)

var dsn = flag.String("dsn", "", "dsn to connect to postgresql")

func main() {
	flag.Parse()
	if len(flag.Args()) > 0 {
		flag.Usage()
		os.Exit(1)
	}

	if *dsn == "" {
		*dsn = os.Getenv("DATABASE_URL")
	}

	_, err := db2.ConnectDB(*dsn)
	if err != nil {
		log.Fatal(err)
	}

}

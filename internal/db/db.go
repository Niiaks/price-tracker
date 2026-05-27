package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var schema = `
CREATE TABLE IF NOT EXISTS products (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    url VARCHAR(255) NOT NULL,
    threshold DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS price_history (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    product_id INTEGER REFERENCES products (id) ON DELETE CASCADE,
    price DECIMAL(10,2) NOT NULL,
    scrapped_at TIMESTAMPTZ DEFAULT NOW()
)
`

func ConnectDB(dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}
	fmt.Println("Connected to database")

	db.MustExec(schema)

	return db, nil
}

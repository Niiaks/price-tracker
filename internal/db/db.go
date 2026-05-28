package db

import (
	"context"
	"fmt"

	"github.com/Niiaks/price-tracker/pkg"
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

func InsertProduct(ctx context.Context, db *sqlx.DB, product pkg.Product) (*pkg.Product, error) {
	query := `INSERT INTO products(name, url, threshold) 
          VALUES (:name, :url, :threshold)
          RETURNING *`

	rows, err := db.NamedQueryContext(ctx, query, product)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if rows.Next() {
		err = rows.StructScan(&product)
		if err != nil {
			return nil, err
		}
	}
	return &product, nil
}

func InsertPriceHistory(ctx context.Context, db *sqlx.DB, price pkg.PriceHistory) error {
	query := `INSERT INTO price_history(product_id, price) VALUES (:product_id,:price)`
	_, err := db.NamedExecContext(ctx, query, price)
	if err != nil {
		return err
	}
	return nil
}

func GetProducts(ctx context.Context, db *sqlx.DB) ([]pkg.Product, error) {
	var product []pkg.Product
	query := `SELECT * FROM products`

	err := db.Select(&product, query)
	if err != nil {
		return nil, err
	}
	return product, nil
}

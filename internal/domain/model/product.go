package model

type Product struct {
	ID          string
	Status      string
	ProductType string `db:"product_type"`
}

package models

import "time"

type V1OrderItemDal struct {
	Id            int64     `db:"id"`
	OrderId       int64     `db:"order_id"`
	ProductId     int64     `db:"product_id"`
	Quantity      int       `db:"quantity"`
	ProductTitle  string    `db:"product_title"`
	ProductUrl    string    `db:"product_url"`
	PriceCents    int64     `db:"price_cents"`
	PriceCurrency string    `db:"price_currency"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}

func (V1OrderItemDal) PgTypeName() string { return "v1_order_item" }

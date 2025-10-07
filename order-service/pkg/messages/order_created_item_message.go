package messages

import "time"

type OrderCreatedItemMessage struct {
	Id            int64     `json:"id"`
	OrderId       int64     `json:"order_id"`
	ProductId     int64     `json:"product_id"`
	Quantity      int       `json:"quantity"`
	ProductTitle  string    `json:"product_title"`
	ProductUrl    string    `json:"product_url"`
	PriceCents    int64     `json:"price_cents"`
	PriceCurrency string    `json:"price_currency"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

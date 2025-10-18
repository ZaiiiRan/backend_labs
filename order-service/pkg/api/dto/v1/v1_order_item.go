package dto

import "time"

type V1OrderItem struct {
	ID           int64     `json:"id"`
	OrderID      int64     `json:"order_id"`
	ProductID    int64     `json:"product_id"`
	Quantity     int       `json:"quantity"`
	ProductTitle string    `json:"product_title"`
	ProductURL   string    `json:"product_url"`
	PriceCents   int64     `json:"price_cents"`
	PriceCurr    string    `json:"price_currency"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

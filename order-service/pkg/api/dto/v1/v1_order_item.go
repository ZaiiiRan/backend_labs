package dto

import "time"

type V1OrderItem struct {
	ID           int64     `json:"id"`
	OrderID      int64     `json:"orderId"`
	ProductID    int64     `json:"productId"`
	Quantity     int       `json:"quantity"`
	ProductTitle string    `json:"productTitle"`
	ProductURL   string    `json:"productUrl"`
	PriceCents   int64     `json:"priceCents"`
	PriceCurr    string    `json:"priceCurrency"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

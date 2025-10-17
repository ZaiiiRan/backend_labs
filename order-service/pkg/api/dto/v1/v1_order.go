package dto

import "time"

type V1Order struct {
	ID              int64         `json:"id"`
	CustomerID      int64         `json:"customer_id"`
	DeliveryAddress string        `json:"delivery_address"`
	TotalPriceCents int64         `json:"total_price_cents"`
	TotalPriceCurr  string        `json:"total_price_currency"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`
	OrderItems      []V1OrderItem `json:"order_items"`
}

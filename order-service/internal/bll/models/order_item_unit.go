package models

import "time"

type OrderItemUnit struct {
	ID           int64
	OrderID      int64
	ProductID    int64
	Quantity     int
	ProductTitle string
	ProductURL   string
	PriceCents   int64
	PriceCurr    string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

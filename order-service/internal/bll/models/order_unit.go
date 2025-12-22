package models

import "time"

type OrderUnit struct {
	ID              int64
	CustomerID      int64
	DeliveryAddress string
	TotalPriceCents int64
	TotalPriceCurr  string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	OrderItems      []OrderItemUnit
	Status          OrderStatus
}

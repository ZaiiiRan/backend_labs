package dto

import "time"

type V1Order struct {
	ID              int64         `json:"id"`
	CustomerID      int64         `json:"customerId"`
	DeliveryAddress string        `json:"deliveryAddress"`
	TotalPriceCents int64         `json:"totalPriceCents"`
	TotalPriceCurr  string        `json:"totalPriceCurrency"`
	CreatedAt       time.Time     `json:"createdAt"`
	UpdatedAt       time.Time     `json:"updatedAt"`
	OrderItems      []V1OrderItem `json:"orderItems"`
}

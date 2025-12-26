package messages

import "time"

type OrderCreatedMessage struct {
	Id              int64                     `json:"id"`
	CustomerID      int64                     `json:"customer_id"`
	DeliveryAddress string                    `json:"delivery_address"`
	TotalPriceCents int64                     `json:"total_price_cents"`
	TotalPriceCurr  string                    `json:"total_price_curr"`
	CreatedAt       time.Time                 `json:"created_at"`
	UpdatedAt       time.Time                 `json:"updated_at"`
	OrderItems      []OrderCreatedItemMessage `json:"order_items"`
}

func (m *OrderCreatedMessage) RoutingKey() string {
	return "order.created"
}

package dto

import "time"

type V1LogOrder struct {
	Id          int64     `json:"id"`
	OrderId     int64     `json:"orderId"`
	OrderItemId int64     `json:"orderItemId"`
	CustomerId  int64     `json:"customerId"`
	OrderStatus string    `json:"orderStatus"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

package dto

import "time"

type V1LogOrder struct {
	Id          int64     `json:"id"`
	OrderId     int64     `json:"order_id"`
	OrderItemId int64     `json:"order_item_id"`
	CustomerId  int64     `json:"customer_id"`
	OrderStatus string    `json:"order_status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

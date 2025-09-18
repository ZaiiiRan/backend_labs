package models

import "time"

type V1OrderDal struct {
	Id                 int64     `db:"id"`
	CustomerId         int64     `json:"customer_id"`
	DeliveryAdress     string    `json:"delivery_address"`
	TotalPriceCents    int64     `json:"total_price_cents"`
	TotalPriceCurrency string    `json:"total_price_currency"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

func (V1OrderDal) PgTypeName() string { return "v1_order" }

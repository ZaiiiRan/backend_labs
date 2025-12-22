package models

import "time"

type V1OrderDal struct {
	ID              int64     `db:"id"`
	CustomerID      int64     `db:"customer_id"`
	DeliveryAddress string    `db:"delivery_address"`
	TotalPriceCents int64     `db:"total_price_cents"`
	TotalPriceCurr  string    `db:"total_price_currency"`
	Status          string    `db:"status"`
	CreatedAt       time.Time `db:"created_at"`
	UpdatedAt       time.Time `db:"updated_at"`
}

func (o V1OrderDal) IsNull() bool { return false }
func (o V1OrderDal) Index(i int) any {
	switch i {
	case 0:
		return o.ID
	case 1:
		return o.CustomerID
	case 2:
		return o.DeliveryAddress
	case 3:
		return o.TotalPriceCents
	case 4:
		return o.TotalPriceCurr
	case 5:
		return o.CreatedAt
	case 6:
		return o.UpdatedAt
	case 7:
		return o.Status
	default:
		return nil
	}
}

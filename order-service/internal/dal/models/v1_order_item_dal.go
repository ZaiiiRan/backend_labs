package models

import "time"

type V1OrderItemDal struct {
	ID           int64     `db:"id"`
	OrderID      int64     `db:"order_id"`
	ProductID    int64     `db:"product_id"`
	Quantity     int       `db:"quantity"`
	ProductTitle string    `db:"product_title"`
	ProductURL   string    `db:"product_url"`
	PriceCents   int64     `db:"price_cents"`
	PriceCurr    string    `db:"price_currency"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

func (i V1OrderItemDal) IsNull() bool { return false }
func (i V1OrderItemDal) Index(idx int) any {
	switch idx {
	case 0:
		return i.ID
	case 1:
		return i.OrderID
	case 2:
		return i.ProductID
	case 3:
		return i.Quantity
	case 4:
		return i.ProductTitle
	case 5:
		return i.ProductURL
	case 6:
		return i.PriceCents
	case 7:
		return i.PriceCurr
	case 8:
		return i.CreatedAt
	case 9:
		return i.UpdatedAt
	default:
		return nil
	}
}

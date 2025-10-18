package models

import "time"

type V1AuditLogOrderDal struct {
	ID          int64     `db:"id"`
	OrderID     int64     `db:"order_id"`
	OrderItemID int64     `db:"order_item_id"`
	CustomerID  int64     `db:"customer_id"`
	OrderStatus string    `db:"order_status"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

func (l V1AuditLogOrderDal) IsNull() bool { return false }
func (l V1AuditLogOrderDal) Index(i int) any {
	switch i {
	case 0:
		return l.ID
	case 1:
		return l.OrderID
	case 2:
		return l.OrderItemID
	case 3:
		return l.CustomerID
	case 4:
		return l.OrderStatus
	case 5:
		return l.CreatedAt
	case 6:
		return l.UpdatedAt
	default:
		return nil
	}
}

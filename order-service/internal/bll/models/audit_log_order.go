package models

import "time"

type AuditLogOrder struct {
	ID          int64
	OrderID     int64
	OrderItemID int64
	CustomerID  int64
	OrderStatus OrderStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

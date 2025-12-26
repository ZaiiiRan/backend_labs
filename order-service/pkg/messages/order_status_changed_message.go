package messages

type OrderStatusChangedMessage struct {
	OrderId     int64  `json:"order_id"`
	CustomerId  int64  `json:"customer_id"`
	OrderStatus string `json:"order_status"`
}

func (m *OrderStatusChangedMessage) RoutingKey() string {
	return "order.status.changed"
}

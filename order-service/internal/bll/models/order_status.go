package models

type OrderStatus string

const (
	ORDER_STATUS_CREATED OrderStatus = "created"
)

func (s OrderStatus) String() string {
	switch s {
	case ORDER_STATUS_CREATED:
		return "created"
	default:
		return ""
	}
}

func StringToOrderStatus(str string) OrderStatus {
	if str == "created" {
		return ORDER_STATUS_CREATED
	}
	return ""
}

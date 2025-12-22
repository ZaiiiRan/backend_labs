package models

type OrderStatus string

const (
	ORDER_STATUS_CREATED    OrderStatus = "created"
	ORDER_STATUS_CANCELLED  OrderStatus = "cancelled"
	ORDER_STATUS_PROCESSING OrderStatus = "processing"
	ORDER_STATUS_COMPLETED  OrderStatus = "completed"
)

func (s OrderStatus) String() string {
	switch s {
	case ORDER_STATUS_CREATED:
		return "created"
	case ORDER_STATUS_CANCELLED:
		return "cancelled"
	case ORDER_STATUS_PROCESSING:
		return "processing"
	case ORDER_STATUS_COMPLETED:
		return "completed"
	default:
		return ""
	}
}

func (s OrderStatus) CanTransition(targetStatus OrderStatus) bool {
	if s == targetStatus {
		return true
	}

	switch s {
	case ORDER_STATUS_CREATED:
		return targetStatus == ORDER_STATUS_CANCELLED || targetStatus == ORDER_STATUS_PROCESSING
	case ORDER_STATUS_PROCESSING:
		return targetStatus == ORDER_STATUS_CANCELLED || targetStatus == ORDER_STATUS_COMPLETED
	case ORDER_STATUS_CANCELLED:
		return false
	case ORDER_STATUS_COMPLETED:
		return false
	default:
		return false
	}
}

func StringToOrderStatus(str string) OrderStatus {
	switch str {
	case "created":
		return ORDER_STATUS_CREATED
	case "cancelled":
		return ORDER_STATUS_CANCELLED
	case "processing":
		return ORDER_STATUS_PROCESSING
	case "completed":
		return ORDER_STATUS_COMPLETED
	default:
		return ""
	}
}

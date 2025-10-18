package mappers

import (
	bll "github.com/ZaiiiRan/backend_labs/order-service/internal/bll/models"
	dal "github.com/ZaiiiRan/backend_labs/order-service/internal/dal/models"
	"github.com/ZaiiiRan/backend_labs/order-service/pkg/api/dto/v1"
)

func DalAuditLogOrderToBll(i dal.V1AuditLogOrderDal) bll.AuditLogOrder {
	return bll.AuditLogOrder{
		ID:          i.ID,
		OrderID:     i.OrderID,
		OrderItemID: i.OrderItemID,
		CustomerID:  i.CustomerID,
		OrderStatus: bll.StringToOrderStatus(i.OrderStatus),
		CreatedAt:   i.CreatedAt,
		UpdatedAt:   i.UpdatedAt,
	}
}

func BllAuditLogOrderToDal(i bll.AuditLogOrder) dal.V1AuditLogOrderDal {
	return dal.V1AuditLogOrderDal{
		ID:          i.ID,
		OrderID:     i.OrderID,
		OrderItemID: i.OrderItemID,
		CustomerID:  i.CustomerID,
		OrderStatus: i.OrderStatus.String(),
		CreatedAt:   i.CreatedAt,
		UpdatedAt:   i.UpdatedAt,
	}
}

func DtoAuditLogOrderToBll(i dto.V1LogOrder) bll.AuditLogOrder {
	return bll.AuditLogOrder{
		ID:          i.Id,
		OrderID:     i.OrderId,
		OrderItemID: i.OrderItemId,
		CustomerID:  i.CustomerId,
		OrderStatus: bll.StringToOrderStatus(i.OrderStatus),
		CreatedAt:   i.CreatedAt,
		UpdatedAt:   i.UpdatedAt,
	}
}

func BllAuditLogOrderToDto(i bll.AuditLogOrder) dto.V1LogOrder {
	return dto.V1LogOrder{
		Id:          i.ID,
		OrderId:     i.OrderID,
		OrderItemId: i.OrderItemID,
		OrderStatus: i.OrderStatus.String(),
		CustomerId:  i.CustomerID,
		CreatedAt:   i.CreatedAt,
		UpdatedAt:   i.UpdatedAt,
	}
}

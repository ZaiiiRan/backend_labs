package mappers

import (
	bll "github.com/ZaiiiRan/backend_labs/order-service/internal/bll/models"
	dal "github.com/ZaiiiRan/backend_labs/order-service/internal/dal/models"
	"github.com/ZaiiiRan/backend_labs/order-service/pkg/api/dto"
)

func DalOrderToBll(o dal.V1OrderDal, items []bll.OrderItemUnit) bll.OrderUnit {
	return bll.OrderUnit{
		ID:              o.ID,
		CustomerID:      o.CustomerID,
		DeliveryAddress: o.DeliveryAddress,
		TotalPriceCents: o.TotalPriceCents,
		TotalPriceCurr:  o.TotalPriceCurr,
		CreatedAt:       o.CreatedAt,
		UpdatedAt:       o.UpdatedAt,
		OrderItems:      items,
	}
}

func BllOrderToDal(o bll.OrderUnit) dal.V1OrderDal {
	return dal.V1OrderDal{
		ID:              o.ID,
		CustomerID:      o.CustomerID,
		DeliveryAddress: o.DeliveryAddress,
		TotalPriceCents: o.TotalPriceCents,
		TotalPriceCurr:  o.TotalPriceCurr,
		CreatedAt:       o.CreatedAt.UTC(),
		UpdatedAt:       o.UpdatedAt.UTC(),
	}
}

func DtoOrderToBll(o dto.V1Order) bll.OrderUnit {
	var items []bll.OrderItemUnit
	for _, it := range o.OrderItems {
		items = append(items, DtoOrderItemToBll(it))
	}

	return bll.OrderUnit{
		ID:              o.ID,
		CustomerID:      o.CustomerID,
		DeliveryAddress: o.DeliveryAddress,
		TotalPriceCents: o.TotalPriceCents,
		TotalPriceCurr:  o.TotalPriceCurr,
		CreatedAt:       o.CreatedAt,
		UpdatedAt:       o.UpdatedAt,
		OrderItems:      items,
	}
}

func BllOrderToDto(o bll.OrderUnit) dto.V1Order {
	var items []dto.V1OrderItem
	for _, it := range o.OrderItems {
		items = append(items, BllOrderItemToDto(it))
	}

	return dto.V1Order{
		ID:              o.ID,
		CustomerID:      o.CustomerID,
		DeliveryAddress: o.DeliveryAddress,
		TotalPriceCents: o.TotalPriceCents,
		TotalPriceCurr:  o.TotalPriceCurr,
		CreatedAt:       o.CreatedAt,
		UpdatedAt:       o.UpdatedAt,
		OrderItems:      items,
	}
}

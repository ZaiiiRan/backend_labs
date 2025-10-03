package mappers

import (
	bll "github.com/ZaiiiRan/backend_labs/order-service/internal/bll/models"
	dal "github.com/ZaiiiRan/backend_labs/order-service/internal/dal/models"
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


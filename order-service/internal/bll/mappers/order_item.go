package mappers

import (
	bll "github.com/ZaiiiRan/backend_labs/order-service/internal/bll/models"
	dal "github.com/ZaiiiRan/backend_labs/order-service/internal/dal/models"
)

func DalOrderItemToBll(i dal.V1OrderItemDal) bll.OrderItemUnit {
	return bll.OrderItemUnit{
		ID:           i.ID,
		OrderID:      i.OrderID,
		ProductID:    i.ProductID,
		Quantity:     i.Quantity,
		ProductTitle: i.ProductTitle,
		ProductURL:   i.ProductURL,
		PriceCents:   i.PriceCents,
		PriceCurr:    i.PriceCurr,
		CreatedAt:    i.CreatedAt,
		UpdatedAt:    i.UpdatedAt,
	}
}

func BllOrderItemToDal(i bll.OrderItemUnit, orderID int64) dal.V1OrderItemDal {
	return dal.V1OrderItemDal{
		ID:           i.ID,
		OrderID:      orderID,
		ProductID:    i.ProductID,
		Quantity:     i.Quantity,
		ProductTitle: i.ProductTitle,
		ProductURL:   i.ProductURL,
		PriceCents:   i.PriceCents,
		PriceCurr:    i.PriceCurr,
		CreatedAt:    i.CreatedAt,
		UpdatedAt:    i.UpdatedAt,
	}
}

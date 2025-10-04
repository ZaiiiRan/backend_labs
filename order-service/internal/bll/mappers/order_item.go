package mappers

import (
	bll "github.com/ZaiiiRan/backend_labs/order-service/internal/bll/models"
	dal "github.com/ZaiiiRan/backend_labs/order-service/internal/dal/models"
	"github.com/ZaiiiRan/backend_labs/order-service/pkg/api/dto"
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
		CreatedAt:    i.CreatedAt.UTC(),
		UpdatedAt:    i.UpdatedAt.UTC(),
	}
}

func DtoOrderItemToBll(it dto.V1OrderItem) bll.OrderItemUnit {
	return bll.OrderItemUnit{
		ID:           it.ID,
		OrderID:      it.OrderID,
		ProductID:    it.ProductID,
		Quantity:     it.Quantity,
		ProductTitle: it.ProductTitle,
		ProductURL:   it.ProductURL,
		PriceCents:   it.PriceCents,
		PriceCurr:    it.PriceCurr,
		CreatedAt:    it.CreatedAt,
		UpdatedAt:    it.UpdatedAt,
	}
}

func BllOrderItemToDto(it bll.OrderItemUnit) dto.V1OrderItem {
	return dto.V1OrderItem{
		ID:           it.ID,
		OrderID:      it.OrderID,
		ProductID:    it.ProductID,
		Quantity:     it.Quantity,
		ProductTitle: it.ProductTitle,
		ProductURL:   it.ProductURL,
		PriceCents:   it.PriceCents,
		PriceCurr:    it.PriceCurr,
		CreatedAt:    it.CreatedAt,
		UpdatedAt:    it.UpdatedAt,
	}
}

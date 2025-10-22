package mappers

import (
	pb "github.com/ZaiiiRan/backend_labs/order-service/gen/go/order-service/v1"
	bll "github.com/ZaiiiRan/backend_labs/order-service/internal/bll/models"
	dal "github.com/ZaiiiRan/backend_labs/order-service/internal/dal/models"
	"github.com/ZaiiiRan/backend_labs/order-service/pkg/messages"
	"google.golang.org/protobuf/types/known/timestamppb"
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

func PbOrderItemToBll(it *pb.OrderItem) bll.OrderItemUnit {
	return bll.OrderItemUnit{
		ID:           it.Id,
		OrderID:      it.OrderId,
		ProductID:    it.ProductId,
		Quantity:     int(it.Quantity),
		ProductTitle: it.ProductTitle,
		ProductURL:   it.ProductUrl,
		PriceCents:   it.PriceCents,
		PriceCurr:    it.PriceCurrency,
		CreatedAt:    it.CreatedAt.AsTime(),
		UpdatedAt:    it.UpdatedAt.AsTime(),
	}
}

func BllOrderItemToPb(it bll.OrderItemUnit) *pb.OrderItem {
	createdAt := timestamppb.New(it.CreatedAt)
	updatedAt := timestamppb.New(it.UpdatedAt)

	return &pb.OrderItem{
		Id:            it.ID,
		OrderId:       it.OrderID,
		ProductId:     it.ProductID,
		Quantity:      int32(it.Quantity),
		ProductTitle:  it.ProductTitle,
		ProductUrl:    it.ProductURL,
		PriceCents:    it.PriceCents,
		PriceCurrency: it.PriceCurr,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}
}

func BllOrderItemToOrderCreatedItemMessage(it bll.OrderItemUnit) messages.OrderCreatedItemMessage {
	return messages.OrderCreatedItemMessage{
		Id:            it.ID,
		OrderId:       it.OrderID,
		ProductId:     it.ProductID,
		Quantity:      it.Quantity,
		ProductTitle:  it.ProductTitle,
		ProductUrl:    it.ProductURL,
		PriceCents:    it.PriceCents,
		PriceCurrency: it.PriceCurr,
		CreatedAt:     it.CreatedAt,
		UpdatedAt:     it.UpdatedAt,
	}
}

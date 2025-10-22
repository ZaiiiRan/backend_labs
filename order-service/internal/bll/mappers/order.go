package mappers

import (
	pb "github.com/ZaiiiRan/backend_labs/order-service/gen/go/order-service/v1"
	bll "github.com/ZaiiiRan/backend_labs/order-service/internal/bll/models"
	dal "github.com/ZaiiiRan/backend_labs/order-service/internal/dal/models"
	"github.com/ZaiiiRan/backend_labs/order-service/pkg/api/dto/v1"
	"github.com/ZaiiiRan/backend_labs/order-service/pkg/messages"
	"google.golang.org/protobuf/types/known/timestamppb"
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

func PbOrderToBll(o *pb.Order) bll.OrderUnit {
	var items []bll.OrderItemUnit
	for _, it := range o.OrderItems {
		items = append(items, PbOrderItemToBll(it))
	}

	return bll.OrderUnit{
		ID:              o.Id,
		CustomerID:      o.CustomerId,
		DeliveryAddress: o.DeliveryAddress,
		TotalPriceCents: o.TotalPriceCents,
		TotalPriceCurr:  o.TotalPriceCurrency,
		CreatedAt:       o.CreatedAt.AsTime(),
		UpdatedAt:       o.UpdatedAt.AsTime(),
		OrderItems:      items,
	}
}

func BllOrderToPb(o bll.OrderUnit) *pb.Order {
	createdAt := timestamppb.New(o.CreatedAt)
	updatedAt := timestamppb.New(o.UpdatedAt)

	var items []*pb.OrderItem
	for _, it := range o.OrderItems {
		items = append(items, BllOrderItemToPb(it))
	}

	return &pb.Order{
		Id:                 o.ID,
		CustomerId:         o.CustomerID,
		DeliveryAddress:    o.DeliveryAddress,
		TotalPriceCents:    o.TotalPriceCents,
		TotalPriceCurrency: o.TotalPriceCurr,
		CreatedAt:          createdAt,
		UpdatedAt:          updatedAt,
		OrderItems:         items,
	}
}

func BllOrderToOrderCreatedMessage(o bll.OrderUnit) messages.OrderCreatedMessage {
	var items []messages.OrderCreatedItemMessage
	for _, it := range o.OrderItems {
		items = append(items, BllOrderItemToOrderCreatedItemMessage(it))
	}

	return messages.OrderCreatedMessage{
		Id:              o.ID,
		CustomerID:      o.CustomerID,
		DeliveryAddress: o.DeliveryAddress,
		TotalPriceCents: o.TotalPriceCents,
		TotalPriceCurr:  o.TotalPriceCurr,
		CreatedAt:       o.CreatedAt,
		UpdatedAt:       o.UpdatedAt,
		OrderItems:      items,
	}
}

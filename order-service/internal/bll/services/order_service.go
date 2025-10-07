package services

import (
	"context"
	"time"

	"github.com/ZaiiiRan/backend_labs/order-service/internal/bll/mappers"
	bll "github.com/ZaiiiRan/backend_labs/order-service/internal/bll/models"
	"github.com/ZaiiiRan/backend_labs/order-service/internal/dal/interfaces"
	dal "github.com/ZaiiiRan/backend_labs/order-service/internal/dal/models"
	"github.com/ZaiiiRan/backend_labs/order-service/internal/dal/rabbitmq"
	unitofwork "github.com/ZaiiiRan/backend_labs/order-service/internal/dal/unit_of_work/postgres"
)

type OrderService struct {
	uow           *unitofwork.UnitOfWork
	orderRepo     interfaces.OrderRepository
	orderItemRepo interfaces.OrderItemRepository
	publisher     *rabbitmq.Publisher
}

func NewOrderService(uow *unitofwork.UnitOfWork, orderRepo interfaces.OrderRepository, orderItemRepo interfaces.OrderItemRepository, publisher *rabbitmq.Publisher) *OrderService {
	return &OrderService{
		uow:           uow,
		orderRepo:     orderRepo,
		orderItemRepo: orderItemRepo,
		publisher:     publisher,
	}
}

func (s *OrderService) BatchInsert(ctx context.Context, orders []bll.OrderUnit) ([]bll.OrderUnit, error) {
	now := time.Now().UTC()

	_, err := s.uow.BeginTransaction(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			s.uow.Rollback(ctx)
		}
	}()

	var dalOrders []dal.V1OrderDal
	for _, o := range orders {
		d := mappers.BllOrderToDal(o)
		d.CreatedAt = now
		d.UpdatedAt = now
		dalOrders = append(dalOrders, d)
	}

	insertedOrders, err := s.orderRepo.BulkInsert(ctx, dalOrders)
	if err != nil {
		return nil, err
	}

	var dalItems []dal.V1OrderItemDal
	for idx, insOrder := range insertedOrders {
		for _, item := range orders[idx].OrderItems {
			d := mappers.BllOrderItemToDal(item, insOrder.ID)
			d.CreatedAt = now
			d.UpdatedAt = now
			dalItems = append(dalItems, d)
		}
	}

	insertedItems, err := s.orderItemRepo.BulkInsert(ctx, dalItems)
	if err != nil {
		return nil, err
	}

	if err := s.uow.Commit(ctx); err != nil {
		return nil, err
	}

	itemLookup := make(map[int64][]bll.OrderItemUnit)
	for _, it := range insertedItems {
		itemLookup[it.OrderID] = append(itemLookup[it.OrderID], mappers.DalOrderItemToBll(it))
	}

	var result []bll.OrderUnit
	for _, o := range insertedOrders {
		result = append(result, mappers.DalOrderToBll(o, itemLookup[o.ID]))
	}

	var msgs []any
	for _, o := range result {
		msgs = append(msgs, mappers.BllOrderToOrderCreatedMessage(o))
	}

	if err := s.publisher.PublishBatch(ctx, s.publisher.QueueOrderCreated(), msgs); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *OrderService) GetOrders(ctx context.Context, query bll.QueryOrderItemsModel) ([]bll.OrderUnit, error) {
	orders, err := s.orderRepo.Query(ctx, dal.QueryOrdersDalModel{
		IDs:         query.IDs,
		CustomerIDs: query.CustomerIDs,
		Limit:       query.PageSize,
		Offset:      query.PageSize * (query.Page - 1),
	})
	if err != nil {
		return nil, err
	}
	if len(orders) == 0 {
		return []bll.OrderUnit{}, nil
	}

	var items []dal.V1OrderItemDal
	if query.IncludeOrderItems {
		ordersIDs := make([]int64, len(orders))
		for i, o := range orders {
			ordersIDs[i] = o.ID
		}
		items, err = s.orderItemRepo.Query(ctx, dal.QueryOrderItemsDalModel{OrderIDs: ordersIDs})
		if err != nil {
			return nil, err
		}
	}

	itemLookup := make(map[int64][]bll.OrderItemUnit)
	for _, it := range items {
		itemLookup[it.OrderID] = append(itemLookup[it.OrderID], mappers.DalOrderItemToBll(it))
	}

	var result []bll.OrderUnit
	for _, o := range orders {
		result = append(result, mappers.DalOrderToBll(o, itemLookup[o.ID]))
	}

	return result, nil
}

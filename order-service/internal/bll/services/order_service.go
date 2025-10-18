package services

import (
	"context"
	"time"

	"github.com/ZaiiiRan/backend_labs/order-service/internal/bll/mappers"
	bll "github.com/ZaiiiRan/backend_labs/order-service/internal/bll/models"
	"github.com/ZaiiiRan/backend_labs/order-service/internal/dal/interfaces"
	dal "github.com/ZaiiiRan/backend_labs/order-service/internal/dal/models"
	publisher "github.com/ZaiiiRan/backend_labs/order-service/internal/dal/publisher/rabbitmq"
	unitofwork "github.com/ZaiiiRan/backend_labs/order-service/internal/dal/unit_of_work/postgres"
	"go.uber.org/zap"
)

type OrderService struct {
	uow                   *unitofwork.UnitOfWork
	orderRepo             interfaces.OrderRepository
	orderItemRepo         interfaces.OrderItemRepository
	orderCreatedPublisher *publisher.Publisher
	log                   *zap.SugaredLogger
}

func NewOrderService(
	uow *unitofwork.UnitOfWork,
	orderRepo interfaces.OrderRepository,
	orderItemRepo interfaces.OrderItemRepository,
	orderCreatedPublisher *publisher.Publisher,
	log *zap.SugaredLogger,
) *OrderService {
	return &OrderService{
		uow:                   uow,
		orderRepo:             orderRepo,
		orderItemRepo:         orderItemRepo,
		orderCreatedPublisher: orderCreatedPublisher,
		log:                   log,
	}
}

func (s *OrderService) BatchInsert(ctx context.Context, orders []bll.OrderUnit) ([]bll.OrderUnit, error) {
	now := time.Now().UTC()
	s.log.Infow("order_service.batch_insert_start", "orders_count", len(orders))

	_, err := s.uow.BeginTransaction(ctx)
	if err != nil {
		s.log.Errorw("order_service.begin_transaction_failed", "err", err)
		return nil, err
	}
	defer func() {
		if err != nil {
			s.uow.Rollback(ctx)
			s.log.Warnw("order_service.transaction_rollback", "err", err)
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
		s.log.Errorw("order_service.bulk_insert_orders_failed", "err", err)
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
		s.log.Errorw("order_service.bulk_insert_order_items_failed", "err", err)
		return nil, err
	}

	if err := s.uow.Commit(ctx); err != nil {
		s.log.Errorw("order_service.commit_transaction_failed", "err", err)
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

	go func() {
		var msgs []any
		for _, o := range result {
			msgs = append(msgs, mappers.BllOrderToOrderCreatedMessage(o))
		}

		ctxPub, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := s.orderCreatedPublisher.PublishBatch(ctxPub, msgs); err != nil {
			s.log.Errorw("order_service.publish_order_created_messages_failed", "err", err, "lost_msgs", msgs)
		}
	}()

	s.log.Infow("order_service.batch_insert_success", "inserted_orders_count", len(result))
	return result, nil
}

func (s *OrderService) GetOrders(ctx context.Context, query bll.QueryOrderItemsModel) ([]bll.OrderUnit, error) {
	s.log.Infow("order_service.get_orders_start", "query", query)

	orders, err := s.orderRepo.Query(ctx, dal.QueryOrdersDalModel{
		IDs:         query.IDs,
		CustomerIDs: query.CustomerIDs,
		Limit:       query.PageSize,
		Offset:      query.PageSize * (query.Page - 1),
	})
	if err != nil {
		s.log.Errorw("order_service.query_orders_failed", "err", err)
		return nil, err
	}
	if len(orders) == 0 {
		s.log.Infow("order_service.get_orders_success", "returned_orders_count", 0)
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

	s.log.Infow("order_service.get_orders_success", "returned_orders_count", len(result))
	return result, nil
}

func (s *OrderService) UnitOfWork() *unitofwork.UnitOfWork {
	return s.uow
}

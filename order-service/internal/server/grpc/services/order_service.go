package services

import (
	"context"

	pb "github.com/ZaiiiRan/backend_labs/order-service/gen/go/order-service/v1"
	"github.com/ZaiiiRan/backend_labs/order-service/internal/bll/mappers"
	"github.com/ZaiiiRan/backend_labs/order-service/internal/bll/models"
	bllServices "github.com/ZaiiiRan/backend_labs/order-service/internal/bll/services"
	"github.com/ZaiiiRan/backend_labs/order-service/internal/dal/postgres"
	producer "github.com/ZaiiiRan/backend_labs/order-service/internal/dal/producer/kafka"
	repositories "github.com/ZaiiiRan/backend_labs/order-service/internal/dal/repositories/postgres"
	unitofwork "github.com/ZaiiiRan/backend_labs/order-service/internal/dal/unit_of_work/postgres"
	"github.com/ZaiiiRan/backend_labs/order-service/internal/utils"
	"github.com/ZaiiiRan/backend_labs/order-service/internal/validators"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrderService struct {
	pb.UnimplementedOrderServiceServer

	log                        *zap.SugaredLogger
	pgClient                   *postgres.PostgresClient
	orderCreatedProducer       *producer.Producer
	orderStatusChangedProducer *producer.Producer
}

func NewOrderService(
	pgClient *postgres.PostgresClient,
	orderCreatedProducer *producer.Producer,
	orderStatusChangedProducer *producer.Producer,
	log *zap.SugaredLogger,
) *OrderService {
	return &OrderService{
		pgClient:                   pgClient,
		orderCreatedProducer:       orderCreatedProducer,
		orderStatusChangedProducer: orderStatusChangedProducer,
		log:                        log,
	}
}

func (s *OrderService) BatchCreate(ctx context.Context, req *pb.BatchCreateRequest) (*pb.BatchCreateResponse, error) {
	l := s.log.With("op", "batch_create")
	l.Infow("order_controller.batch_create_start")

	if errs := validators.ValidateBatchCreateRequest(req); errs != nil {
		l.Errorw("order_controller.batch_create_request_validation_failed", "err", errs)
		return nil, errs.ToStatus()
	}

	var orders []models.OrderUnit
	for _, o := range req.Orders {
		order := mappers.PbOrderToBll(o)
		order.Status = models.ORDER_STATUS_CREATED
		orders = append(orders, order)
	}

	orderSvc := s.createBllOrderService(l)
	defer orderSvc.UnitOfWork().Close()

	result, err := orderSvc.BatchInsert(ctx, orders)
	if err != nil {
		l.Errorw("order_controller.batch_insert_failed", "err", "Internal server error")
		return nil, status.Errorf(codes.Internal, "Internal server error")
	}

	l.Infow("order_controller.batch_create_success")

	var resp pb.BatchCreateResponse
	for _, o := range result {
		resp.Orders = append(resp.Orders, mappers.BllOrderToPb(o))
	}

	return &resp, nil
}

func (s *OrderService) QueryOrders(ctx context.Context, req *pb.QueryOrdersRequest) (*pb.QueryOrdersResponse, error) {
	l := s.log.With("op", "query_orders")
	l.Infow("order_controller.query_orders_start")

	if errs := validators.ValidateQueryOrdersRequest(req); errs != nil {
		l.Errorw("order_controller.query_orders_request_validation_failed", "err", errs)
		return nil, errs.ToStatus()
	}

	orderSvc := s.createBllOrderService(l)
	defer orderSvc.UnitOfWork().Close()

	result, err := orderSvc.GetOrders(ctx, mappers.PbQueryOrderItemsToBll(req))
	if err != nil {
		l.Errorw("order_controller.get_orders_failed", "err", "Internal server error")
		return nil, status.Errorf(codes.Internal, "Internal server error")
	}

	l.Infow("order_controller.query_orders_success")

	var resp pb.QueryOrdersResponse
	for _, o := range result {
		resp.Orders = append(resp.Orders, mappers.BllOrderToPb(o))
	}

	return &resp, nil
}

func (s *OrderService) UpdateOrdersStatus(ctx context.Context, req *pb.UpdateOrdersStatusRequest) (*pb.UpdateOrdersStatusResponse, error) {
	l := s.log.With("op", "update_orders_status")
	l.Infow("order_controller.update_orders_status_start")

	if errs := validators.ValidateUpdateOrdersStatusRequest(req); errs != nil {
		l.Errorw("order_controller.update_orders_status_request_validation_failed", "err", errs)
		return nil, errs.ToStatus()
	}

	orderSvc := s.createBllOrderService(l)
	defer orderSvc.UnitOfWork().Close()

	_, err := orderSvc.UpdateOrdersStatus(ctx, req.OrderIds, models.OrderStatus(req.NewStatus))
	if err != nil {
		l.Errorw("order_controller.update_orders_status_failed", "err", err)
		if utils.IsGrpcError(err) {
			return nil, err
		}
		return nil, status.Errorf(codes.Internal, "Internal server error")
	}

	l.Infow("order_controller.update_orders_status_success")
	return &pb.UpdateOrdersStatusResponse{}, nil
}

func (s *OrderService) AuditLogOrderBatchCreate(ctx context.Context, req *pb.AuditLogOrderBatchCreateRequest) (*pb.AuditLogOrderBatchCreateResponse, error) {
	l := s.log.With("op", "audit_log_order_batch_create")
	l.Infow("order_controller.audit_log_order_batch_create_start")

	if errs := validators.ValidateAuditLogOrderBatchCreateRequest(req); errs != nil {
		l.Errorw("order_controller.audit_log_order_batch_create_request_validation_failed", "err", errs)
		return nil, errs.ToStatus()
	}

	var logs []models.AuditLogOrder
	for _, i := range req.Orders {
		logs = append(logs, mappers.PbAuditLogOrderToBll(i))
	}

	auditLogOrderReporderSvc := s.createBllAuditLogOrderService(l)
	defer auditLogOrderReporderSvc.UnitOfWork().Close()

	result, err := auditLogOrderReporderSvc.BatchInsert(ctx, logs)
	if err != nil {
		l.Errorw("order_controller.audit_log_order_batch_insert_failed", "err", "Internal server error")
		return nil, status.Errorf(codes.Internal, "Internal server error")
	}

	l.Infow("order_controller.audit_log_order_batch_create_success")

	var resp pb.AuditLogOrderBatchCreateResponse
	for _, i := range result {
		resp.Orders = append(resp.Orders, mappers.BllAuditLogOrderToPb(i))
	}

	return &resp, nil
}

func (s *OrderService) createBllOrderService(log *zap.SugaredLogger) *bllServices.OrderService {
	uow := unitofwork.New(s.pgClient)
	orderRepo := repositories.NewOrderRepository(uow)
	orderItemRepo := repositories.NewOrderItemRepository(uow)
	return bllServices.NewOrderService(uow, orderRepo, orderItemRepo, s.orderCreatedProducer, s.orderStatusChangedProducer, log)
}

func (s *OrderService) createBllAuditLogOrderService(log *zap.SugaredLogger) *bllServices.AuditLogOrderService {
	uow := unitofwork.New(s.pgClient)
	repo := repositories.NewAuditLogOrderRepository(uow)
	return bllServices.NewAuditLogOrderService(uow, repo, log)
}

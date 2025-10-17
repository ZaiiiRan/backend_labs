package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/ZaiiiRan/backend_labs/order-service/internal/bll/mappers"
	"github.com/ZaiiiRan/backend_labs/order-service/internal/bll/models"
	"github.com/ZaiiiRan/backend_labs/order-service/internal/bll/services"
	"github.com/ZaiiiRan/backend_labs/order-service/internal/dal/postgres"
	publisher "github.com/ZaiiiRan/backend_labs/order-service/internal/dal/publisher/rabbitmq"
	repositories "github.com/ZaiiiRan/backend_labs/order-service/internal/dal/repositories/postgres"
	unitofwork "github.com/ZaiiiRan/backend_labs/order-service/internal/dal/unit_of_work/postgres"
	"github.com/ZaiiiRan/backend_labs/order-service/internal/validators"
	"github.com/ZaiiiRan/backend_labs/order-service/pkg/api/dto/v1"
	"go.uber.org/zap"
)

type OrderController struct {
	log                   *zap.SugaredLogger
	pgClient              *postgres.PostgresClient
	orderCreatedPublisher *publisher.Publisher
}

func NewOrderController(pgClient *postgres.PostgresClient, orderCreatedPublisher *publisher.Publisher, log *zap.SugaredLogger) *OrderController {
	return &OrderController{pgClient: pgClient, orderCreatedPublisher: orderCreatedPublisher, log: log}
}

// BatchCreate godoc
// @Summary Create orders batch
// @Description Creates orders with order items
// @Tags orders
// @Accept json
// @Produce json
// @Param request body dto.V1CreateOrderRequest true "Create Order Request"
// @Success 200 {object} dto.V1CreateOrderResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/order/batch-create [post]
func (c *OrderController) BatchCreate(w http.ResponseWriter, r *http.Request) {
	l := c.log.With("op", "batch_create")
	l.Infow("order_controller.batch_create_start")

	if r.Method != http.MethodPost {
		l.Warnw("order_controller.batch_create_failed", "err", "Method not allowed", "method", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req dto.V1CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		l.Errorw("order_controller.batch_create_failed", "err", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if errs := validators.ValidateV1CreateOrderRequest(&req); errs != nil {
		l.Errorw("order_controller.create_order_request_validation_failed", "err", errs)
		c.writeJSON(w, http.StatusBadRequest, errs)
		return
	}

	var orders []models.OrderUnit
	for _, o := range req.Orders {
		orders = append(orders, mappers.DtoOrderToBll(o))
	}

	orderService := c.createOrderService(l)
	defer orderService.UnitOfWork().Close()

	result, err := orderService.BatchInsert(r.Context(), orders)
	if err != nil {
		l.Errorw("order_controller.batch_insert_failed", "err", "Internal server error")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	l.Infow("order_controller.batch_create_success")

	var resp dto.V1CreateOrderResponse
	for _, o := range result {
		resp.Orders = append(resp.Orders, mappers.BllOrderToDto(o))
	}

	c.writeJSON(w, http.StatusOK, resp)
}

// QueryOrders godoc
// @Summary Query orders
// @Description Returns orders with optional order items
// @Tags orders
// @Accept json
// @Produce json
// @Param request body dto.V1QueryOrdersRequest true "Query Orders Request"
// @Success 200 {object} dto.V1QueryOrdersResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/order/query [post]
func (c *OrderController) QueryOrders(w http.ResponseWriter, r *http.Request) {
	l := c.log.With("op", "query_orders")
	l.Infow("order_controller.query_orders_start")

	if r.Method != http.MethodPost {
		l.Warnw("order_controller.query_orders_failed", "err", "Method not allowed", "method", r.Method)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req dto.V1QueryOrdersRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		l.Errorw("order_controller.query_orders_failed", "err", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if errs := validators.ValidateV1QueryOrdersRequest(req); errs != nil {
		l.Errorw("order_controller.quey_orders_request_validation_failed", "err", errs)
		c.writeJSON(w, http.StatusBadRequest, errs)
		return
	}

	orderService := c.createOrderService(l)
	defer orderService.UnitOfWork().Close()

	result, err := orderService.GetOrders(r.Context(), mappers.DtoQueryOrderItemsToBll(req))
	if err != nil {
		l.Errorw("order_controller.get_orders_failed", "err", "Internal server error")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	l.Infow("order_controller.query_orders_success")

	var resp dto.V1QueryOrdersResponse
	for _, o := range result {
		resp.Orders = append(resp.Orders, mappers.BllOrderToDto(o))
	}

	c.writeJSON(w, http.StatusOK, resp)
}

// AuditLogOrderBatchCreate godoc
// @Summary Create audit logs for orders batch
// @Description Creates audit logs for orders
// @Tags audit_logs
// @Accept json
// @Produce json
// @Param request body dto.V1CreateAuditLogOrderRequest true "Create Audit Log Order Request"
// @Success 200 {object} dto.V1CreateAuditLogOrderResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/audit-log/order/batch-create [post]
func (c *OrderController) AuditLogOrderBatchCreate(w http.ResponseWriter, r *http.Request) {
	l := c.log.With("op", "audit_log_order_batch_create")
	l.Infow("order_controller.audit_log_order_batch_create_start")

	if r.Method != http.MethodPost {
		l.Warnw("order_controller.audit_log_order_batch_create_failed", "err", "Method not allowed", "method", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req dto.V1CreateAuditLogOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		l.Errorw("order_controller.audit_log_order_batch_create_failed", "err", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if errs := validators.ValidateV1CreateAuditLogOrderRequest(req); errs != nil {
		l.Errorw("order_controller.create_audit_log_order_request_validation_failed", "err", errs)
		c.writeJSON(w, http.StatusBadRequest, errs)
		return
	}

	var logs []models.AuditLogOrder
	for _, i := range req.Orders {
		logs = append(logs, mappers.DtoAuditLogOrderToBll(i))
	}

	auditLogOrderReporderService := c.createAuditLogOrderService(l)
	defer auditLogOrderReporderService.UnitOfWork().Close()

	result, err := auditLogOrderReporderService.BatchInsert(r.Context(), logs)
	if err != nil {
		l.Errorw("order_controller.audit_log_order_batch_insert_failed", "err", "Internal server error")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	l.Infow("order_controller.audit_log_order_batch_create_success")

	var resp dto.V1CreateAuditLogOrderResponse
	for _, i := range result {
		resp.Orders = append(resp.Orders, mappers.BllAuditLogOrderToDto(i))
	}

	c.writeJSON(w, http.StatusOK, resp)
}

func (c *OrderController) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func (c *OrderController) createOrderService(log *zap.SugaredLogger) *services.OrderService {
	uow := unitofwork.New(c.pgClient)
	orderRepo := repositories.NewOrderRepository(uow)
	orderItemRepo := repositories.NewOrderItemRepository(uow)

	orderService := services.NewOrderService(uow, orderRepo, orderItemRepo, c.orderCreatedPublisher, log)
	return orderService
}

func (c *OrderController) createAuditLogOrderService(log *zap.SugaredLogger) *services.AuditLogOrderService {
	uow := unitofwork.New(c.pgClient)
	auditLogOrderRepo := repositories.NewAuditLogOrderRepository(uow)

	auditLogOrderService := services.NewAuditLogOrderService(uow, auditLogOrderRepo, log)
	return auditLogOrderService
}

package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/ZaiiiRan/backend_labs/order-service/internal/bll/models"
	"github.com/ZaiiiRan/backend_labs/order-service/internal/bll/services"
	"github.com/ZaiiiRan/backend_labs/order-service/pkg/api/dto"
)

type OrderController struct {
	orderService *services.OrderService
}

func NewOrderController(orderService *services.OrderService) *OrderController {
	return &OrderController{orderService: orderService}
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
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req dto.V1CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var orders []models.OrderUnit
	for _, o := range req.Orders {
		var items []models.OrderItemUnit
		for _, it := range o.OrderItems {
			items = append(items, models.OrderItemUnit{
				ProductID:    it.ProductID,
				Quantity:     it.Quantity,
				ProductTitle: it.ProductTitle,
				ProductURL:   it.ProductURL,
				PriceCents:   it.PriceCents,
				PriceCurr:    it.PriceCurr,
			})
		}
		orders = append(orders, models.OrderUnit{
			CustomerID:      o.CustomerID,
			DeliveryAddress: o.DeliveryAddress,
			TotalPriceCents: o.TotalPriceCents,
			TotalPriceCurr:  o.TotalPriceCurr,
			OrderItems:      items,
		})
	}

	result, err := c.orderService.BatchInsert(r.Context(), orders)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var resp dto.V1CreateOrderResponse
	for _, o := range result {
		var items []dto.V1OrderItem
		for _, it := range o.OrderItems {
			items = append(items, dto.V1OrderItem{
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
			})
		}
		resp.Orders = append(resp.Orders, dto.V1Order{
			ID:              o.ID,
			CustomerID:      o.CustomerID,
			DeliveryAddress: o.DeliveryAddress,
			TotalPriceCents: o.TotalPriceCents,
			TotalPriceCurr:  o.TotalPriceCurr,
			CreatedAt:       o.CreatedAt,
			UpdatedAt:       o.UpdatedAt,
			OrderItems:      items,
		})
	}

	writeJSON(w, http.StatusOK, resp)
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
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req dto.V1QueryOrdersRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	result, err := c.orderService.GetOrders(r.Context(), models.QueryOrderItemsModel{
		IDs:               req.IDs,
		CustomerIDs:       req.CustomerIDs,
		Page:              req.Page,
		PageSize:          req.PageSize,
		IncludeOrderItems: req.IncludeOrderItems,
	})
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var resp dto.V1QueryOrdersResponse
	for _, o := range result {
		var items []dto.V1OrderItem
		for _, it := range o.OrderItems {
			items = append(items, dto.V1OrderItem{
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
			})
		}
		resp.Orders = append(resp.Orders, dto.V1Order{
			ID:              o.ID,
			CustomerID:      o.CustomerID,
			DeliveryAddress: o.DeliveryAddress,
			TotalPriceCents: o.TotalPriceCents,
			TotalPriceCurr:  o.TotalPriceCurr,
			CreatedAt:       o.CreatedAt,
			UpdatedAt:       o.UpdatedAt,
			OrderItems:      items,
		})
	}

	writeJSON(w, http.StatusOK, resp)
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

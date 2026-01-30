package handler

import (
	"enterprise-microservice-system/common/logger"
	"enterprise-microservice-system/common/response"
	"enterprise-microservice-system/services/order-service/internal/model"
	"enterprise-microservice-system/services/order-service/internal/service"
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// OrderHandler handles HTTP requests for orders
type OrderHandler struct {
	service service.OrderService
	logger  *logger.Logger
}

// NewOrderHandler creates a new order handler
func NewOrderHandler(service service.OrderService, logger *logger.Logger) *OrderHandler {
	return &OrderHandler{
		service: service,
		logger:  logger,
	}
}

// CreateOrder handles order creation
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req model.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid request body", zap.Error(err))
		response.Error(c, err)
		return
	}

	order, err := h.service.CreateOrder(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to create order", zap.Error(err))
		response.Error(c, err)
		return
	}

	h.logger.Info("Order created successfully", zap.Uint("order_id", order.ID))
	response.Created(c, order)
}

// GetOrder handles retrieving an order by ID
func (h *OrderHandler) GetOrder(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		h.logger.Warn("Invalid order ID", zap.Error(err))
		response.Error(c, err)
		return
	}

	order, err := h.service.GetOrder(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("Failed to get order", zap.Uint64("order_id", id), zap.Error(err))
		response.Error(c, err)
		return
	}

	response.Success(c, order)
}

// UpdateOrder handles updating an order
func (h *OrderHandler) UpdateOrder(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		h.logger.Warn("Invalid order ID", zap.Error(err))
		response.Error(c, err)
		return
	}

	var req model.UpdateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid request body", zap.Error(err))
		response.Error(c, err)
		return
	}

	order, err := h.service.UpdateOrder(c.Request.Context(), uint(id), &req)
	if err != nil {
		h.logger.Error("Failed to update order", zap.Uint64("order_id", id), zap.Error(err))
		response.Error(c, err)
		return
	}

	h.logger.Info("Order updated successfully", zap.Uint64("order_id", id))
	response.Success(c, order)
}

// DeleteOrder handles deleting an order
func (h *OrderHandler) DeleteOrder(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		h.logger.Warn("Invalid order ID", zap.Error(err))
		response.Error(c, err)
		return
	}

	if err := h.service.DeleteOrder(c.Request.Context(), uint(id)); err != nil {
		h.logger.Error("Failed to delete order", zap.Uint64("order_id", id), zap.Error(err))
		response.Error(c, err)
		return
	}

	h.logger.Info("Order deleted successfully", zap.Uint64("order_id", id))
	response.Success(c, gin.H{"message": "order deleted successfully"})
}

// ListOrders handles listing orders with pagination
func (h *OrderHandler) ListOrders(c *gin.Context) {
	var query model.ListOrdersQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		h.logger.Warn("Invalid query parameters", zap.Error(err))
		response.Error(c, err)
		return
	}

	orders, total, err := h.service.ListOrders(c.Request.Context(), &query)
	if err != nil {
		h.logger.Error("Failed to list orders", zap.Error(err))
		response.Error(c, err)
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(query.PageSize)))
	meta := &response.Meta{
		Page:       query.Page,
		PageSize:   query.PageSize,
		TotalPages: totalPages,
		TotalCount: total,
	}

	response.SuccessWithMeta(c, orders, meta)
}

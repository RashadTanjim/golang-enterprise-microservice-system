package handler

import (
	"enterprise-microservice-system/common/logger"
	"enterprise-microservice-system/common/middleware"
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
// @Summary Create a new order
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param order body model.CreateOrderRequest true "Order data"
// @Success 201 {object} response.Response{data=model.Order}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /orders [post]
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req model.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid request body", zap.Error(err))
		response.Error(c, err)
		return
	}

	actor := resolveActor(c)
	order, err := h.service.CreateOrder(c.Request.Context(), &req, actor)
	if err != nil {
		h.logger.Error("Failed to create order", zap.Error(err))
		response.Error(c, err)
		return
	}

	h.logger.Info("Order created successfully", zap.Uint("order_id", order.ID))
	response.Created(c, order)
}

// GetOrder handles retrieving an order by ID
// @Summary Get an order by ID
// @Tags orders
// @Produce json
// @Security BearerAuth
// @Param id path int true "Order ID"
// @Success 200 {object} response.Response{data=model.OrderWithUser}
// @Failure 404 {object} response.Response
// @Router /orders/{id} [get]
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
// @Summary Update an order
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Order ID"
// @Param order body model.UpdateOrderRequest true "Order update data"
// @Success 200 {object} response.Response{data=model.Order}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /orders/{id} [put]
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

	actor := resolveActor(c)
	order, err := h.service.UpdateOrder(c.Request.Context(), uint(id), &req, actor)
	if err != nil {
		h.logger.Error("Failed to update order", zap.Uint64("order_id", id), zap.Error(err))
		response.Error(c, err)
		return
	}

	h.logger.Info("Order updated successfully", zap.Uint64("order_id", id))
	response.Success(c, order)
}

// DeleteOrder handles deleting an order
// @Summary Delete an order
// @Tags orders
// @Produce json
// @Security BearerAuth
// @Param id path int true "Order ID"
// @Success 200 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /orders/{id} [delete]
func (h *OrderHandler) DeleteOrder(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		h.logger.Warn("Invalid order ID", zap.Error(err))
		response.Error(c, err)
		return
	}

	actor := resolveActor(c)
	if err := h.service.DeleteOrder(c.Request.Context(), uint(id), actor); err != nil {
		h.logger.Error("Failed to delete order", zap.Uint64("order_id", id), zap.Error(err))
		response.Error(c, err)
		return
	}

	h.logger.Info("Order deleted successfully", zap.Uint64("order_id", id))
	response.Success(c, gin.H{"message": "order deleted successfully"})
}

// ListOrders handles listing orders with pagination
// @Summary List orders
// @Tags orders
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Param user_id query int false "Filter by user ID"
// @Param order_status query string false "Filter by order status"
// @Param status query string false "Filter by record status (active/deleted)"
// @Param product_id query string false "Filter by product ID"
// @Success 200 {object} response.Response{data=[]model.Order}
// @Router /orders [get]
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

func resolveActor(c *gin.Context) string {
	if subject, ok := middleware.GetAuthSubject(c); ok && subject != "" {
		return subject
	}
	return "system"
}

package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"enterprise-microservice-system/common/cache"
	"enterprise-microservice-system/common/errors"
	"enterprise-microservice-system/services/order-service/internal/client"
	"enterprise-microservice-system/services/order-service/internal/model"
	"enterprise-microservice-system/services/order-service/internal/repository"

	"gorm.io/gorm"
)

// OrderService defines the business logic interface for orders
type OrderService interface {
	CreateOrder(ctx context.Context, req *model.CreateOrderRequest, actor string) (*model.OrderWithUser, error)
	GetOrder(ctx context.Context, id uint) (*model.OrderWithUser, error)
	UpdateOrder(ctx context.Context, id uint, req *model.UpdateOrderRequest, actor string) (*model.Order, error)
	DeleteOrder(ctx context.Context, id uint, actor string) error
	ListOrders(ctx context.Context, query *model.ListOrdersQuery) ([]*model.Order, int64, error)
}

// orderService implements OrderService
type orderService struct {
	repo       repository.OrderRepository
	userClient client.UserServiceClient
	cache      *cache.Cache
}

// NewOrderService creates a new order service
func NewOrderService(repo repository.OrderRepository, userClient client.UserServiceClient, cacheClient *cache.Cache) OrderService {
	return &orderService{
		repo:       repo,
		userClient: userClient,
		cache:      cacheClient,
	}
}

// CreateOrder creates a new order with user validation
func (s *orderService) CreateOrder(ctx context.Context, req *model.CreateOrderRequest, actor string) (*model.OrderWithUser, error) {
	// Validate user exists via circuit breaker protected call
	user, err := s.userClient.GetUser(ctx, req.UserID)
	if err != nil {
		// If circuit breaker is open or user service is down, still create order
		// but return without user data (graceful degradation)
		if appErr, ok := err.(*errors.AppError); ok {
			if appErr.Code == errors.ErrCodeCircuitOpen || appErr.Code == errors.ErrCodeServiceUnavail {
				// Create order without user validation
				order := &model.Order{
					UserID:      req.UserID,
					ProductID:   req.ProductID,
					Quantity:    req.Quantity,
					TotalPrice:  req.TotalPrice,
					OrderStatus: model.OrderStatusPending,
					Status:      model.OrderRecordStatusActive,
				}

				if actor == "" {
					actor = "system"
				}
				order.CreatedBy = actor
				order.UpdatedBy = actor

				if err := s.repo.Create(ctx, order); err != nil {
					return nil, errors.NewInternal("failed to create order", err)
				}

				result := &model.OrderWithUser{
					Order: *order,
					User:  nil, // User data unavailable
				}
				s.cacheSetOrder(ctx, result)
				return result, nil
			}
		}
		return nil, err
	}

	// Check if user is active
	if user.Status != "active" {
		return nil, errors.NewBadRequest("user is not active")
	}

	// Create order
	order := &model.Order{
		UserID:      req.UserID,
		ProductID:   req.ProductID,
		Quantity:    req.Quantity,
		TotalPrice:  req.TotalPrice,
		OrderStatus: model.OrderStatusPending,
		Status:      model.OrderRecordStatusActive,
	}

	if actor == "" {
		actor = "system"
	}
	order.CreatedBy = actor
	order.UpdatedBy = actor

	if err := s.repo.Create(ctx, order); err != nil {
		return nil, errors.NewInternal("failed to create order", err)
	}

	s.cacheSetOrder(ctx, &model.OrderWithUser{Order: *order, User: user})

	return &model.OrderWithUser{
		Order: *order,
		User:  user,
	}, nil
}

// GetOrder retrieves an order by ID with user data
func (s *orderService) GetOrder(ctx context.Context, id uint) (*model.OrderWithUser, error) {
	if cached := s.cacheGetOrder(ctx, id); cached != nil {
		return cached, nil
	}

	order, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewNotFound("order")
		}
		return nil, errors.NewInternal("failed to get order", err)
	}

	// Try to fetch user data (graceful degradation if user service is down)
	user, err := s.userClient.GetUser(ctx, order.UserID)
	if err != nil {
		// Log error but continue without user data
		user = nil
	}

	result := &model.OrderWithUser{
		Order: *order,
		User:  user,
	}
	s.cacheSetOrder(ctx, result)
	return result, nil
}

// UpdateOrder updates an order
func (s *orderService) UpdateOrder(ctx context.Context, id uint, req *model.UpdateOrderRequest, actor string) (*model.Order, error) {
	// Get existing order
	order, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewNotFound("order")
		}
		return nil, errors.NewInternal("failed to get order", err)
	}

	// Update fields if provided
	if req.OrderStatus != nil {
		order.OrderStatus = *req.OrderStatus
	}
	if actor == "" {
		actor = "system"
	}
	order.UpdatedBy = actor

	// Save updates
	if err := s.repo.Update(ctx, order); err != nil {
		return nil, errors.NewInternal("failed to update order", err)
	}

	s.cacheDeleteOrder(ctx, id)

	return order, nil
}

// DeleteOrder deletes an order
func (s *orderService) DeleteOrder(ctx context.Context, id uint, actor string) error {
	// Check if order exists
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.NewNotFound("order")
		}
		return errors.NewInternal("failed to get order", err)
	}

	if actor == "" {
		actor = "system"
	}

	// Delete order (status-based soft delete)
	if err := s.repo.Delete(ctx, id, actor); err != nil {
		return errors.NewInternal("failed to delete order", err)
	}

	s.cacheDeleteOrder(ctx, id)

	return nil
}

// ListOrders retrieves a paginated list of orders
func (s *orderService) ListOrders(ctx context.Context, query *model.ListOrdersQuery) ([]*model.Order, int64, error) {
	query.ApplyDefaults()

	if cachedOrders, cachedTotal, ok := s.cacheGetOrderList(ctx, query); ok {
		return cachedOrders, cachedTotal, nil
	}

	orders, total, err := s.repo.List(ctx, query)
	if err != nil {
		return nil, 0, errors.NewInternal("failed to list orders", err)
	}

	s.cacheSetOrderList(ctx, query, orders, total)
	return orders, total, nil
}

func (s *orderService) cacheGetOrder(ctx context.Context, id uint) *model.OrderWithUser {
	if s.cache == nil || !s.cache.Enabled() {
		return nil
	}

	var payload model.OrderWithUser
	found, err := s.cache.GetJSON(ctx, fmt.Sprintf("order:%d", id), &payload)
	if err != nil || !found {
		return nil
	}
	return &payload
}

func (s *orderService) cacheSetOrder(ctx context.Context, order *model.OrderWithUser) {
	if s.cache == nil || !s.cache.Enabled() || order == nil {
		return
	}

	_ = s.cache.SetJSON(ctx, fmt.Sprintf("order:%d", order.ID), order, 0)
}

func (s *orderService) cacheDeleteOrder(ctx context.Context, id uint) {
	if s.cache == nil || !s.cache.Enabled() {
		return
	}
	_ = s.cache.Delete(ctx, fmt.Sprintf("order:%d", id))
}

func (s *orderService) cacheGetOrderList(ctx context.Context, query *model.ListOrdersQuery) ([]*model.Order, int64, bool) {
	if s.cache == nil || !s.cache.Enabled() {
		return nil, 0, false
	}

	key := s.orderListCacheKey(query)
	var payload struct {
		Orders []*model.Order `json:"orders"`
		Total  int64          `json:"total"`
	}

	found, err := s.cache.GetJSON(ctx, key, &payload)
	if err != nil || !found {
		return nil, 0, false
	}
	return payload.Orders, payload.Total, true
}

func (s *orderService) cacheSetOrderList(ctx context.Context, query *model.ListOrdersQuery, orders []*model.Order, total int64) {
	if s.cache == nil || !s.cache.Enabled() {
		return
	}

	key := s.orderListCacheKey(query)
	payload := struct {
		Orders []*model.Order `json:"orders"`
		Total  int64          `json:"total"`
	}{
		Orders: orders,
		Total:  total,
	}

	_ = s.cache.SetJSON(ctx, key, payload, 60*time.Second)
}

func (s *orderService) orderListCacheKey(query *model.ListOrdersQuery) string {
	userID := "any"
	if query.UserID != nil {
		userID = fmt.Sprintf("%d", *query.UserID)
	}

	orderStatus := "all"
	if query.OrderStatus != nil {
		orderStatus = string(*query.OrderStatus)
	}

	recordStatus := "any"
	if query.Status != nil {
		recordStatus = *query.Status
	}

	product := strings.TrimSpace(query.ProductID)
	if product == "" {
		product = "all"
	}

	return fmt.Sprintf(
		"orders:list:p%d:ps%d:user:%s:order_status:%s:record_status:%s:product:%s",
		query.Page,
		query.PageSize,
		userID,
		orderStatus,
		recordStatus,
		product,
	)
}

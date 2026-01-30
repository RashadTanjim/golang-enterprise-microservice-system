package service

import (
	"context"
	"enterprise-microservice-system/common/errors"
	"enterprise-microservice-system/services/order-service/internal/client"
	"enterprise-microservice-system/services/order-service/internal/model"
	"enterprise-microservice-system/services/order-service/internal/repository"

	"gorm.io/gorm"
)

// OrderService defines the business logic interface for orders
type OrderService interface {
	CreateOrder(ctx context.Context, req *model.CreateOrderRequest) (*model.OrderWithUser, error)
	GetOrder(ctx context.Context, id uint) (*model.OrderWithUser, error)
	UpdateOrder(ctx context.Context, id uint, req *model.UpdateOrderRequest) (*model.Order, error)
	DeleteOrder(ctx context.Context, id uint) error
	ListOrders(ctx context.Context, query *model.ListOrdersQuery) ([]*model.Order, int64, error)
}

// orderService implements OrderService
type orderService struct {
	repo       repository.OrderRepository
	userClient *client.UserClient
}

// NewOrderService creates a new order service
func NewOrderService(repo repository.OrderRepository, userClient *client.UserClient) OrderService {
	return &orderService{
		repo:       repo,
		userClient: userClient,
	}
}

// CreateOrder creates a new order with user validation
func (s *orderService) CreateOrder(ctx context.Context, req *model.CreateOrderRequest) (*model.OrderWithUser, error) {
	// Validate user exists via circuit breaker protected call
	user, err := s.userClient.GetUser(ctx, req.UserID)
	if err != nil {
		// If circuit breaker is open or user service is down, still create order
		// but return without user data (graceful degradation)
		if errors, ok := err.(*errors.AppError); ok {
			if errors.Code == errors.ErrCodeCircuitOpen || errors.Code == errors.ErrCodeServiceUnavail {
				// Create order without user validation
				order := &model.Order{
					UserID:     req.UserID,
					ProductID:  req.ProductID,
					Quantity:   req.Quantity,
					TotalPrice: req.TotalPrice,
					Status:     model.OrderStatusPending,
				}

				if err := s.repo.Create(ctx, order); err != nil {
					return nil, errors.NewInternal("failed to create order", err)
				}

				return &model.OrderWithUser{
					Order: *order,
					User:  nil, // User data unavailable
				}, nil
			}
		}
		return nil, err
	}

	// Check if user is active
	if !user.Active {
		return nil, errors.NewBadRequest("user is not active")
	}

	// Create order
	order := &model.Order{
		UserID:     req.UserID,
		ProductID:  req.ProductID,
		Quantity:   req.Quantity,
		TotalPrice: req.TotalPrice,
		Status:     model.OrderStatusPending,
	}

	if err := s.repo.Create(ctx, order); err != nil {
		return nil, errors.NewInternal("failed to create order", err)
	}

	return &model.OrderWithUser{
		Order: *order,
		User:  user,
	}, nil
}

// GetOrder retrieves an order by ID with user data
func (s *orderService) GetOrder(ctx context.Context, id uint) (*model.OrderWithUser, error) {
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

	return &model.OrderWithUser{
		Order: *order,
		User:  user,
	}, nil
}

// UpdateOrder updates an order
func (s *orderService) UpdateOrder(ctx context.Context, id uint, req *model.UpdateOrderRequest) (*model.Order, error) {
	// Get existing order
	order, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewNotFound("order")
		}
		return nil, errors.NewInternal("failed to get order", err)
	}

	// Update fields if provided
	if req.Status != nil {
		order.Status = *req.Status
	}

	// Save updates
	if err := s.repo.Update(ctx, order); err != nil {
		return nil, errors.NewInternal("failed to update order", err)
	}

	return order, nil
}

// DeleteOrder deletes an order
func (s *orderService) DeleteOrder(ctx context.Context, id uint) error {
	// Check if order exists
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.NewNotFound("order")
		}
		return errors.NewInternal("failed to get order", err)
	}

	// Delete order
	if err := s.repo.Delete(ctx, id); err != nil {
		return errors.NewInternal("failed to delete order", err)
	}

	return nil
}

// ListOrders retrieves a paginated list of orders
func (s *orderService) ListOrders(ctx context.Context, query *model.ListOrdersQuery) ([]*model.Order, int64, error) {
	query.ApplyDefaults()

	orders, total, err := s.repo.List(ctx, query)
	if err != nil {
		return nil, 0, errors.NewInternal("failed to list orders", err)
	}

	return orders, total, nil
}

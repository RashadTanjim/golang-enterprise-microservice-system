package repository

import (
	"context"
	"time"

	"enterprise-microservice-system/services/order-service/internal/model"

	"gorm.io/gorm"
)

// OrderRepository defines the interface for order data operations
type OrderRepository interface {
	Create(ctx context.Context, order *model.Order) error
	FindByID(ctx context.Context, id uint) (*model.Order, error)
	Update(ctx context.Context, order *model.Order) error
	Delete(ctx context.Context, id uint, updatedBy string) error
	List(ctx context.Context, query *model.ListOrdersQuery) ([]*model.Order, int64, error)
	FindByUserID(ctx context.Context, userID uint) ([]*model.Order, error)
}

// orderRepository implements OrderRepository
type orderRepository struct {
	db *gorm.DB
}

// NewOrderRepository creates a new order repository
func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

// Create creates a new order
func (r *orderRepository) Create(ctx context.Context, order *model.Order) error {
	return r.db.WithContext(ctx).Create(order).Error
}

// FindByID finds an order by ID
func (r *orderRepository) FindByID(ctx context.Context, id uint) (*model.Order, error) {
	var order model.Order
	err := r.db.WithContext(ctx).
		Where("status <> ?", model.OrderRecordStatusDeleted).
		First(&order, id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// Update updates an order
func (r *orderRepository) Update(ctx context.Context, order *model.Order) error {
	return r.db.WithContext(ctx).Save(order).Error
}

// Delete soft deletes an order
func (r *orderRepository) Delete(ctx context.Context, id uint, updatedBy string) error {
	if updatedBy == "" {
		updatedBy = "system"
	}

	return r.db.WithContext(ctx).
		Model(&model.Order{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     model.OrderRecordStatusDeleted,
			"updated_by": updatedBy,
			"updated_at": time.Now().UTC(),
		}).Error
}

// List retrieves a paginated list of orders
func (r *orderRepository) List(ctx context.Context, query *model.ListOrdersQuery) ([]*model.Order, int64, error) {
	var orders []*model.Order
	var total int64

	db := r.db.WithContext(ctx).Model(&model.Order{})

	// Apply filters
	if query.UserID != nil {
		db = db.Where("user_id = ?", *query.UserID)
	}

	if query.OrderStatus != nil {
		db = db.Where("order_status = ?", *query.OrderStatus)
	}

	if query.ProductID != "" {
		db = db.Where("product_id = ?", query.ProductID)
	}

	if query.Status != nil {
		db = db.Where("status = ?", *query.Status)
	} else {
		db = db.Where("status <> ?", model.OrderRecordStatusDeleted)
	}

	// Count total
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	err := db.Offset(query.Offset()).
		Limit(query.PageSize).
		Order("created_at DESC").
		Find(&orders).Error

	if err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

// FindByUserID finds all orders for a user
func (r *orderRepository) FindByUserID(ctx context.Context, userID uint) ([]*model.Order, error) {
	var orders []*model.Order
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND status <> ?", userID, model.OrderRecordStatusDeleted).
		Order("created_at DESC").
		Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}

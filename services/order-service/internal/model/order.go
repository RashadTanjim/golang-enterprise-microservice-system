package model

import "time"

// OrderStatus represents the status of an order
type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusConfirmed OrderStatus = "confirmed"
	OrderStatusShipped   OrderStatus = "shipped"
	OrderStatusDelivered OrderStatus = "delivered"
	OrderStatusCancelled OrderStatus = "cancelled"
)

const (
	OrderRecordStatusActive  = "active"
	OrderRecordStatusDeleted = "deleted"
)

// Order represents an order in the system
type Order struct {
	ID          uint        `gorm:"primarykey" json:"id"`
	UserID      uint        `gorm:"not null;index" json:"user_id"`
	ProductID   string      `gorm:"not null" json:"product_id"`
	Quantity    int         `gorm:"not null" json:"quantity"`
	TotalPrice  float64     `gorm:"not null" json:"total_price"`
	OrderStatus OrderStatus `gorm:"column:order_status;type:varchar(20);not null;default:'pending'" json:"order_status"`
	Status      string      `gorm:"type:varchar(20);not null;default:'active';index" json:"status"`
	CreatedBy   string      `gorm:"type:varchar(100);not null;default:'system'" json:"created_by"`
	UpdatedBy   string      `gorm:"type:varchar(100);not null;default:'system'" json:"updated_by"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

// TableName overrides the default table name
func (Order) TableName() string {
	return "orders"
}

// CreateOrderRequest represents the request to create an order
type CreateOrderRequest struct {
	UserID     uint    `json:"user_id" binding:"required"`
	ProductID  string  `json:"product_id" binding:"required"`
	Quantity   int     `json:"quantity" binding:"required,min=1"`
	TotalPrice float64 `json:"total_price" binding:"required,min=0"`
}

// UpdateOrderRequest represents the request to update an order
type UpdateOrderRequest struct {
	OrderStatus *OrderStatus `json:"order_status" binding:"omitempty,oneof=pending confirmed shipped delivered cancelled"`
}

// ListOrdersQuery represents query parameters for listing orders
type ListOrdersQuery struct {
	Page        int          `form:"page" binding:"omitempty,min=1"`
	PageSize    int          `form:"page_size" binding:"omitempty,min=1,max=100"`
	UserID      *uint        `form:"user_id"`
	OrderStatus *OrderStatus `form:"order_status" binding:"omitempty,oneof=pending confirmed shipped delivered cancelled"`
	Status      *string      `form:"status" binding:"omitempty,oneof=active deleted"`
	ProductID   string       `form:"product_id"`
}

// ApplyDefaults applies default values to the query
func (q *ListOrdersQuery) ApplyDefaults() {
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.PageSize <= 0 {
		q.PageSize = 10
	}
}

// Offset calculates the offset for pagination
func (q *ListOrdersQuery) Offset() int {
	return (q.Page - 1) * q.PageSize
}

// User represents user data from user service
type User struct {
	ID     uint   `json:"id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	Age    int    `json:"age"`
	Status string `json:"status"`
}

// OrderWithUser combines order with user data
type OrderWithUser struct {
	Order
	User *User `json:"user,omitempty"`
}

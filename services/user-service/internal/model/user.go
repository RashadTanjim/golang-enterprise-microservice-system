package model

import "time"

const (
	UserStatusActive   = "active"
	UserStatusInactive = "inactive"
	UserStatusDeleted  = "deleted"
)

// User represents a user in the system
type User struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Email     string    `gorm:"uniqueIndex;not null" json:"email"`
	Name      string    `gorm:"not null" json:"name"`
	Age       int       `json:"age"`
	Status    string    `gorm:"type:varchar(20);not null;default:'active';index" json:"status"`
	CreatedBy string    `gorm:"type:varchar(100);not null;default:'system'" json:"created_by"`
	UpdatedBy string    `gorm:"type:varchar(100);not null;default:'system'" json:"updated_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName overrides the default table name
func (User) TableName() string {
	return "users"
}

// CreateUserRequest represents the request to create a user
type CreateUserRequest struct {
	Email string `json:"email" binding:"required,email"`
	Name  string `json:"name" binding:"required,min=2,max=100"`
	Age   int    `json:"age" binding:"required,min=1,max=150"`
}

// UpdateUserRequest represents the request to update a user
type UpdateUserRequest struct {
	Name   *string `json:"name" binding:"omitempty,min=2,max=100"`
	Age    *int    `json:"age" binding:"omitempty,min=1,max=150"`
	Status *string `json:"status" binding:"omitempty,oneof=active inactive"`
}

// ListUsersQuery represents query parameters for listing users
type ListUsersQuery struct {
	Page     int     `form:"page" binding:"omitempty,min=1"`
	PageSize int     `form:"page_size" binding:"omitempty,min=1,max=100"`
	Search   string  `form:"search"`
	Status   *string `form:"status" binding:"omitempty,oneof=active inactive deleted"`
}

// ApplyDefaults applies default values to the query
func (q *ListUsersQuery) ApplyDefaults() {
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.PageSize <= 0 {
		q.PageSize = 10
	}
}

// Offset calculates the offset for pagination
func (q *ListUsersQuery) Offset() int {
	return (q.Page - 1) * q.PageSize
}

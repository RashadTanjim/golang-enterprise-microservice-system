package model

import "time"

const (
	RepositoryStatusActive   = "active"
	RepositoryStatusInactive = "inactive"
	RepositoryStatusDeleted  = "deleted"
)

// Repository represents a repository/project in the system
type Repository struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	Name        string    `gorm:"not null;uniqueIndex" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	OwnerID     uint      `gorm:"not null;index" json:"owner_id"`
	Visibility  string    `gorm:"type:varchar(20);default:'private'" json:"visibility"`
	URL         string    `json:"url"`
	Status      string    `gorm:"type:varchar(20);not null;default:'active';index" json:"status"`
	CreatedBy   string    `gorm:"type:varchar(100);not null;default:'system'" json:"created_by"`
	UpdatedBy   string    `gorm:"type:varchar(100);not null;default:'system'" json:"updated_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName overrides the default table name
func (Repository) TableName() string {
	return "repositories"
}

// CreateRepositoryRequest represents the request to create a repository
type CreateRepositoryRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=100"`
	Description string `json:"description" binding:"omitempty,max=1000"`
	OwnerID     uint   `json:"owner_id" binding:"required"`
	Visibility  string `json:"visibility" binding:"omitempty,oneof=public private"`
	URL         string `json:"url" binding:"omitempty,url"`
}

// UpdateRepositoryRequest represents the request to update a repository
type UpdateRepositoryRequest struct {
	Name        *string `json:"name" binding:"omitempty,min=2,max=100"`
	Description *string `json:"description" binding:"omitempty,max=1000"`
	Visibility  *string `json:"visibility" binding:"omitempty,oneof=public private"`
	URL         *string `json:"url" binding:"omitempty,url"`
	Status      *string `json:"status" binding:"omitempty,oneof=active inactive"`
}

// ListRepositoriesQuery represents query parameters for listing repositories
type ListRepositoriesQuery struct {
	Page       int     `form:"page" binding:"omitempty,min=1"`
	PageSize   int     `form:"page_size" binding:"omitempty,min=1,max=100"`
	Search     string  `form:"search"`
	OwnerID    *uint   `form:"owner_id"`
	Visibility string  `form:"visibility" binding:"omitempty,oneof=public private"`
	Status     *string `form:"status" binding:"omitempty,oneof=active inactive deleted"`
}

// ApplyDefaults applies default values to the query
func (q *ListRepositoriesQuery) ApplyDefaults() {
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.PageSize <= 0 {
		q.PageSize = 10
	}
}

// Offset calculates the offset for pagination
func (q *ListRepositoriesQuery) Offset() int {
	return (q.Page - 1) * q.PageSize
}

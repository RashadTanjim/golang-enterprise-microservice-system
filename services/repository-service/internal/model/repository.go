package model

import (
	"time"

	"gorm.io/gorm"
)

// Repository represents a repository/project in the system
type Repository struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	Name        string         `gorm:"not null;uniqueIndex" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	OwnerID     uint           `gorm:"not null;index" json:"owner_id"`
	Visibility  string         `gorm:"type:varchar(20);default:'private'" json:"visibility"`
	URL         string         `json:"url"`
	Active      bool           `gorm:"default:true" json:"active"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
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
	Active      *bool   `json:"active"`
}

// ListRepositoriesQuery represents query parameters for listing repositories
type ListRepositoriesQuery struct {
	Page       int    `form:"page" binding:"omitempty,min=1"`
	PageSize   int    `form:"page_size" binding:"omitempty,min=1,max=100"`
	Search     string `form:"search"`
	OwnerID    *uint  `form:"owner_id"`
	Visibility string `form:"visibility" binding:"omitempty,oneof=public private"`
	Active     *bool  `form:"active"`
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

package model

import "time"

const (
	AuditLogStatusActive  = "active"
	AuditLogStatusDeleted = "deleted"
)

// AuditLog represents an audit log entry in the system
// capturing who did what on which resource.
type AuditLog struct {
	ID           uint      `gorm:"primarykey" json:"id"`
	Actor        string    `gorm:"type:varchar(120);not null;index" json:"actor"`
	Action       string    `gorm:"type:varchar(100);not null;index" json:"action"`
	ResourceType string    `gorm:"type:varchar(100);not null;index" json:"resource_type"`
	ResourceID   string    `gorm:"type:varchar(100);not null;index" json:"resource_id"`
	Description  string    `gorm:"type:text" json:"description"`
	Metadata     string    `gorm:"type:text" json:"metadata"`
	Status       string    `gorm:"type:varchar(20);not null;default:'active';index" json:"status"`
	CreatedBy    string    `gorm:"type:varchar(100);not null;default:'system'" json:"created_by"`
	UpdatedBy    string    `gorm:"type:varchar(100);not null;default:'system'" json:"updated_by"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// TableName overrides the default table name
func (AuditLog) TableName() string {
	return "audit_logs"
}

// CreateAuditLogRequest represents the request to create an audit log entry
// Actor is optional and will default to the authenticated subject.
type CreateAuditLogRequest struct {
	Actor        string `json:"actor" binding:"omitempty,max=120"`
	Action       string `json:"action" binding:"required,min=2,max=100"`
	ResourceType string `json:"resource_type" binding:"required,min=2,max=100"`
	ResourceID   string `json:"resource_id" binding:"required,min=1,max=100"`
	Description  string `json:"description" binding:"omitempty,max=2000"`
	Metadata     string `json:"metadata" binding:"omitempty,max=5000"`
}

// UpdateAuditLogRequest represents the request to update an audit log entry
// (for demo purposes, allows updating description/metadata/status).
type UpdateAuditLogRequest struct {
	Description *string `json:"description" binding:"omitempty,max=2000"`
	Metadata    *string `json:"metadata" binding:"omitempty,max=5000"`
	Status      *string `json:"status" binding:"omitempty,oneof=active deleted"`
}

// ListAuditLogsQuery represents query parameters for listing audit logs
// Status filters visibility (active/deleted).
type ListAuditLogsQuery struct {
	Page         int     `form:"page" binding:"omitempty,min=1"`
	PageSize     int     `form:"page_size" binding:"omitempty,min=1,max=100"`
	Search       string  `form:"search"`
	Actor        string  `form:"actor"`
	Action       string  `form:"action"`
	ResourceType string  `form:"resource_type"`
	ResourceID   string  `form:"resource_id"`
	Status       *string `form:"status" binding:"omitempty,oneof=active deleted"`
}

// ApplyDefaults applies default values to the query
func (q *ListAuditLogsQuery) ApplyDefaults() {
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.PageSize <= 0 {
		q.PageSize = 10
	}
}

// Offset calculates the offset for pagination
func (q *ListAuditLogsQuery) Offset() int {
	return (q.Page - 1) * q.PageSize
}

package repository

import (
	"context"
	"time"

	"enterprise-microservice-system/services/audit-log-service/internal/model"

	"gorm.io/gorm"
)

// AuditLogRepository defines the interface for audit log data operations
type AuditLogRepository interface {
	Create(ctx context.Context, entry *model.AuditLog) error
	FindByID(ctx context.Context, id uint) (*model.AuditLog, error)
	Update(ctx context.Context, entry *model.AuditLog) error
	Delete(ctx context.Context, id uint, updatedBy string) error
	List(ctx context.Context, query *model.ListAuditLogsQuery) ([]*model.AuditLog, int64, error)
}

// auditLogRepository implements AuditLogRepository
type auditLogRepository struct {
	db *gorm.DB
}

// NewAuditLogRepository creates a new audit log repository
func NewAuditLogRepository(db *gorm.DB) AuditLogRepository {
	return &auditLogRepository{db: db}
}

// Create creates a new audit log entry
func (r *auditLogRepository) Create(ctx context.Context, entry *model.AuditLog) error {
	return r.db.WithContext(ctx).Create(entry).Error
}

// FindByID finds an audit log entry by ID
func (r *auditLogRepository) FindByID(ctx context.Context, id uint) (*model.AuditLog, error) {
	var entry model.AuditLog
	err := r.db.WithContext(ctx).
		Where("status <> ?", model.AuditLogStatusDeleted).
		First(&entry, id).Error
	if err != nil {
		return nil, err
	}
	return &entry, nil
}

// Update updates an audit log entry
func (r *auditLogRepository) Update(ctx context.Context, entry *model.AuditLog) error {
	return r.db.WithContext(ctx).Save(entry).Error
}

// Delete soft deletes an audit log entry
func (r *auditLogRepository) Delete(ctx context.Context, id uint, updatedBy string) error {
	if updatedBy == "" {
		updatedBy = "system"
	}

	return r.db.WithContext(ctx).
		Model(&model.AuditLog{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     model.AuditLogStatusDeleted,
			"updated_by": updatedBy,
			"updated_at": time.Now().UTC(),
		}).Error
}

// List retrieves a paginated list of audit log entries
func (r *auditLogRepository) List(ctx context.Context, query *model.ListAuditLogsQuery) ([]*model.AuditLog, int64, error) {
	var entries []*model.AuditLog
	var total int64

	db := r.db.WithContext(ctx).Model(&model.AuditLog{})

	if query.Search != "" {
		searchPattern := "%" + query.Search + "%"
		db = db.Where("description ILIKE ? OR metadata ILIKE ?", searchPattern, searchPattern)
	}

	if query.Actor != "" {
		db = db.Where("actor = ?", query.Actor)
	}

	if query.Action != "" {
		db = db.Where("action = ?", query.Action)
	}

	if query.ResourceType != "" {
		db = db.Where("resource_type = ?", query.ResourceType)
	}

	if query.ResourceID != "" {
		db = db.Where("resource_id = ?", query.ResourceID)
	}

	if query.Status != nil {
		db = db.Where("status = ?", *query.Status)
	} else {
		db = db.Where("status <> ?", model.AuditLogStatusDeleted)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Offset(query.Offset()).
		Limit(query.PageSize).
		Order("created_at DESC").
		Find(&entries).Error
	if err != nil {
		return nil, 0, err
	}

	return entries, total, nil
}

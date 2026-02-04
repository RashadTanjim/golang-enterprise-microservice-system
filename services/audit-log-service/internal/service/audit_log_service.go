package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"enterprise-microservice-system/common/cache"
	"enterprise-microservice-system/common/errors"
	"enterprise-microservice-system/services/audit-log-service/internal/model"
	"enterprise-microservice-system/services/audit-log-service/internal/repository"

	"gorm.io/gorm"
)

// AuditLogService defines the business logic interface for audit logs
type AuditLogService interface {
	CreateAuditLog(ctx context.Context, req *model.CreateAuditLogRequest, actor string) (*model.AuditLog, error)
	GetAuditLog(ctx context.Context, id uint) (*model.AuditLog, error)
	UpdateAuditLog(ctx context.Context, id uint, req *model.UpdateAuditLogRequest, actor string) (*model.AuditLog, error)
	DeleteAuditLog(ctx context.Context, id uint, actor string) error
	ListAuditLogs(ctx context.Context, query *model.ListAuditLogsQuery) ([]*model.AuditLog, int64, error)
}

// auditLogService implements AuditLogService
type auditLogService struct {
	repo  repository.AuditLogRepository
	cache *cache.Cache
}

// NewAuditLogService creates a new audit log service
func NewAuditLogService(repo repository.AuditLogRepository, cacheClient *cache.Cache) AuditLogService {
	return &auditLogService{
		repo:  repo,
		cache: cacheClient,
	}
}

// CreateAuditLog creates a new audit log entry
func (s *auditLogService) CreateAuditLog(ctx context.Context, req *model.CreateAuditLogRequest, actor string) (*model.AuditLog, error) {
	if actor == "" {
		actor = "system"
	}

	entry := &model.AuditLog{
		Actor:        req.Actor,
		Action:       req.Action,
		ResourceType: req.ResourceType,
		ResourceID:   req.ResourceID,
		Description:  req.Description,
		Metadata:     req.Metadata,
		Status:       model.AuditLogStatusActive,
		CreatedBy:    actor,
		UpdatedBy:    actor,
	}

	if entry.Actor == "" {
		entry.Actor = actor
	}

	if err := s.repo.Create(ctx, entry); err != nil {
		return nil, errors.NewInternal("failed to create audit log", err)
	}

	s.cacheSetAuditLog(ctx, entry)
	return entry, nil
}

// GetAuditLog retrieves an audit log entry by ID
func (s *auditLogService) GetAuditLog(ctx context.Context, id uint) (*model.AuditLog, error) {
	if cached := s.cacheGetAuditLog(ctx, id); cached != nil {
		return cached, nil
	}

	entry, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewNotFound("audit log")
		}
		return nil, errors.NewInternal("failed to get audit log", err)
	}

	s.cacheSetAuditLog(ctx, entry)
	return entry, nil
}

// UpdateAuditLog updates an audit log entry
func (s *auditLogService) UpdateAuditLog(ctx context.Context, id uint, req *model.UpdateAuditLogRequest, actor string) (*model.AuditLog, error) {
	entry, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewNotFound("audit log")
		}
		return nil, errors.NewInternal("failed to get audit log", err)
	}

	if req.Description != nil {
		entry.Description = *req.Description
	}
	if req.Metadata != nil {
		entry.Metadata = *req.Metadata
	}
	if req.Status != nil {
		entry.Status = *req.Status
	}
	if actor == "" {
		actor = "system"
	}
	entry.UpdatedBy = actor

	if err := s.repo.Update(ctx, entry); err != nil {
		return nil, errors.NewInternal("failed to update audit log", err)
	}

	s.cacheSetAuditLog(ctx, entry)
	return entry, nil
}

// DeleteAuditLog deletes an audit log entry
func (s *auditLogService) DeleteAuditLog(ctx context.Context, id uint, actor string) error {
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.NewNotFound("audit log")
		}
		return errors.NewInternal("failed to get audit log", err)
	}

	if actor == "" {
		actor = "system"
	}

	if err := s.repo.Delete(ctx, id, actor); err != nil {
		return errors.NewInternal("failed to delete audit log", err)
	}

	s.cacheDeleteAuditLog(ctx, id)
	return nil
}

// ListAuditLogs retrieves a paginated list of audit logs
func (s *auditLogService) ListAuditLogs(ctx context.Context, query *model.ListAuditLogsQuery) ([]*model.AuditLog, int64, error) {
	query.ApplyDefaults()

	if cachedEntries, cachedTotal, ok := s.cacheGetAuditLogList(ctx, query); ok {
		return cachedEntries, cachedTotal, nil
	}

	entries, total, err := s.repo.List(ctx, query)
	if err != nil {
		return nil, 0, errors.NewInternal("failed to list audit logs", err)
	}

	s.cacheSetAuditLogList(ctx, query, entries, total)
	return entries, total, nil
}

func (s *auditLogService) cacheGetAuditLog(ctx context.Context, id uint) *model.AuditLog {
	if s.cache == nil || !s.cache.Enabled() {
		return nil
	}

	var entry model.AuditLog
	found, err := s.cache.GetJSON(ctx, fmt.Sprintf("audit-log:%d", id), &entry)
	if err != nil || !found {
		return nil
	}
	return &entry
}

func (s *auditLogService) cacheSetAuditLog(ctx context.Context, entry *model.AuditLog) {
	if s.cache == nil || !s.cache.Enabled() || entry == nil {
		return
	}

	_ = s.cache.SetJSON(ctx, fmt.Sprintf("audit-log:%d", entry.ID), entry, 0)
}

func (s *auditLogService) cacheDeleteAuditLog(ctx context.Context, id uint) {
	if s.cache == nil || !s.cache.Enabled() {
		return
	}
	_ = s.cache.Delete(ctx, fmt.Sprintf("audit-log:%d", id))
}

func (s *auditLogService) cacheGetAuditLogList(ctx context.Context, query *model.ListAuditLogsQuery) ([]*model.AuditLog, int64, bool) {
	if s.cache == nil || !s.cache.Enabled() {
		return nil, 0, false
	}

	key := s.auditLogListCacheKey(query)
	var payload struct {
		Entries []*model.AuditLog `json:"entries"`
		Total   int64             `json:"total"`
	}

	found, err := s.cache.GetJSON(ctx, key, &payload)
	if err != nil || !found {
		return nil, 0, false
	}
	return payload.Entries, payload.Total, true
}

func (s *auditLogService) cacheSetAuditLogList(ctx context.Context, query *model.ListAuditLogsQuery, entries []*model.AuditLog, total int64) {
	if s.cache == nil || !s.cache.Enabled() {
		return
	}

	key := s.auditLogListCacheKey(query)
	payload := struct {
		Entries []*model.AuditLog `json:"entries"`
		Total   int64             `json:"total"`
	}{
		Entries: entries,
		Total:   total,
	}

	_ = s.cache.SetJSON(ctx, key, payload, 60*time.Second)
}

func (s *auditLogService) auditLogListCacheKey(query *model.ListAuditLogsQuery) string {
	status := "any"
	if query.Status != nil {
		status = *query.Status
	}

	actor := strings.TrimSpace(query.Actor)
	if actor == "" {
		actor = "any"
	}

	action := strings.TrimSpace(query.Action)
	if action == "" {
		action = "any"
	}

	resourceType := strings.TrimSpace(query.ResourceType)
	if resourceType == "" {
		resourceType = "any"
	}

	resourceID := strings.TrimSpace(query.ResourceID)
	if resourceID == "" {
		resourceID = "any"
	}

	search := strings.TrimSpace(query.Search)
	if search == "" {
		search = "all"
	}

	return fmt.Sprintf(
		"audit-logs:list:p%d:ps%d:actor:%s:action:%s:rtype:%s:rid:%s:status:%s:search:%s",
		query.Page,
		query.PageSize,
		actor,
		action,
		resourceType,
		resourceID,
		status,
		search,
	)
}

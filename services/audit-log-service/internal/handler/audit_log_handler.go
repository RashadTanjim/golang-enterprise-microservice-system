package handler

import (
	"enterprise-microservice-system/common/logger"
	"enterprise-microservice-system/common/middleware"
	"enterprise-microservice-system/common/response"
	"enterprise-microservice-system/services/audit-log-service/internal/model"
	"enterprise-microservice-system/services/audit-log-service/internal/service"
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// AuditLogHandler handles HTTP requests for audit logs
type AuditLogHandler struct {
	service service.AuditLogService
	logger  *logger.Logger
}

// NewAuditLogHandler creates a new audit log handler
func NewAuditLogHandler(service service.AuditLogService, logger *logger.Logger) *AuditLogHandler {
	return &AuditLogHandler{
		service: service,
		logger:  logger,
	}
}

// CreateAuditLog handles audit log creation
// @Summary Create a new audit log entry
// @Tags audit-logs
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param audit_log body model.CreateAuditLogRequest true "Audit log data"
// @Success 201 {object} response.Response{data=model.AuditLog}
// @Failure 400 {object} response.Response
// @Router /audit-logs [post]
func (h *AuditLogHandler) CreateAuditLog(c *gin.Context) {
	var req model.CreateAuditLogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid request body", zap.Error(err))
		response.Error(c, err)
		return
	}

	actor := resolveActor(c)
	entry, err := h.service.CreateAuditLog(c.Request.Context(), &req, actor)
	if err != nil {
		h.logger.Error("Failed to create audit log", zap.Error(err))
		response.Error(c, err)
		return
	}

	h.logger.Info("Audit log created successfully", zap.Uint("audit_log_id", entry.ID))
	response.Created(c, entry)
}

// GetAuditLog handles retrieving an audit log by ID
// @Summary Get an audit log entry by ID
// @Tags audit-logs
// @Produce json
// @Security BearerAuth
// @Param id path int true "Audit log ID"
// @Success 200 {object} response.Response{data=model.AuditLog}
// @Failure 404 {object} response.Response
// @Router /audit-logs/{id} [get]
func (h *AuditLogHandler) GetAuditLog(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		h.logger.Warn("Invalid audit log ID", zap.Error(err))
		response.Error(c, err)
		return
	}

	entry, err := h.service.GetAuditLog(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("Failed to get audit log", zap.Uint64("audit_log_id", id), zap.Error(err))
		response.Error(c, err)
		return
	}

	response.Success(c, entry)
}

// UpdateAuditLog handles updating an audit log entry
// @Summary Update an audit log entry
// @Tags audit-logs
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Audit log ID"
// @Param audit_log body model.UpdateAuditLogRequest true "Audit log update data"
// @Success 200 {object} response.Response{data=model.AuditLog}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /audit-logs/{id} [put]
func (h *AuditLogHandler) UpdateAuditLog(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		h.logger.Warn("Invalid audit log ID", zap.Error(err))
		response.Error(c, err)
		return
	}

	var req model.UpdateAuditLogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid request body", zap.Error(err))
		response.Error(c, err)
		return
	}

	actor := resolveActor(c)
	entry, err := h.service.UpdateAuditLog(c.Request.Context(), uint(id), &req, actor)
	if err != nil {
		h.logger.Error("Failed to update audit log", zap.Uint64("audit_log_id", id), zap.Error(err))
		response.Error(c, err)
		return
	}

	h.logger.Info("Audit log updated successfully", zap.Uint64("audit_log_id", id))
	response.Success(c, entry)
}

// DeleteAuditLog handles deleting an audit log entry
// @Summary Delete an audit log entry
// @Tags audit-logs
// @Produce json
// @Security BearerAuth
// @Param id path int true "Audit log ID"
// @Success 200 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /audit-logs/{id} [delete]
func (h *AuditLogHandler) DeleteAuditLog(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		h.logger.Warn("Invalid audit log ID", zap.Error(err))
		response.Error(c, err)
		return
	}

	actor := resolveActor(c)
	if err := h.service.DeleteAuditLog(c.Request.Context(), uint(id), actor); err != nil {
		h.logger.Error("Failed to delete audit log", zap.Uint64("audit_log_id", id), zap.Error(err))
		response.Error(c, err)
		return
	}

	h.logger.Info("Audit log deleted successfully", zap.Uint64("audit_log_id", id))
	response.Success(c, gin.H{"message": "audit log deleted successfully"})
}

// ListAuditLogs handles listing audit log entries with pagination
// @Summary List audit log entries
// @Tags audit-logs
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Param search query string false "Search term"
// @Param actor query string false "Filter by actor"
// @Param action query string false "Filter by action"
// @Param resource_type query string false "Filter by resource type"
// @Param resource_id query string false "Filter by resource ID"
// @Param status query string false "Filter by status (active/deleted)"
// @Success 200 {object} response.Response{data=[]model.AuditLog}
// @Router /audit-logs [get]
func (h *AuditLogHandler) ListAuditLogs(c *gin.Context) {
	var query model.ListAuditLogsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		h.logger.Warn("Invalid query parameters", zap.Error(err))
		response.Error(c, err)
		return
	}

	entries, total, err := h.service.ListAuditLogs(c.Request.Context(), &query)
	if err != nil {
		h.logger.Error("Failed to list audit logs", zap.Error(err))
		response.Error(c, err)
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(query.PageSize)))
	meta := &response.Meta{
		Page:       query.Page,
		PageSize:   query.PageSize,
		TotalPages: totalPages,
		TotalCount: total,
	}

	response.SuccessWithMeta(c, entries, meta)
}

func resolveActor(c *gin.Context) string {
	if subject, ok := middleware.GetAuthSubject(c); ok && subject != "" {
		return subject
	}
	return "system"
}

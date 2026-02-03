package handler

import (
	"enterprise-microservice-system/common/logger"
	"enterprise-microservice-system/common/middleware"
	"enterprise-microservice-system/common/response"
	"enterprise-microservice-system/services/repository-service/internal/model"
	"enterprise-microservice-system/services/repository-service/internal/service"
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RepositoryHandler handles HTTP requests for repositories
type RepositoryHandler struct {
	service service.RepositoryService
	logger  *logger.Logger
}

// NewRepositoryHandler creates a new repository handler
func NewRepositoryHandler(service service.RepositoryService, logger *logger.Logger) *RepositoryHandler {
	return &RepositoryHandler{
		service: service,
		logger:  logger,
	}
}

// CreateRepository handles repository creation
// @Summary Create a new repository
// @Tags repositories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param repository body model.CreateRepositoryRequest true "Repository data"
// @Success 201 {object} response.Response{data=model.Repository}
// @Failure 400 {object} response.Response
// @Failure 409 {object} response.Response
// @Router /repositories [post]
func (h *RepositoryHandler) CreateRepository(c *gin.Context) {
	var req model.CreateRepositoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid request body", zap.Error(err))
		response.Error(c, err)
		return
	}

	actor := resolveActor(c)
	repo, err := h.service.CreateRepository(c.Request.Context(), &req, actor)
	if err != nil {
		h.logger.Error("Failed to create repository", zap.Error(err))
		response.Error(c, err)
		return
	}

	h.logger.Info("Repository created successfully", zap.Uint("repository_id", repo.ID))
	response.Created(c, repo)
}

// GetRepository handles retrieving a repository by ID
// @Summary Get a repository by ID
// @Tags repositories
// @Produce json
// @Security BearerAuth
// @Param id path int true "Repository ID"
// @Success 200 {object} response.Response{data=model.Repository}
// @Failure 404 {object} response.Response
// @Router /repositories/{id} [get]
func (h *RepositoryHandler) GetRepository(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		h.logger.Warn("Invalid repository ID", zap.Error(err))
		response.Error(c, err)
		return
	}

	repo, err := h.service.GetRepository(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("Failed to get repository", zap.Uint64("repository_id", id), zap.Error(err))
		response.Error(c, err)
		return
	}

	response.Success(c, repo)
}

// UpdateRepository handles updating a repository
// @Summary Update a repository
// @Tags repositories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Repository ID"
// @Param repository body model.UpdateRepositoryRequest true "Repository update data"
// @Success 200 {object} response.Response{data=model.Repository}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /repositories/{id} [put]
func (h *RepositoryHandler) UpdateRepository(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		h.logger.Warn("Invalid repository ID", zap.Error(err))
		response.Error(c, err)
		return
	}

	var req model.UpdateRepositoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid request body", zap.Error(err))
		response.Error(c, err)
		return
	}

	actor := resolveActor(c)
	repo, err := h.service.UpdateRepository(c.Request.Context(), uint(id), &req, actor)
	if err != nil {
		h.logger.Error("Failed to update repository", zap.Uint64("repository_id", id), zap.Error(err))
		response.Error(c, err)
		return
	}

	h.logger.Info("Repository updated successfully", zap.Uint64("repository_id", id))
	response.Success(c, repo)
}

// DeleteRepository handles deleting a repository
// @Summary Delete a repository
// @Tags repositories
// @Produce json
// @Security BearerAuth
// @Param id path int true "Repository ID"
// @Success 200 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /repositories/{id} [delete]
func (h *RepositoryHandler) DeleteRepository(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		h.logger.Warn("Invalid repository ID", zap.Error(err))
		response.Error(c, err)
		return
	}

	actor := resolveActor(c)
	if err := h.service.DeleteRepository(c.Request.Context(), uint(id), actor); err != nil {
		h.logger.Error("Failed to delete repository", zap.Uint64("repository_id", id), zap.Error(err))
		response.Error(c, err)
		return
	}

	h.logger.Info("Repository deleted successfully", zap.Uint64("repository_id", id))
	response.Success(c, gin.H{"message": "repository deleted successfully"})
}

// ListRepositories handles listing repositories with pagination
// @Summary List repositories
// @Tags repositories
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Param search query string false "Search term"
// @Param owner_id query int false "Filter by owner ID"
// @Param visibility query string false "Filter by visibility (public/private)"
// @Param status query string false "Filter by status (active/inactive/deleted)"
// @Success 200 {object} response.Response{data=[]model.Repository}
// @Router /repositories [get]
func (h *RepositoryHandler) ListRepositories(c *gin.Context) {
	var query model.ListRepositoriesQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		h.logger.Warn("Invalid query parameters", zap.Error(err))
		response.Error(c, err)
		return
	}

	repos, total, err := h.service.ListRepositories(c.Request.Context(), &query)
	if err != nil {
		h.logger.Error("Failed to list repositories", zap.Error(err))
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

	response.SuccessWithMeta(c, repos, meta)
}

func resolveActor(c *gin.Context) string {
	if subject, ok := middleware.GetAuthSubject(c); ok && subject != "" {
		return subject
	}
	return "system"
}

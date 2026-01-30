package handler

import (
	"enterprise-microservice-system/common/logger"
	"enterprise-microservice-system/common/response"
	"enterprise-microservice-system/services/user-service/internal/model"
	"enterprise-microservice-system/services/user-service/internal/service"
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// UserHandler handles HTTP requests for users
type UserHandler struct {
	service service.UserService
	logger  *logger.Logger
}

// NewUserHandler creates a new user handler
func NewUserHandler(service service.UserService, logger *logger.Logger) *UserHandler {
	return &UserHandler{
		service: service,
		logger:  logger,
	}
}

// CreateUser handles user creation
// @Summary Create a new user
// @Tags users
// @Accept json
// @Produce json
// @Param user body model.CreateUserRequest true "User data"
// @Success 201 {object} response.Response{data=model.User}
// @Failure 400 {object} response.Response
// @Failure 409 {object} response.Response
// @Router /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req model.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid request body", zap.Error(err))
		response.Error(c, err)
		return
	}

	user, err := h.service.CreateUser(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to create user", zap.Error(err))
		response.Error(c, err)
		return
	}

	h.logger.Info("User created successfully", zap.Uint("user_id", user.ID))
	response.Created(c, user)
}

// GetUser handles retrieving a user by ID
// @Summary Get a user by ID
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} response.Response{data=model.User}
// @Failure 404 {object} response.Response
// @Router /users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		h.logger.Warn("Invalid user ID", zap.Error(err))
		response.Error(c, err)
		return
	}

	user, err := h.service.GetUser(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("Failed to get user", zap.Uint64("user_id", id), zap.Error(err))
		response.Error(c, err)
		return
	}

	response.Success(c, user)
}

// UpdateUser handles updating a user
// @Summary Update a user
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body model.UpdateUserRequest true "User update data"
// @Success 200 {object} response.Response{data=model.User}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		h.logger.Warn("Invalid user ID", zap.Error(err))
		response.Error(c, err)
		return
	}

	var req model.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid request body", zap.Error(err))
		response.Error(c, err)
		return
	}

	user, err := h.service.UpdateUser(c.Request.Context(), uint(id), &req)
	if err != nil {
		h.logger.Error("Failed to update user", zap.Uint64("user_id", id), zap.Error(err))
		response.Error(c, err)
		return
	}

	h.logger.Info("User updated successfully", zap.Uint64("user_id", id))
	response.Success(c, user)
}

// DeleteUser handles deleting a user
// @Summary Delete a user
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		h.logger.Warn("Invalid user ID", zap.Error(err))
		response.Error(c, err)
		return
	}

	if err := h.service.DeleteUser(c.Request.Context(), uint(id)); err != nil {
		h.logger.Error("Failed to delete user", zap.Uint64("user_id", id), zap.Error(err))
		response.Error(c, err)
		return
	}

	h.logger.Info("User deleted successfully", zap.Uint64("user_id", id))
	response.Success(c, gin.H{"message": "user deleted successfully"})
}

// ListUsers handles listing users with pagination
// @Summary List users
// @Tags users
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Param search query string false "Search term"
// @Param active query bool false "Filter by active status"
// @Success 200 {object} response.Response{data=[]model.User}
// @Router /users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	var query model.ListUsersQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		h.logger.Warn("Invalid query parameters", zap.Error(err))
		response.Error(c, err)
		return
	}

	users, total, err := h.service.ListUsers(c.Request.Context(), &query)
	if err != nil {
		h.logger.Error("Failed to list users", zap.Error(err))
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

	response.SuccessWithMeta(c, users, meta)
}

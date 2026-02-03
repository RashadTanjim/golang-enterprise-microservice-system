package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"enterprise-microservice-system/common/cache"
	"enterprise-microservice-system/common/errors"
	"enterprise-microservice-system/services/user-service/internal/model"
	"enterprise-microservice-system/services/user-service/internal/repository"

	"gorm.io/gorm"
)

// UserService defines the business logic interface for users
type UserService interface {
	CreateUser(ctx context.Context, req *model.CreateUserRequest, actor string) (*model.User, error)
	GetUser(ctx context.Context, id uint) (*model.User, error)
	UpdateUser(ctx context.Context, id uint, req *model.UpdateUserRequest, actor string) (*model.User, error)
	DeleteUser(ctx context.Context, id uint, actor string) error
	ListUsers(ctx context.Context, query *model.ListUsersQuery) ([]*model.User, int64, error)
}

// userService implements UserService
type userService struct {
	repo  repository.UserRepository
	cache *cache.Cache
}

// NewUserService creates a new user service
func NewUserService(repo repository.UserRepository, cacheClient *cache.Cache) UserService {
	return &userService{
		repo:  repo,
		cache: cacheClient,
	}
}

// CreateUser creates a new user
func (s *userService) CreateUser(ctx context.Context, req *model.CreateUserRequest, actor string) (*model.User, error) {
	// Check if email already exists
	existingUser, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, errors.NewInternal("failed to check email uniqueness", err)
	}
	if existingUser != nil {
		return nil, errors.NewConflict("email already exists")
	}

	// Create user
	user := &model.User{
		Email:  req.Email,
		Name:   req.Name,
		Age:    req.Age,
		Status: model.UserStatusActive,
	}

	if actor == "" {
		actor = "system"
	}
	user.CreatedBy = actor
	user.UpdatedBy = actor

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, errors.NewInternal("failed to create user", err)
	}

	s.cacheSetUser(ctx, user)

	return user, nil
}

// GetUser retrieves a user by ID
func (s *userService) GetUser(ctx context.Context, id uint) (*model.User, error) {
	if cached := s.cacheGetUser(ctx, id); cached != nil {
		return cached, nil
	}

	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewNotFound("user")
		}
		return nil, errors.NewInternal("failed to get user", err)
	}

	s.cacheSetUser(ctx, user)
	return user, nil
}

// UpdateUser updates a user
func (s *userService) UpdateUser(ctx context.Context, id uint, req *model.UpdateUserRequest, actor string) (*model.User, error) {
	// Get existing user
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewNotFound("user")
		}
		return nil, errors.NewInternal("failed to get user", err)
	}

	// Update fields if provided
	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.Age != nil {
		user.Age = *req.Age
	}
	if req.Status != nil {
		user.Status = *req.Status
	}
	if actor == "" {
		actor = "system"
	}
	user.UpdatedBy = actor

	// Save updates
	if err := s.repo.Update(ctx, user); err != nil {
		return nil, errors.NewInternal("failed to update user", err)
	}

	s.cacheSetUser(ctx, user)

	return user, nil
}

// DeleteUser deletes a user
func (s *userService) DeleteUser(ctx context.Context, id uint, actor string) error {
	// Check if user exists
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.NewNotFound("user")
		}
		return errors.NewInternal("failed to get user", err)
	}

	if actor == "" {
		actor = "system"
	}

	// Delete user (status-based soft delete)
	if err := s.repo.Delete(ctx, id, actor); err != nil {
		return errors.NewInternal("failed to delete user", err)
	}

	s.cacheDeleteUser(ctx, id)

	return nil
}

// ListUsers retrieves a paginated list of users
func (s *userService) ListUsers(ctx context.Context, query *model.ListUsersQuery) ([]*model.User, int64, error) {
	query.ApplyDefaults()

	if cachedUsers, cachedTotal, ok := s.cacheGetUserList(ctx, query); ok {
		return cachedUsers, cachedTotal, nil
	}

	users, total, err := s.repo.List(ctx, query)
	if err != nil {
		return nil, 0, errors.NewInternal("failed to list users", err)
	}

	s.cacheSetUserList(ctx, query, users, total)
	return users, total, nil
}

func (s *userService) cacheGetUser(ctx context.Context, id uint) *model.User {
	if s.cache == nil || !s.cache.Enabled() {
		return nil
	}

	var user model.User
	found, err := s.cache.GetJSON(ctx, fmt.Sprintf("user:%d", id), &user)
	if err != nil || !found {
		return nil
	}
	return &user
}

func (s *userService) cacheSetUser(ctx context.Context, user *model.User) {
	if s.cache == nil || !s.cache.Enabled() || user == nil {
		return
	}

	_ = s.cache.SetJSON(ctx, fmt.Sprintf("user:%d", user.ID), user, 0)
}

func (s *userService) cacheDeleteUser(ctx context.Context, id uint) {
	if s.cache == nil || !s.cache.Enabled() {
		return
	}
	_ = s.cache.Delete(ctx, fmt.Sprintf("user:%d", id))
}

func (s *userService) cacheGetUserList(ctx context.Context, query *model.ListUsersQuery) ([]*model.User, int64, bool) {
	if s.cache == nil || !s.cache.Enabled() {
		return nil, 0, false
	}

	key := s.userListCacheKey(query)
	var payload struct {
		Users []*model.User `json:"users"`
		Total int64         `json:"total"`
	}

	found, err := s.cache.GetJSON(ctx, key, &payload)
	if err != nil || !found {
		return nil, 0, false
	}
	return payload.Users, payload.Total, true
}

func (s *userService) cacheSetUserList(ctx context.Context, query *model.ListUsersQuery, users []*model.User, total int64) {
	if s.cache == nil || !s.cache.Enabled() {
		return
	}

	key := s.userListCacheKey(query)
	payload := struct {
		Users []*model.User `json:"users"`
		Total int64         `json:"total"`
	}{
		Users: users,
		Total: total,
	}

	_ = s.cache.SetJSON(ctx, key, payload, 60*time.Second)
}

func (s *userService) userListCacheKey(query *model.ListUsersQuery) string {
	status := "any"
	if query.Status != nil {
		status = *query.Status
	}

	search := strings.TrimSpace(query.Search)
	if search == "" {
		search = "all"
	}

	return fmt.Sprintf("users:list:p%d:ps%d:search:%s:status:%s", query.Page, query.PageSize, search, status)
}

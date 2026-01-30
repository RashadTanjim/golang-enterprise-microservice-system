package service

import (
	"context"
	"enterprise-microservice-system/common/errors"
	"enterprise-microservice-system/services/user-service/internal/model"
	"enterprise-microservice-system/services/user-service/internal/repository"

	"gorm.io/gorm"
)

// UserService defines the business logic interface for users
type UserService interface {
	CreateUser(ctx context.Context, req *model.CreateUserRequest) (*model.User, error)
	GetUser(ctx context.Context, id uint) (*model.User, error)
	UpdateUser(ctx context.Context, id uint, req *model.UpdateUserRequest) (*model.User, error)
	DeleteUser(ctx context.Context, id uint) error
	ListUsers(ctx context.Context, query *model.ListUsersQuery) ([]*model.User, int64, error)
}

// userService implements UserService
type userService struct {
	repo repository.UserRepository
}

// NewUserService creates a new user service
func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

// CreateUser creates a new user
func (s *userService) CreateUser(ctx context.Context, req *model.CreateUserRequest) (*model.User, error) {
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
		Active: true,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, errors.NewInternal("failed to create user", err)
	}

	return user, nil
}

// GetUser retrieves a user by ID
func (s *userService) GetUser(ctx context.Context, id uint) (*model.User, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewNotFound("user")
		}
		return nil, errors.NewInternal("failed to get user", err)
	}
	return user, nil
}

// UpdateUser updates a user
func (s *userService) UpdateUser(ctx context.Context, id uint, req *model.UpdateUserRequest) (*model.User, error) {
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
	if req.Active != nil {
		user.Active = *req.Active
	}

	// Save updates
	if err := s.repo.Update(ctx, user); err != nil {
		return nil, errors.NewInternal("failed to update user", err)
	}

	return user, nil
}

// DeleteUser deletes a user
func (s *userService) DeleteUser(ctx context.Context, id uint) error {
	// Check if user exists
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.NewNotFound("user")
		}
		return errors.NewInternal("failed to get user", err)
	}

	// Delete user
	if err := s.repo.Delete(ctx, id); err != nil {
		return errors.NewInternal("failed to delete user", err)
	}

	return nil
}

// ListUsers retrieves a paginated list of users
func (s *userService) ListUsers(ctx context.Context, query *model.ListUsersQuery) ([]*model.User, int64, error) {
	query.ApplyDefaults()

	users, total, err := s.repo.List(ctx, query)
	if err != nil {
		return nil, 0, errors.NewInternal("failed to list users", err)
	}

	return users, total, nil
}

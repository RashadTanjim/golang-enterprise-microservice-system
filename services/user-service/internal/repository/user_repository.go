package repository

import (
	"context"
	"time"

	"enterprise-microservice-system/services/user-service/internal/model"

	"gorm.io/gorm"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	FindByID(ctx context.Context, id uint) (*model.User, error)
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id uint, updatedBy string) error
	List(ctx context.Context, query *model.ListUsersQuery) ([]*model.User, int64, error)
}

// userRepository implements UserRepository
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// Create creates a new user
func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// FindByID finds a user by ID
func (r *userRepository) FindByID(ctx context.Context, id uint) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).
		Where("status <> ?", model.UserStatusDeleted).
		First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByEmail finds a user by email
func (r *userRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).
		Where("email = ? AND status <> ?", email, model.UserStatusDeleted).
		First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update updates a user
func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

// Delete soft deletes a user
func (r *userRepository) Delete(ctx context.Context, id uint, updatedBy string) error {
	if updatedBy == "" {
		updatedBy = "system"
	}

	return r.db.WithContext(ctx).
		Model(&model.User{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     model.UserStatusDeleted,
			"updated_by": updatedBy,
			"updated_at": time.Now().UTC(),
		}).Error
}

// List retrieves a paginated list of users
func (r *userRepository) List(ctx context.Context, query *model.ListUsersQuery) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	db := r.db.WithContext(ctx).Model(&model.User{})

	// Apply filters
	if query.Search != "" {
		searchPattern := "%" + query.Search + "%"
		db = db.Where("name ILIKE ? OR email ILIKE ?", searchPattern, searchPattern)
	}

	if query.Status != nil {
		db = db.Where("status = ?", *query.Status)
	} else {
		db = db.Where("status <> ?", model.UserStatusDeleted)
	}

	// Count total
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	err := db.Offset(query.Offset()).
		Limit(query.PageSize).
		Order("created_at DESC").
		Find(&users).Error

	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

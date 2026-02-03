package repository

import (
	"context"
	"enterprise-microservice-system/services/repository-service/internal/model"

	"gorm.io/gorm"
)

// RepositoryRepository defines the interface for repository data operations
type RepositoryRepository interface {
	Create(ctx context.Context, repo *model.Repository) error
	FindByID(ctx context.Context, id uint) (*model.Repository, error)
	FindByName(ctx context.Context, name string) (*model.Repository, error)
	Update(ctx context.Context, repo *model.Repository) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, query *model.ListRepositoriesQuery) ([]*model.Repository, int64, error)
}

// repositoryRepository implements RepositoryRepository
type repositoryRepository struct {
	db *gorm.DB
}

// NewRepositoryRepository creates a new repository repository
func NewRepositoryRepository(db *gorm.DB) RepositoryRepository {
	return &repositoryRepository{db: db}
}

// Create creates a new repository
func (r *repositoryRepository) Create(ctx context.Context, repo *model.Repository) error {
	return r.db.WithContext(ctx).Create(repo).Error
}

// FindByID finds a repository by ID
func (r *repositoryRepository) FindByID(ctx context.Context, id uint) (*model.Repository, error) {
	var repo model.Repository
	err := r.db.WithContext(ctx).First(&repo, id).Error
	if err != nil {
		return nil, err
	}
	return &repo, nil
}

// FindByName finds a repository by name
func (r *repositoryRepository) FindByName(ctx context.Context, name string) (*model.Repository, error) {
	var repo model.Repository
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&repo).Error
	if err != nil {
		return nil, err
	}
	return &repo, nil
}

// Update updates a repository
func (r *repositoryRepository) Update(ctx context.Context, repo *model.Repository) error {
	return r.db.WithContext(ctx).Save(repo).Error
}

// Delete soft deletes a repository
func (r *repositoryRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Repository{}, id).Error
}

// List retrieves a paginated list of repositories
func (r *repositoryRepository) List(ctx context.Context, query *model.ListRepositoriesQuery) ([]*model.Repository, int64, error) {
	var repos []*model.Repository
	var total int64

	db := r.db.WithContext(ctx).Model(&model.Repository{})

	// Apply filters
	if query.Search != "" {
		searchPattern := "%" + query.Search + "%"
		db = db.Where("name ILIKE ? OR description ILIKE ?", searchPattern, searchPattern)
	}

	if query.OwnerID != nil {
		db = db.Where("owner_id = ?", *query.OwnerID)
	}

	if query.Visibility != "" {
		db = db.Where("visibility = ?", query.Visibility)
	}

	if query.Active != nil {
		db = db.Where("active = ?", *query.Active)
	}

	// Count total
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	err := db.Offset(query.Offset()).
		Limit(query.PageSize).
		Order("created_at DESC").
		Find(&repos).Error

	if err != nil {
		return nil, 0, err
	}

	return repos, total, nil
}

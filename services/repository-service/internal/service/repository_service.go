package service

import (
	"context"
	"enterprise-microservice-system/common/errors"
	"enterprise-microservice-system/services/repository-service/internal/model"
	"enterprise-microservice-system/services/repository-service/internal/repository"

	"gorm.io/gorm"
)

// RepositoryService defines the business logic interface for repositories
type RepositoryService interface {
	CreateRepository(ctx context.Context, req *model.CreateRepositoryRequest) (*model.Repository, error)
	GetRepository(ctx context.Context, id uint) (*model.Repository, error)
	UpdateRepository(ctx context.Context, id uint, req *model.UpdateRepositoryRequest) (*model.Repository, error)
	DeleteRepository(ctx context.Context, id uint) error
	ListRepositories(ctx context.Context, query *model.ListRepositoriesQuery) ([]*model.Repository, int64, error)
}

// repositoryService implements RepositoryService
type repositoryService struct {
	repo repository.RepositoryRepository
}

// NewRepositoryService creates a new repository service
func NewRepositoryService(repo repository.RepositoryRepository) RepositoryService {
	return &repositoryService{repo: repo}
}

// CreateRepository creates a new repository
func (s *repositoryService) CreateRepository(ctx context.Context, req *model.CreateRepositoryRequest) (*model.Repository, error) {
	// Check if repository name already exists
	existingRepo, err := s.repo.FindByName(ctx, req.Name)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, errors.NewInternal("failed to check repository name uniqueness", err)
	}
	if existingRepo != nil {
		return nil, errors.NewConflict("repository name already exists")
	}

	// Set default visibility if not provided
	visibility := req.Visibility
	if visibility == "" {
		visibility = "private"
	}

	// Create repository
	repo := &model.Repository{
		Name:        req.Name,
		Description: req.Description,
		OwnerID:     req.OwnerID,
		Visibility:  visibility,
		URL:         req.URL,
		Active:      true,
	}

	if err := s.repo.Create(ctx, repo); err != nil {
		return nil, errors.NewInternal("failed to create repository", err)
	}

	return repo, nil
}

// GetRepository retrieves a repository by ID
func (s *repositoryService) GetRepository(ctx context.Context, id uint) (*model.Repository, error) {
	repo, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewNotFound("repository")
		}
		return nil, errors.NewInternal("failed to get repository", err)
	}
	return repo, nil
}

// UpdateRepository updates a repository
func (s *repositoryService) UpdateRepository(ctx context.Context, id uint, req *model.UpdateRepositoryRequest) (*model.Repository, error) {
	// Get existing repository
	repo, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewNotFound("repository")
		}
		return nil, errors.NewInternal("failed to get repository", err)
	}

	// Check if name is being updated and if it conflicts
	if req.Name != nil && *req.Name != repo.Name {
		existingRepo, err := s.repo.FindByName(ctx, *req.Name)
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, errors.NewInternal("failed to check repository name uniqueness", err)
		}
		if existingRepo != nil {
			return nil, errors.NewConflict("repository name already exists")
		}
		repo.Name = *req.Name
	}

	// Update fields if provided
	if req.Description != nil {
		repo.Description = *req.Description
	}
	if req.Visibility != nil {
		repo.Visibility = *req.Visibility
	}
	if req.URL != nil {
		repo.URL = *req.URL
	}
	if req.Active != nil {
		repo.Active = *req.Active
	}

	// Save updates
	if err := s.repo.Update(ctx, repo); err != nil {
		return nil, errors.NewInternal("failed to update repository", err)
	}

	return repo, nil
}

// DeleteRepository deletes a repository
func (s *repositoryService) DeleteRepository(ctx context.Context, id uint) error {
	// Check if repository exists
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.NewNotFound("repository")
		}
		return errors.NewInternal("failed to get repository", err)
	}

	// Delete repository
	if err := s.repo.Delete(ctx, id); err != nil {
		return errors.NewInternal("failed to delete repository", err)
	}

	return nil
}

// ListRepositories retrieves a paginated list of repositories
func (s *repositoryService) ListRepositories(ctx context.Context, query *model.ListRepositoriesQuery) ([]*model.Repository, int64, error) {
	query.ApplyDefaults()

	repos, total, err := s.repo.List(ctx, query)
	if err != nil {
		return nil, 0, errors.NewInternal("failed to list repositories", err)
	}

	return repos, total, nil
}

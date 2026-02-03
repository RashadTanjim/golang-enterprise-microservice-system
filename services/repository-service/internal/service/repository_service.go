package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"enterprise-microservice-system/common/cache"
	"enterprise-microservice-system/common/errors"
	"enterprise-microservice-system/services/repository-service/internal/model"
	"enterprise-microservice-system/services/repository-service/internal/repository"

	"gorm.io/gorm"
)

// RepositoryService defines the business logic interface for repositories
type RepositoryService interface {
	CreateRepository(ctx context.Context, req *model.CreateRepositoryRequest, actor string) (*model.Repository, error)
	GetRepository(ctx context.Context, id uint) (*model.Repository, error)
	UpdateRepository(ctx context.Context, id uint, req *model.UpdateRepositoryRequest, actor string) (*model.Repository, error)
	DeleteRepository(ctx context.Context, id uint, actor string) error
	ListRepositories(ctx context.Context, query *model.ListRepositoriesQuery) ([]*model.Repository, int64, error)
}

// repositoryService implements RepositoryService
type repositoryService struct {
	repo  repository.RepositoryRepository
	cache *cache.Cache
}

// NewRepositoryService creates a new repository service
func NewRepositoryService(repo repository.RepositoryRepository, cacheClient *cache.Cache) RepositoryService {
	return &repositoryService{
		repo:  repo,
		cache: cacheClient,
	}
}

// CreateRepository creates a new repository
func (s *repositoryService) CreateRepository(ctx context.Context, req *model.CreateRepositoryRequest, actor string) (*model.Repository, error) {
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
		Status:      model.RepositoryStatusActive,
	}

	if actor == "" {
		actor = "system"
	}
	repo.CreatedBy = actor
	repo.UpdatedBy = actor

	if err := s.repo.Create(ctx, repo); err != nil {
		return nil, errors.NewInternal("failed to create repository", err)
	}

	s.cacheSetRepository(ctx, repo)

	return repo, nil
}

// GetRepository retrieves a repository by ID
func (s *repositoryService) GetRepository(ctx context.Context, id uint) (*model.Repository, error) {
	if cached := s.cacheGetRepository(ctx, id); cached != nil {
		return cached, nil
	}

	repo, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewNotFound("repository")
		}
		return nil, errors.NewInternal("failed to get repository", err)
	}
	s.cacheSetRepository(ctx, repo)
	return repo, nil
}

// UpdateRepository updates a repository
func (s *repositoryService) UpdateRepository(ctx context.Context, id uint, req *model.UpdateRepositoryRequest, actor string) (*model.Repository, error) {
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
	if req.Status != nil {
		repo.Status = *req.Status
	}
	if actor == "" {
		actor = "system"
	}
	repo.UpdatedBy = actor

	// Save updates
	if err := s.repo.Update(ctx, repo); err != nil {
		return nil, errors.NewInternal("failed to update repository", err)
	}

	s.cacheSetRepository(ctx, repo)

	return repo, nil
}

// DeleteRepository deletes a repository
func (s *repositoryService) DeleteRepository(ctx context.Context, id uint, actor string) error {
	// Check if repository exists
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.NewNotFound("repository")
		}
		return errors.NewInternal("failed to get repository", err)
	}

	if actor == "" {
		actor = "system"
	}

	// Delete repository (status-based soft delete)
	if err := s.repo.Delete(ctx, id, actor); err != nil {
		return errors.NewInternal("failed to delete repository", err)
	}

	s.cacheDeleteRepository(ctx, id)

	return nil
}

// ListRepositories retrieves a paginated list of repositories
func (s *repositoryService) ListRepositories(ctx context.Context, query *model.ListRepositoriesQuery) ([]*model.Repository, int64, error) {
	query.ApplyDefaults()

	if cachedRepos, cachedTotal, ok := s.cacheGetRepositoryList(ctx, query); ok {
		return cachedRepos, cachedTotal, nil
	}

	repos, total, err := s.repo.List(ctx, query)
	if err != nil {
		return nil, 0, errors.NewInternal("failed to list repositories", err)
	}

	s.cacheSetRepositoryList(ctx, query, repos, total)
	return repos, total, nil
}

func (s *repositoryService) cacheGetRepository(ctx context.Context, id uint) *model.Repository {
	if s.cache == nil || !s.cache.Enabled() {
		return nil
	}

	var repo model.Repository
	found, err := s.cache.GetJSON(ctx, fmt.Sprintf("repository:%d", id), &repo)
	if err != nil || !found {
		return nil
	}
	return &repo
}

func (s *repositoryService) cacheSetRepository(ctx context.Context, repo *model.Repository) {
	if s.cache == nil || !s.cache.Enabled() || repo == nil {
		return
	}

	_ = s.cache.SetJSON(ctx, fmt.Sprintf("repository:%d", repo.ID), repo, 0)
}

func (s *repositoryService) cacheDeleteRepository(ctx context.Context, id uint) {
	if s.cache == nil || !s.cache.Enabled() {
		return
	}
	_ = s.cache.Delete(ctx, fmt.Sprintf("repository:%d", id))
}

func (s *repositoryService) cacheGetRepositoryList(ctx context.Context, query *model.ListRepositoriesQuery) ([]*model.Repository, int64, bool) {
	if s.cache == nil || !s.cache.Enabled() {
		return nil, 0, false
	}

	key := s.repositoryListCacheKey(query)
	var payload struct {
		Repositories []*model.Repository `json:"repositories"`
		Total        int64               `json:"total"`
	}

	found, err := s.cache.GetJSON(ctx, key, &payload)
	if err != nil || !found {
		return nil, 0, false
	}
	return payload.Repositories, payload.Total, true
}

func (s *repositoryService) cacheSetRepositoryList(ctx context.Context, query *model.ListRepositoriesQuery, repos []*model.Repository, total int64) {
	if s.cache == nil || !s.cache.Enabled() {
		return
	}

	key := s.repositoryListCacheKey(query)
	payload := struct {
		Repositories []*model.Repository `json:"repositories"`
		Total        int64               `json:"total"`
	}{
		Repositories: repos,
		Total:        total,
	}

	_ = s.cache.SetJSON(ctx, key, payload, 60*time.Second)
}

func (s *repositoryService) repositoryListCacheKey(query *model.ListRepositoriesQuery) string {
	search := strings.TrimSpace(query.Search)
	if search == "" {
		search = "all"
	}

	visibility := query.Visibility
	if visibility == "" {
		visibility = "all"
	}

	status := "any"
	if query.Status != nil {
		status = *query.Status
	}

	owner := "any"
	if query.OwnerID != nil {
		owner = fmt.Sprintf("%d", *query.OwnerID)
	}

	return fmt.Sprintf(
		"repositories:list:p%d:ps%d:search:%s:owner:%s:visibility:%s:status:%s",
		query.Page,
		query.PageSize,
		search,
		owner,
		visibility,
		status,
	)
}

package tests

import (
	"context"
	"enterprise-microservice-system/services/repository-service/internal/model"
	"enterprise-microservice-system/services/repository-service/internal/service"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockRepositoryRepository is a mock implementation of RepositoryRepository
type MockRepositoryRepository struct {
	mock.Mock
}

func (m *MockRepositoryRepository) Create(ctx context.Context, repo *model.Repository) error {
	args := m.Called(ctx, repo)
	return args.Error(0)
}

func (m *MockRepositoryRepository) FindByID(ctx context.Context, id uint) (*model.Repository, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Repository), args.Error(1)
}

func (m *MockRepositoryRepository) FindByName(ctx context.Context, name string) (*model.Repository, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Repository), args.Error(1)
}

func (m *MockRepositoryRepository) Update(ctx context.Context, repo *model.Repository) error {
	args := m.Called(ctx, repo)
	return args.Error(0)
}

func (m *MockRepositoryRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRepositoryRepository) List(ctx context.Context, query *model.ListRepositoriesQuery) ([]*model.Repository, int64, error) {
	args := m.Called(ctx, query)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*model.Repository), args.Get(1).(int64), args.Error(2)
}

func TestCreateRepository_Success(t *testing.T) {
	mockRepo := new(MockRepositoryRepository)
	service := service.NewRepositoryService(mockRepo)

	req := &model.CreateRepositoryRequest{
		Name:        "test-repo",
		Description: "Test repository description",
		OwnerID:     1,
		Visibility:  "public",
		URL:         "https://github.com/test/test-repo",
	}

	mockRepo.On("FindByName", mock.Anything, req.Name).Return(nil, gorm.ErrRecordNotFound)
	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.Repository")).Return(nil)

	repo, err := service.CreateRepository(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, repo)
	assert.Equal(t, req.Name, repo.Name)
	assert.Equal(t, req.Description, repo.Description)
	assert.Equal(t, req.OwnerID, repo.OwnerID)
	assert.Equal(t, req.Visibility, repo.Visibility)
	mockRepo.AssertExpectations(t)
}

func TestCreateRepository_DuplicateName(t *testing.T) {
	mockRepo := new(MockRepositoryRepository)
	service := service.NewRepositoryService(mockRepo)

	req := &model.CreateRepositoryRequest{
		Name:        "test-repo",
		Description: "Test repository description",
		OwnerID:     1,
	}

	existingRepo := &model.Repository{
		ID:   1,
		Name: req.Name,
	}

	mockRepo.On("FindByName", mock.Anything, req.Name).Return(existingRepo, nil)

	repo, err := service.CreateRepository(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, repo)
	mockRepo.AssertExpectations(t)
}

func TestUpdateRepository_Success(t *testing.T) {
	mockRepo := new(MockRepositoryRepository)
	service := service.NewRepositoryService(mockRepo)

	existingRepo := &model.Repository{
		ID:          1,
		Name:        "test-repo",
		Description: "Old description",
		OwnerID:     1,
	}

	newDesc := "New description"
	req := &model.UpdateRepositoryRequest{
		Description: &newDesc,
	}

	mockRepo.On("FindByID", mock.Anything, uint(1)).Return(existingRepo, nil)
	mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*model.Repository")).Return(nil)

	repo, err := service.UpdateRepository(context.Background(), 1, req)

	assert.NoError(t, err)
	assert.NotNil(t, repo)
	assert.Equal(t, newDesc, repo.Description)
	mockRepo.AssertExpectations(t)
}

func TestGetRepository_Success(t *testing.T) {
	mockRepo := new(MockRepositoryRepository)
	service := service.NewRepositoryService(mockRepo)

	expectedRepo := &model.Repository{
		ID:          1,
		Name:        "test-repo",
		Description: "Test repository",
		OwnerID:     1,
	}

	mockRepo.On("FindByID", mock.Anything, uint(1)).Return(expectedRepo, nil)

	repo, err := service.GetRepository(context.Background(), 1)

	assert.NoError(t, err)
	assert.NotNil(t, repo)
	assert.Equal(t, expectedRepo.ID, repo.ID)
	assert.Equal(t, expectedRepo.Name, repo.Name)
	mockRepo.AssertExpectations(t)
}

func TestGetRepository_NotFound(t *testing.T) {
	mockRepo := new(MockRepositoryRepository)
	service := service.NewRepositoryService(mockRepo)

	mockRepo.On("FindByID", mock.Anything, uint(999)).Return(nil, gorm.ErrRecordNotFound)

	repo, err := service.GetRepository(context.Background(), 999)

	assert.Error(t, err)
	assert.Nil(t, repo)
	mockRepo.AssertExpectations(t)
}

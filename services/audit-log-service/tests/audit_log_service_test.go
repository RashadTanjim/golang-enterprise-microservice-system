package tests

import (
	"context"
	"testing"

	"enterprise-microservice-system/common/errors"
	"enterprise-microservice-system/services/audit-log-service/internal/model"
	"enterprise-microservice-system/services/audit-log-service/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockAuditLogRepository is a mock implementation of AuditLogRepository
type MockAuditLogRepository struct {
	mock.Mock
}

func (m *MockAuditLogRepository) Create(ctx context.Context, entry *model.AuditLog) error {
	args := m.Called(ctx, entry)
	return args.Error(0)
}

func (m *MockAuditLogRepository) FindByID(ctx context.Context, id uint) (*model.AuditLog, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.AuditLog), args.Error(1)
}

func (m *MockAuditLogRepository) Update(ctx context.Context, entry *model.AuditLog) error {
	args := m.Called(ctx, entry)
	return args.Error(0)
}

func (m *MockAuditLogRepository) Delete(ctx context.Context, id uint, updatedBy string) error {
	args := m.Called(ctx, id, updatedBy)
	return args.Error(0)
}

func (m *MockAuditLogRepository) List(ctx context.Context, query *model.ListAuditLogsQuery) ([]*model.AuditLog, int64, error) {
	args := m.Called(ctx, query)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*model.AuditLog), args.Get(1).(int64), args.Error(2)
}

func TestCreateAuditLog_Success(t *testing.T) {
	mockRepo := new(MockAuditLogRepository)
	svc := service.NewAuditLogService(mockRepo, nil)

	req := &model.CreateAuditLogRequest{
		Action:       "user.login",
		ResourceType: "user",
		ResourceID:   "1",
		Description:  "User signed in",
		Metadata:     "{\"ip\":\"127.0.0.1\"}",
	}

	mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(entry *model.AuditLog) bool {
		return entry.Action == req.Action &&
			entry.ResourceType == req.ResourceType &&
			entry.ResourceID == req.ResourceID &&
			entry.Actor == "tester" &&
			entry.CreatedBy == "tester" &&
			entry.UpdatedBy == "tester" &&
			entry.Status == model.AuditLogStatusActive
	})).Return(nil)

	entry, err := svc.CreateAuditLog(context.Background(), req, "tester")

	assert.NoError(t, err)
	assert.NotNil(t, entry)
	assert.Equal(t, req.Action, entry.Action)
	assert.Equal(t, "tester", entry.Actor)
	mockRepo.AssertExpectations(t)
}

func TestGetAuditLog_NotFound(t *testing.T) {
	mockRepo := new(MockAuditLogRepository)
	svc := service.NewAuditLogService(mockRepo, nil)

	mockRepo.On("FindByID", mock.Anything, uint(42)).Return(nil, gorm.ErrRecordNotFound)

	entry, err := svc.GetAuditLog(context.Background(), 42)

	assert.Error(t, err)
	assert.Nil(t, entry)
	var appErr *errors.AppError
	assert.ErrorAs(t, err, &appErr)
	assert.Equal(t, errors.ErrCodeNotFound, appErr.Code)
	mockRepo.AssertExpectations(t)
}

func TestUpdateAuditLog_Success(t *testing.T) {
	mockRepo := new(MockAuditLogRepository)
	svc := service.NewAuditLogService(mockRepo, nil)

	existing := &model.AuditLog{
		ID:           10,
		Actor:        "system",
		Action:       "order.created",
		ResourceType: "order",
		ResourceID:   "42",
		Description:  "Order created",
		Status:       model.AuditLogStatusActive,
	}

	newDesc := "Order created via API"
	newStatus := model.AuditLogStatusActive
	req := &model.UpdateAuditLogRequest{
		Description: &newDesc,
		Status:      &newStatus,
	}

	mockRepo.On("FindByID", mock.Anything, uint(10)).Return(existing, nil)
	mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(entry *model.AuditLog) bool {
		return entry.ID == 10 && entry.Description == newDesc && entry.UpdatedBy == "admin"
	})).Return(nil)

	entry, err := svc.UpdateAuditLog(context.Background(), 10, req, "admin")

	assert.NoError(t, err)
	assert.NotNil(t, entry)
	assert.Equal(t, newDesc, entry.Description)
	assert.Equal(t, "admin", entry.UpdatedBy)
	mockRepo.AssertExpectations(t)
}

func TestDeleteAuditLog_Success(t *testing.T) {
	mockRepo := new(MockAuditLogRepository)
	svc := service.NewAuditLogService(mockRepo, nil)

	existing := &model.AuditLog{
		ID:     7,
		Action: "user.logout",
		Status: model.AuditLogStatusActive,
	}

	mockRepo.On("FindByID", mock.Anything, uint(7)).Return(existing, nil)
	mockRepo.On("Delete", mock.Anything, uint(7), "admin").Return(nil)

	err := svc.DeleteAuditLog(context.Background(), 7, "admin")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestListAuditLogs_Success(t *testing.T) {
	mockRepo := new(MockAuditLogRepository)
	svc := service.NewAuditLogService(mockRepo, nil)

	query := &model.ListAuditLogsQuery{
		Page:     1,
		PageSize: 2,
	}

	entries := []*model.AuditLog{
		{ID: 1, Action: "user.login"},
		{ID: 2, Action: "order.created"},
	}

	mockRepo.On("List", mock.Anything, query).Return(entries, int64(2), nil)

	result, total, err := svc.ListAuditLogs(context.Background(), query)

	assert.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, result, 2)
	mockRepo.AssertExpectations(t)
}

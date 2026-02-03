package tests

import (
	"context"
	"testing"

	"enterprise-microservice-system/services/order-service/internal/client"
	"enterprise-microservice-system/services/order-service/internal/model"
	"enterprise-microservice-system/services/order-service/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) Create(ctx context.Context, order *model.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func (m *MockOrderRepository) FindByID(ctx context.Context, id uint) (*model.Order, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Order), args.Error(1)
}

func (m *MockOrderRepository) Update(ctx context.Context, order *model.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func (m *MockOrderRepository) Delete(ctx context.Context, id uint, updatedBy string) error {
	args := m.Called(ctx, id, updatedBy)
	return args.Error(0)
}

func (m *MockOrderRepository) List(ctx context.Context, query *model.ListOrdersQuery) ([]*model.Order, int64, error) {
	args := m.Called(ctx, query)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*model.Order), args.Get(1).(int64), args.Error(2)
}

func (m *MockOrderRepository) FindByUserID(ctx context.Context, userID uint) ([]*model.Order, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Order), args.Error(1)
}

type MockUserClient struct {
	mock.Mock
}

func (m *MockUserClient) GetUser(ctx context.Context, userID uint) (*model.User, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserClient) GetCircuitBreakerState() float64 {
	args := m.Called()
	return args.Get(0).(float64)
}

var _ client.UserServiceClient = (*MockUserClient)(nil)

func TestCreateOrder_SetsAuditAndStatus(t *testing.T) {
	repo := new(MockOrderRepository)
	userClient := new(MockUserClient)
	svc := service.NewOrderService(repo, userClient, nil)

	req := &model.CreateOrderRequest{
		UserID:     1,
		ProductID:  "PROD-1",
		Quantity:   2,
		TotalPrice: 25.5,
	}

	userClient.On("GetUser", mock.Anything, uint(1)).Return(&model.User{ID: 1, Status: "active"}, nil)
	repo.On("Create", mock.Anything, mock.MatchedBy(func(order *model.Order) bool {
		return order.OrderStatus == model.OrderStatusPending &&
			order.Status == model.OrderRecordStatusActive &&
			order.CreatedBy == "tester" &&
			order.UpdatedBy == "tester"
	})).Return(nil)

	result, err := svc.CreateOrder(context.Background(), req, "tester")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, model.OrderStatusPending, result.OrderStatus)
	assert.Equal(t, model.OrderRecordStatusActive, result.Status)
	repo.AssertExpectations(t)
	userClient.AssertExpectations(t)
}

func TestUpdateOrder_UpdatesStatus(t *testing.T) {
	repo := new(MockOrderRepository)
	userClient := new(MockUserClient)
	svc := service.NewOrderService(repo, userClient, nil)

	existing := &model.Order{ID: 10, OrderStatus: model.OrderStatusPending, Status: model.OrderRecordStatusActive}
	newStatus := model.OrderStatusConfirmed

	repo.On("FindByID", mock.Anything, uint(10)).Return(existing, nil)
	repo.On("Update", mock.Anything, mock.MatchedBy(func(order *model.Order) bool {
		return order.OrderStatus == newStatus && order.UpdatedBy == "tester"
	})).Return(nil)

	result, err := svc.UpdateOrder(context.Background(), 10, &model.UpdateOrderRequest{OrderStatus: &newStatus}, "tester")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, newStatus, result.OrderStatus)
	repo.AssertExpectations(t)
}

func TestDeleteOrder_UsesActor(t *testing.T) {
	repo := new(MockOrderRepository)
	userClient := new(MockUserClient)
	svc := service.NewOrderService(repo, userClient, nil)

	repo.On("FindByID", mock.Anything, uint(7)).Return(&model.Order{ID: 7}, nil)
	repo.On("Delete", mock.Anything, uint(7), "tester").Return(nil)

	err := svc.DeleteOrder(context.Background(), 7, "tester")

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

package api

import (
	"subscriptions/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type mockStorage struct {
	mock mock.Mock
}

func NewMockStorage() *mockStorage {
	return &mockStorage{}
}

func (m *mockStorage) GetSubscription(userID uuid.UUID, serviceName string) (*models.Subscription, error) {
	args := m.mock.Called(userID, serviceName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Subscription), args.Error(1)
}

func (m *mockStorage) CreateSubscription(subscription *models.Subscription) (*models.Subscription, error) {
	args := m.mock.Called(subscription)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Subscription), args.Error(1)
}

func (m *mockStorage) UpdateSubscription(subscription *models.Subscription) (*models.Subscription, error) {
	args := m.mock.Called(subscription)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Subscription), args.Error(1)
}

func (m *mockStorage) DeleteSubscription(userID uuid.UUID, serviceName string) error {
	args := m.mock.Called(userID, serviceName)
	return args.Error(0)
}

func (m *mockStorage) ListSubscriptions(userID uuid.UUID, page int) ([]*models.Subscription, error) {
	args := m.mock.Called(userID, page)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Subscription), args.Error(1)
}

func (m *mockStorage) SubscriptionTotalCost(userID uuid.UUID, serviceName, startDate, endDate string) (*models.TotalSubscriptionCost, error) {
	args := m.mock.Called(userID, serviceName, startDate, endDate)
	return &models.TotalSubscriptionCost{TotalCost: args.Int(0)}, args.Error(1)
}

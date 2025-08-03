package storage

import (
	"subscriptions/internal/models"

	"github.com/google/uuid"
)

type Storage interface {
	GetSubscription(userID uuid.UUID, serviceName string) (*models.Subscription, error)
	CreateSubscription(subscription *models.Subscription) (*models.Subscription, error)
	UpdateSubscription(subscription *models.Subscription) (*models.Subscription, error)
	DeleteSubscription(userID uuid.UUID, serviceName string) error
	ListSubscriptions(userID uuid.UUID, page int) ([]*models.Subscription, error)
	SubscriptionTotalCost(userID uuid.UUID, serviceName, startDate, endDate string) (*models.TotalSubscriptionCost, error)
}

package models

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	UserID      uuid.UUID `json:"id" validate:"required, uuid"`
	ServiceName string    `json:"service_name" validate:"required"`
	Price       int32     `json:"price" validate:"required"`
	StartDate   time.Time `json:"start_date" validate:"required"`
	EndDate     time.Time `json:"end_date"`
}

type TotalSubscriptionCost struct {
	TotalCost int `json:"total_cost"`
}

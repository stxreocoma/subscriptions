package database

import (
	"subscriptions/internal/models"
	"subscriptions/internal/storage/database"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock/v4"
)

func TestGetSubscription(t *testing.T) {
	// Mock the database connection
	db, mock, err := database.NewMock()
	if err != nil {
		t.Fatalf("failed to create mock database: %v", err)
	}
	defer db.Close()

	// Define the expected behavior
	userID := uuid.New()
	serviceName := "Yandex Plus"
	mock.ExpectQuery("SELECT").WithArgs(userID, serviceName).WillReturnRows(pgxmock.NewRows([]string{"service_name", "user_id", "start_date", "end_date"}).
		AddRow(serviceName, userID, "2025-07", "2026-07"))

	// Call the method under test
	subscription, err := db.GetSubscription(userID, serviceName)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Validate the result
	if subscription == nil || subscription.UserID != userID {
		t.Errorf("expected subscription for user %s, got %v", userID, subscription)
	}

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCreateSubscription(t *testing.T) {
	db, mock, err := database.NewMock()
	if err != nil {
		t.Fatalf("failed to create mock database: %v", err)
	}
	defer db.Close()

	userID := uuid.New()
	subscription := &models.Subscription{
		UserID:      userID,
		ServiceName: "VK Music",
		Price:       299,
		StartDate:   time.Now(),
		EndDate:     time.Now().AddDate(0, 1, 0),
	}

	mock.ExpectExec("INSERT INTO subscriptions").
		WithArgs(subscription.UserID, subscription.ServiceName, subscription.Price, subscription.StartDate.Format("2006-01"), subscription.EndDate.Format("2006-01")).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	createdSubscription, err := db.CreateSubscription(subscription)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if createdSubscription == nil || createdSubscription.UserID != userID {
		t.Errorf("expected subscription for user %s, got %v", userID, createdSubscription)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestUpdateSubscription(t *testing.T) {
	db, mock, err := database.NewMock()
	if err != nil {
		t.Fatalf("failed to create mock database: %v", err)
	}
	defer db.Close()

	subscription := &models.Subscription{
		UserID:      uuid.New(),
		ServiceName: "Spotify",
		Price:       499,
		StartDate:   time.Now(),
		EndDate:     time.Now().AddDate(0, 1, 0),
	}

	mock.ExpectExec("UPDATE subscriptions").
		WithArgs(subscription.Price, subscription.StartDate.Format("2006-01"), subscription.EndDate.Format("2006-01"), subscription.UserID, subscription.ServiceName).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	updatedSubscription, err := db.UpdateSubscription(subscription)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if updatedSubscription == nil || updatedSubscription.UserID != subscription.UserID {
		t.Errorf("expected subscription for user %s, got %v", subscription.UserID, updatedSubscription)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestDeleteSubscription(t *testing.T) {
	db, mock, err := database.NewMock()
	if err != nil {
		t.Fatalf("failed to create mock database: %v", err)
	}
	defer db.Close()

	userID := uuid.New()
	serviceName := "Spotify"

	mock.ExpectExec("DELETE FROM subscriptions").
		WithArgs(userID, serviceName).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	err = db.DeleteSubscription(userID, serviceName)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestListSubscriptions(t *testing.T) {
	db, mock, err := database.NewMock()
	if err != nil {
		t.Fatalf("failed to create mock database: %v", err)
	}
	defer db.Close()

	userID := uuid.New()
	page := 1

	mock.ExpectQuery("SELECT").
		WithArgs(userID, page).
		WillReturnRows(pgxmock.NewRows([]string{"service_name", "user_id", "price", "start_date", "end_date"}).
			AddRow("Yandex Plus", userID, 399, "2024-01", "2024-12").
			AddRow("Kion", userID, 199, "2024-02", "2024-11"))

	subscriptions, err := db.ListSubscriptions(userID, page)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(subscriptions) != 2 {
		t.Errorf("expected 2 subscriptions, got %d", len(subscriptions))
	}
	if subscriptions[0].UserID != userID || subscriptions[1].UserID != userID {
		t.Errorf("expected subscriptions for user %s, got %v and %v", userID, subscriptions[0], subscriptions[1])
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestSubscriptionTotalCost(t *testing.T) {
	db, mock, err := database.NewMock()
	if err != nil {
		t.Fatalf("failed to create mock database: %v", err)
	}
	defer db.Close()

	userID := uuid.New()
	serviceName := "Yandex Plus"
	startDate := "2024-01"
	endDate := "2024-12"

	mock.ExpectQuery("SELECT").
		WithArgs(endDate, startDate, userID, endDate).
		WillReturnRows(pgxmock.NewRows([]string{"total_price"}).
			AddRow(299))

	totalCost, err := db.SubscriptionTotalCost(userID, serviceName, startDate, endDate)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if totalCost.TotalCost != 299 {
		t.Errorf("expected total cost 299, got %d", totalCost)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

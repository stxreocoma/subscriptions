package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"subscriptions/internal/api"
	"subscriptions/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestGetSubscriptionHandler(t *testing.T) {
	userID := uuid.New()
	serviceName := "serviceName"

	mock := NewMockStorage()

	mock.mock.On("GetSubscription", userID, serviceName).Return(&models.Subscription{
		UserID:      userID,
		ServiceName: serviceName,
	}, nil)

	api := &api.API{
		Router:  nil,
		Storage: mock,
	}

	server := httptest.NewServer(api.GetSubscriptionHandler())
	defer server.Close()

	client := server.Client()
	resp, err := client.Get(server.URL + "/subscription/" + userID.String() + "/" + serviceName)
	require.NoError(t, err, "Failed to send request")
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "Failed to read response body")

	require.Equal(t, http.StatusOK, resp.StatusCode, "Expected status OK")

	var subscription models.Subscription
	err = json.Unmarshal(respBody, &subscription)
	require.NoError(t, err, "Failed to unmarshal response body")

	mock.mock.AssertCalled(t, "GetSubscription", userID, serviceName)
	mock.mock.AssertNumberOfCalls(t, "GetSubscription", 1)
	mock.mock.AssertExpectations(t)

}

func TestCreateSubscriptionHandler(t *testing.T) {
	sub := &models.Subscription{
		UserID:      uuid.New(),
		ServiceName: "Netflix",
		Price:       999,
		StartDate:   time.Now(),
		EndDate:     time.Now().AddDate(0, 1, 0),
	}

	mock := NewMockStorage()
	mock.mock.On("CreateSubscription", sub).Return(sub, nil)

	api := &api.API{
		Router:  nil,
		Storage: mock,
	}

	server := httptest.NewServer(api.CreateSubscriptionHandler())
	defer server.Close()

	client := server.Client()
	body, err := json.Marshal(sub)
	require.NoError(t, err, "Failed to marshal body")
	resp, err := client.Post(server.URL+"/subscription", "application/json", bytes.NewReader(body))
	require.NoError(t, err, "Failed to send request")
	defer resp.Body.Close()

	require.Equal(t, http.StatusCreated, resp.StatusCode, "Expected status Created")

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "Failed to read response body")

	var createdSub models.Subscription
	err = json.Unmarshal(respBody, &createdSub)
	require.NoError(t, err, "Failed to unmarshal response body")

	mock.mock.AssertCalled(t, "CreateSubscription", sub)
	mock.mock.AssertNumberOfCalls(t, "CreateSubscription", 1)
	mock.mock.AssertExpectations(t)
}

func TestUpdateSubscriptionHandler(t *testing.T) {
	sub := &models.Subscription{
		UserID:      uuid.New(),
		ServiceName: "Spotify",
		Price:       499,
		StartDate:   time.Now(),
		EndDate:     time.Now().AddDate(0, 3, 0),
	}

	mock := NewMockStorage()
	mock.mock.On("UpdateSubscription", sub).Return(sub, nil)

	api := &api.API{
		Router:  nil,
		Storage: mock,
	}

	server := httptest.NewServer(api.UpdateSubscriptionHandler())
	defer server.Close()

	client := server.Client()

	body, err := json.Marshal(sub)
	require.NoError(t, err, "Failed to marshal body")

	req, err := http.NewRequest(http.MethodPut, server.URL+"/subscription/"+sub.UserID.String()+"/"+sub.ServiceName, bytes.NewReader(body))
	require.NoError(t, err, "Failed to create requesr")

	req.Header.Set("Content-Type", "application-json")

	resp, err := client.Do(req)
	require.NoError(t, err, "Failed to send request")

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "Failed to read response body")

	var updatedSub models.Subscription
	err = json.Unmarshal(respBody, &updatedSub)
	require.NoError(t, err, "Failed to unmarshal response body")

	mock.mock.AssertCalled(t, "UpdateSubscription", sub)
	mock.mock.AssertNumberOfCalls(t, "UpdateSubscription", 1)
	mock.mock.AssertExpectations(t)
}

func TestDeleteSubscriptionHandler(t *testing.T) {
	userID := uuid.New()
	serviceName := "Yandex Music"

	mock := NewMockStorage()
	mock.mock.On("DeleteSubscription", userID, serviceName).Return(nil)

	api := &api.API{
		Router:  nil,
		Storage: mock,
	}

	server := httptest.NewServer(api.DeleteSubscriptionHandler())
	defer server.Close()

	client := server.Client()
	req, err := http.NewRequest(http.MethodDelete, server.URL+"/subscription/"+userID.String()+"/"+serviceName, nil)
	require.NoError(t, err, "Failed to create request")

	resp, err := client.Do(req)
	require.NoError(t, err, "Failed to send request")
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode, "Expected status OK")

	mock.mock.AssertCalled(t, "DeleteSubscription", userID, serviceName)
	mock.mock.AssertNumberOfCalls(t, "DeleteSubscription", 1)
	mock.mock.AssertExpectations(t)
}

func TestListSubscriptionsHandler(t *testing.T) {
	userID := uuid.New()
	subs := []*models.Subscription{
		{
			UserID:      userID,
			ServiceName: "HBO Max",
			Price:       1499,
			StartDate:   time.Now(),
			EndDate:     time.Now().AddDate(0, 1, 0)},
		{
			UserID:      userID,
			ServiceName: "Amazon Prime",
			Price:       1299,
			StartDate:   time.Now(),
		},
	}

	mock := NewMockStorage()
	mock.mock.On("ListSubscriptions", userID).Return(subs, nil)

	api := &api.API{
		Router:  nil,
		Storage: mock,
	}

	server := httptest.NewServer(api.ListSubscriptionsHandler())
	defer server.Close()

	client := server.Client()

	resp, err := client.Get(server.URL + "/subscriptions?userID=" + userID.String())
	require.NoError(t, err, "Failed to send request")

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "Failed to read response body")

	subsList := make([]*models.Subscription, 2)
	err = json.Unmarshal(respBody, &subs)
	require.NoError(t, err, "Failed to unmarshal response body")

	require.Equal(t, http.StatusOK, resp.StatusCode, "Expected status OK")
	mock.mock.AssertCalled(t, "ListSubscriptions", userID)
	mock.mock.AssertNumberOfCalls(t, "ListSubscriptions", 1)
	mock.mock.AssertExpectations(t)
	require.Len(t, subsList, 2, "Expected 2 subscriptions in response")
	require.Equal(t, subsList[0].UserID, userID, "Expected userID to match")
}

func TestSubscriptionTotalCostHandler(t *testing.T) {
	userID := uuid.New()
	serviceName := "VK+Music"
	startDate := "2025-03"
	endDate := "2025-07"

	mock := NewMockStorage()
	mock.mock.On("GetSubscriptionTotalCost", userID, serviceName, startDate, endDate).Return(&models.TotalSubscriptionCost{TotalCost: 1497})

	api := &api.API{
		Router:  nil,
		Storage: mock,
	}

	server := httptest.NewServer(api.GetSubscriptionTotalCostHandler())
	defer server.Close()

	client := server.Client()
	resp, err := client.Get(server.URL + "/subscriptions/total/" + userID.String() + "?service_name=" + serviceName + "&start_date=" + startDate + "&end_date=" + endDate)
	require.NoError(t, err, "Failed to send request")
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "Failed to read response body")

	var totalCost models.TotalSubscriptionCost

	err = json.Unmarshal(body, &totalCost)
	require.NoError(t, err, "Failed to unmarshal body")

	require.Equal(t, http.StatusOK, resp.Status, "Expected status OK")

	mock.mock.AssertCalled(t, "SubscriptionTotalCost", userID, serviceName, startDate, endDate)
	mock.mock.AssertNumberOfCalls(t, "SubscriptionTotalCost", 1)
	mock.mock.AssertExpectations(t)
}

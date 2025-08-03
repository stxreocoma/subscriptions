package api

import (
	"subscriptions/internal/storage"

	"github.com/go-chi/chi/v5"
)

type API struct {
	Router  chi.Router
	Storage storage.Storage
}

func NewAPI(storage storage.Storage) *API {
	api := &API{
		Router:  chi.NewRouter(),
		Storage: storage,
	}

	api.InitRoutes()
	return api
}

func (a *API) InitRoutes() {
	a.Router.Get("/subscription/{userID}/{serviceName}", a.GetSubscriptionHandler())
	a.Router.Post("/subscription", a.CreateSubscriptionHandler())
	a.Router.Put("/subscription/{userID}/{serviceName}", a.UpdateSubscriptionHandler())
	a.Router.Delete("/subscription/{userID}/{serviceName}", a.DeleteSubscriptionHandler())
	a.Router.Get("/subscriptions", a.ListSubscriptionsHandler())
	a.Router.Get("/subscription/total-cost", a.GetSubscriptionTotalCostHandler())
}

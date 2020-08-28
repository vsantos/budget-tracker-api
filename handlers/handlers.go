package handlers

import (
	"budget-tracker/controllers"
	"net/http"
)

// Handlers struct defines all handlers used by backend
type Handlers struct {
	HealthCheckHandler http.Handler

	CreateUserHandler http.Handler
	GetUserHandler    http.Handler
	DeleteUserHandler http.Handler

	CreateCardHandler  http.Handler
	GetAllCardsHandler http.Handler
	DeleteCardHandler  http.Handler
	GetCardsHandler    http.Handler

	CreateBalanceHandler http.Handler
	GetBalanceHandler    http.Handler

	CreateSpendHandler http.Handler
}

// GetHandlers will return all backend handlers initialized
func GetHandlers() (h Handlers) {
	h.HealthCheckHandler = http.HandlerFunc(controllers.HealthCheck)
	h.CreateUserHandler = http.HandlerFunc(controllers.CreateUserEndpoint)
	h.GetUserHandler = http.HandlerFunc(controllers.GetUserEndpoint)
	h.DeleteUserHandler = http.HandlerFunc(controllers.DeleteUserEndpoint)
	h.CreateCardHandler = http.HandlerFunc(controllers.CreateCardEndpoint)
	h.GetAllCardsHandler = http.HandlerFunc(controllers.GetAllCardsEndpoint)
	h.DeleteCardHandler = http.HandlerFunc(controllers.DeleteCardEndpoint)
	h.GetCardsHandler = http.HandlerFunc(controllers.GetCardsEndpoint)
	h.CreateBalanceHandler = http.HandlerFunc(controllers.CreateBalanceEndpoint)
	h.GetBalanceHandler = http.HandlerFunc(controllers.GetBalanceEndpoint)
	h.CreateSpendHandler = http.HandlerFunc(controllers.CreateSpendEndpoint)
	return h
}

package handlers

import (
	"budget-tracker-api/controllers"
	"net/http"
)

// Handlers struct defines all handlers used by backend
type Handlers struct {
	HealthCheckHandler http.Handler
	SwaggerHandler     http.Handler

	OptionsOauthJWTTokenHandler http.Handler
	CreateOauthJWTTokenHandler  http.Handler
	GoogleCallbackHandler       http.Handler

	OptionsJWTTokenHandler http.Handler
	CreateJWTTokenHandler  http.Handler
	RefreshJWTTokenHandler http.Handler

	GetUsersHandler   http.Handler
	CreateUserHandler http.Handler
	GetUserHandler    http.Handler
	DeleteUserHandler http.Handler

	OptionsCardsHandler http.Handler
	CreateCardHandler   http.Handler
	GetAllCardsHandler  http.Handler
	DeleteCardHandler   http.Handler
	GetCardsHandler     http.Handler

	CreateBalanceHandler http.Handler
	GetBalanceHandler    http.Handler

	GetSpendsHandler   http.Handler
	CreateSpendHandler http.Handler
}

// GetHandlers will return all backend handlers initialized
func GetHandlers() (h Handlers) {
	h.HealthCheckHandler = http.HandlerFunc(controllers.HealthCheck)
	h.SwaggerHandler = http.HandlerFunc(controllers.Swagger)

	h.OptionsOauthJWTTokenHandler = http.HandlerFunc(controllers.Oauth2TokenOptionsEndpoint)
	h.CreateOauthJWTTokenHandler = http.HandlerFunc(controllers.CreateOauthTokenEndpoint)
	h.GoogleCallbackHandler = http.HandlerFunc(controllers.GoogleCallbackEndpoint)
	h.OptionsJWTTokenHandler = http.HandlerFunc(controllers.JWTTokenOptionsEndpoint)
	h.CreateJWTTokenHandler = http.HandlerFunc(controllers.CreateJWTTokenEndpoint)
	h.RefreshJWTTokenHandler = http.HandlerFunc(controllers.CreateJWTTokenEndpoint)

	h.CreateUserHandler = http.HandlerFunc(controllers.CreateUserEndpoint)
	h.GetUsersHandler = http.HandlerFunc(controllers.GetUsersEndpoint)
	h.GetUserHandler = http.HandlerFunc(controllers.GetUserEndpoint)
	h.DeleteUserHandler = http.HandlerFunc(controllers.DeleteUserEndpoint)

	h.OptionsCardsHandler = http.HandlerFunc(controllers.CardsOptionsEndpoint)
	h.CreateCardHandler = http.HandlerFunc(controllers.CreateCardEndpoint)
	h.GetAllCardsHandler = http.HandlerFunc(controllers.GetAllCardsEndpoint)
	h.DeleteCardHandler = http.HandlerFunc(controllers.DeleteCardEndpoint)
	h.GetCardsHandler = http.HandlerFunc(controllers.GetCardsEndpoint)

	h.CreateBalanceHandler = http.HandlerFunc(controllers.CreateBalanceEndpoint)
	h.GetBalanceHandler = http.HandlerFunc(controllers.GetBalanceEndpoint)

	h.GetSpendsHandler = http.HandlerFunc(controllers.GetSpendsEndpoint)
	h.CreateSpendHandler = http.HandlerFunc(controllers.CreateSpendEndpoint)
	return h
}

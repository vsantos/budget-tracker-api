package routes

import (
	"budget-tracker/handlers"

	"github.com/gorilla/mux"
)

// InitRoutes will initiate all routes
func InitRoutes(router *mux.Router) {
	m := handlers.GetMiddlewares()
	h := handlers.GetHandlers()

	router.Handle("/health", h.HealthCheckHandler).Methods("GET")

	router.Handle("/api/v1/jwt/issue", m.JSON(h.CreateJWTTokenHandler)).Methods("POST")
	router.Handle("/api/v1/jwt/refresh", m.JSON(h.CreateJWTTokenHandler)).Methods("POST")

	router.Handle("/api/v1/users", m.JSON(m.Auth(h.CreateUserHandler))).Methods("POST")
	router.Handle("/api/v1/users", m.JSON(m.Auth(h.GetUsersHandler))).Methods("GET")
	router.Handle("/api/v1/users/{id}", m.JSON(m.Auth(h.GetUserHandler))).Methods("GET")
	router.Handle("/api/v1/users/{id}", m.JSON(m.Auth(h.DeleteUserHandler))).Methods("DELETE")

	router.Handle("/api/v1/cards", m.JSON(m.Auth(h.CreateCardHandler))).Methods("POST")
	router.Handle("/api/v1/cards", m.JSON(m.Auth(h.GetAllCardsHandler))).Methods("GET")
	router.Handle("/api/v1/cards/{id}", m.JSON(m.Auth(h.DeleteCardHandler))).Methods("DELETE")
	router.Handle("/api/v1/cards/{owner_id}", m.JSON(m.Auth(h.GetCardsHandler))).Methods("GET")

	router.Handle("/api/v1/balance", m.JSON(m.Auth(h.CreateBalanceHandler))).Methods("POST")
	router.Handle("/api/v1/balance/{owner_id}", m.JSON(m.Auth(h.GetBalanceHandler))).Methods("GET")

	router.Handle("/api/v1/spends", m.JSON(m.Auth(h.CreateSpendHandler))).Methods("POST")
	router.Handle("/api/v1/spends/{owner_id}", m.JSON(m.Auth(h.GetSpendsHandler))).Methods("GET")
}

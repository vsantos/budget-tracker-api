package routes

import (
	"github.com/gorilla/mux"
)

// InitRoutes will initiate all routes
func InitRoutes(router *mux.Router) {
	router.HandleFunc("/health", HealthCheck).Methods("GET")

	router.HandleFunc("/api/v1/user", CreateUserEndpoint).Methods("POST")
	router.HandleFunc("/api/v1/user/{id}", GetUserEndpoint).Methods("GET")
	router.HandleFunc("/api/v1/user/{id}", DeleteUserEndpoint).Methods("DELETE")
	router.HandleFunc("/api/v1/users", GetUsersEndpoint).Methods("GET")

	router.HandleFunc("/api/v1/cards", CreateCardEndpoint).Methods("POST")
	router.HandleFunc("/api/v1/cards", GetAllCardsEndpoint).Methods("GET")
	router.HandleFunc("/api/v1/cards/{id}", DeleteCardEndpoint).Methods("DELETE")
	router.HandleFunc("/api/v1/cards/{owner_id}", GetCardsEndpoint).Methods("GET")

	router.HandleFunc("/api/v1/balance", CreateBalanceEndpoint).Methods("POST")
	router.HandleFunc("/api/v1/balance/{owner_id}", GetBalanceEndpoint).Methods("GET")
}

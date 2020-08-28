package routes

import (
	"budget-tracker/controllers"
	"budget-tracker/handlers"

	"github.com/gorilla/mux"
)

// InitRoutes will initiate all routes
func InitRoutes(router *mux.Router) {
	m := handlers.GetMiddlewares()
	h := handlers.GetHandlers()

	router.Handle("/health", h.HealthCheckHandler).Methods("GET")

	router.Handle("/api/v1/auth/token", m.JSON(m.Auth(h.HealthCheckHandler))).Methods("GET")

	router.HandleFunc("/api/v1/user", controllers.CreateUserEndpoint).Methods("POST")
	router.HandleFunc("/api/v1/user/{id}", controllers.GetUserEndpoint).Methods("GET")
	router.HandleFunc("/api/v1/user/{id}", controllers.DeleteUserEndpoint).Methods("DELETE")
	router.HandleFunc("/api/v1/users", controllers.GetUsersEndpoint).Methods("GET")

	router.HandleFunc("/api/v1/cards", controllers.CreateCardEndpoint).Methods("POST")
	router.HandleFunc("/api/v1/cards", controllers.GetAllCardsEndpoint).Methods("GET")
	router.HandleFunc("/api/v1/cards/{id}", controllers.DeleteCardEndpoint).Methods("DELETE")
	router.HandleFunc("/api/v1/cards/{owner_id}", controllers.GetCardsEndpoint).Methods("GET")

	router.HandleFunc("/api/v1/balance", controllers.CreateBalanceEndpoint).Methods("POST")
	router.HandleFunc("/api/v1/balance/{owner_id}", controllers.GetBalanceEndpoint).Methods("GET")

	router.HandleFunc("/api/v1/spends", controllers.CreateSpendEndpoint).Methods("POST")
}

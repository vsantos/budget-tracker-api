package main

import (
	"budget-tracker/routes"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

const port = ":8080"

func init() {
	log.SetFormatter(&log.JSONFormatter{})
}

func main() {
	log.Infoln("Started Application at port", port)

	router := mux.NewRouter()

	router.HandleFunc("/health", routes.HealthCheck).Methods("GET")
	router.HandleFunc("/api/v1/user", routes.CreateUserEndpoint).Methods("POST")
	router.HandleFunc("/api/v1/user/{id}", routes.GetUserEndpoint).Methods("GET")
	router.HandleFunc("/api/v1/user/{id}", routes.DeleteUserEndpoint).Methods("DELETE")
	router.HandleFunc("/api/v1/users", routes.GetUsersEndpoint).Methods("GET")
	router.HandleFunc("/api/v1/cards", routes.CreateCardEndpoint).Methods("POST")
	router.HandleFunc("/api/v1/cards", routes.GetAllCardsEndpoint).Methods("GET")
	router.HandleFunc("/api/v1/cards/{owner_id}", routes.GetCardsEndpoint).Methods("GET")
	router.HandleFunc("/api/v1/cards/{id}", routes.DeleteCardEndpoint).Methods("DELETE")

	http.ListenAndServe(port, router)
}

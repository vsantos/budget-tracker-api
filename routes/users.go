package routes

import (
	"budget-tracker/models"
	"budget-tracker/services"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// HealthCheck will validate if external core components are working
func HealthCheck(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")

	err := services.DatabaseHealth()
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"mongodb": "unhealthy", "details": "` + err.Error() + `"}`))
		return
	}

	response.Write([]byte(`{"mongodb": "healthy"}`))
}

// CreateUserEndpoint creates an user
func CreateUserEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")

	var user models.User

	_ = json.NewDecoder(request.Body).Decode(&user)

	err := models.CreateUser(&user)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "could not create user", "details": "` + err.Error() + `"}`))
		return
	}

	response.Write([]byte(`{"message": "created user '` + user.Login + `'"}`))
}

// DeleteUserEndpoint deletes an user
func DeleteUserEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")

	params := mux.Vars(request)

	err := models.DeleteUser(params["login"])
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "could not delete user", "details": "` + err.Error() + `"}`))
		return
	}

	response.Write([]byte(`{"message": "deleted user '` + params["login"] + `'"}`))
}

// GetUserEndpoint an unique user
func GetUserEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	params := mux.Vars(request)

	user, err := models.GetUser(params["login"])
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}

	json.NewEncoder(response).Encode(user)
}

// GetUsersEndpoint returns a collection of user
func GetUsersEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")

	users, err := models.GetUsers()
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}

	json.NewEncoder(response).Encode(users)
}

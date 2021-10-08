package controllers

import (
	"budget-tracker-api/models"
	"budget-tracker-api/repository"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// CreateUserEndpoint creates an user
func CreateUserEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")

	var user repository.User

	_ = json.NewDecoder(request.Body).Decode(&user)

	result, err := models.CreateUser(request.Context(), user)
	if err != nil {
		if strings.Contains(err.Error(), "user already exists") {
			response.WriteHeader(http.StatusConflict)
			response.Write([]byte(`{"message": "could not create user", "details": "` + err.Error() + `"}`))
			return
		}

		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "could not create user", "details": "` + err.Error() + `"}`))
		return
	}

	response.WriteHeader(http.StatusCreated)
	response.Write([]byte(`{"message": "created user '` + user.Login + `'", "id": "` + result + `"}`))
}

// GetUsersEndpoint returns a collection of user
func GetUsersEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")

	users, err := models.GetUsers(request.Context())
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}

	if len(users) == 0 {
		response.Write([]byte(`[]`))
		return
	}

	json.NewEncoder(response).Encode(users)
}

// GetUserEndpoint an unique user
func GetUserEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	params := mux.Vars(request)

	user, err := models.GetUser(request.Context(), params["id"])
	if err != nil {
		if strings.Contains(err.Error(), "could not find user") {
			response.WriteHeader(http.StatusNotFound)
			response.Write([]byte(`{"message": "could not find user", "id": "` + params["id"] + `"}`))
			return
		}
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}

	json.NewEncoder(response).Encode(user)
}

// DeleteUserEndpoint deletes an user
func DeleteUserEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")

	params := mux.Vars(request)

	err := models.DeleteUser(request.Context(), params["id"])
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "could not delete user", "details": "` + err.Error() + `"}`))
		return
	}

	response.WriteHeader(http.StatusCreated)
	response.Write([]byte(`{"message": "deleted user '` + params["id"] + `'"}`))
}

package controllers

import (
	"net/http"
)

// JWTTokenOptionsEndpoint will return a set of headers for UI purposes
func JWTTokenOptionsEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	response.Header().Set("Access-Control-Allow-Origin", "*")
	response.Header().Set("Access-Control-Allow-Credentials", "true")
	response.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	response.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

// Oauth2TokenOptionsEndpoint will return a set of headers for UI purposes
func Oauth2TokenOptionsEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	response.Header().Set("Access-Control-Allow-Origin", "*")
	response.Header().Set("Access-Control-Allow-Credentials", "true")
	response.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, PUT, PATCH, OPTIONS")
	response.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

// CardsOptionsEndpoint will return a set of headers for UI purposes
func CardsOptionsEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	response.Header().Set("Access-Control-Allow-Origin", "*")
	response.Header().Set("Access-Control-Allow-Credentials", "true")
	response.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, PUT, PATCH, OPTIONS")
	response.Header().Set("Access-Control-Allow-Headers", "Content-Type, api_key, Authorization")
}

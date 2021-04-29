package controllers

import "net/http"

// Swagger will return a swagger json file
func Swagger(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/yaml")
	response.Header().Set("Access-Control-Allow-Origin", "*")
	response.Header().Set("Access-Control-Allow-Credentials", "true")
	response.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, PUT, PATCH, OPTIONS")
	response.Header().Set("Access-Control-Allow-Headers", "Content-Type, api_key, Authorization")
	http.ServeFile(response, request, "/app/docs/swagger.yaml")
}

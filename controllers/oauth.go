package controllers

import (
	"net/http"
)

// validateProvider validates if a given provider is supported
func validateProvider(provider string) bool {
	providers := []string{"google"}
	for _, p := range providers {
		if provider == p {
			return true
		}
	}

	return false
}

// CreateOauthTokenEndpoint creates a token based on user credentials
func CreateOauthTokenEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	response.Header().Set("Access-Control-Allow-Origin", "*")

	v := request.URL.Query()

	ifProvider := validateProvider(v.Get("provider"))
	if !ifProvider {
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte(`{"message": "non-supported oauth2 provider '` + v.Get("provider") + `'"}`))
		return
	}

	url := googleOauthConfig.AuthCodeURL("pseudo-random")
	http.Redirect(response, request, url, http.StatusTemporaryRedirect)

	response.WriteHeader(http.StatusInternalServerError)
	response.Write([]byte(`{"message": "google oauth"}`))
	return
}

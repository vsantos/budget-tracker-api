package controllers

import "net/http"

// CreateSpendEndpoint will create a spend and add to the current month balance
func CreateSpendEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	response.Header().Add("backend", "budget-tracker")

}

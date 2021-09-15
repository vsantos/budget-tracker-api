package controllers

import (
	"budget-tracker-api/services"

	"net/http"
)

// HealthCheck will validate if external core components are working
func HealthCheck(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")

	var d services.DatabaseManager
	d = services.MongoCfg{
		URI: services.DatabaseURI,
	}
	err := d.Health()
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"database": "unhealthy", "details": "` + err.Error() + `"}`))
		return
	}

	response.Write([]byte(`{"database": "healthy"}`))
}

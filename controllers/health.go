package controllers

import (
	"budget-tracker/services"
	"context"

	"net/http"

	oteltrace "go.opentelemetry.io/otel/trace"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

var tracer = otel.Tracer("budget-tracker-api")

// HealthCheck will validate if external core components are working
func HealthCheck(response http.ResponseWriter, request *http.Request) {
	ctx := context.TODO()
	_, span := tracer.Start(ctx, "getUser", oteltrace.WithAttributes(attribute.String("id", "1231")))
	defer span.End()
	response.Header().Add("content-type", "application/json")

	err := services.DatabaseHealth()
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"database": "unhealthy", "details": "` + err.Error() + `"}`))
		return
	}

	response.Write([]byte(`{"database": "healthy"}`))
}

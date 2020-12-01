package routes

import (
	"budget-tracker/handlers"

	"github.com/gorilla/mux"
)

// InitRoutes will initiate all routes
func InitRoutes(router *mux.Router) {
	m := handlers.GetMiddlewares()
	h := handlers.GetHandlers()

	// swagger:operation GET /health Healthchecks health
	//
	// Returns the API can be considered operational
	// ---
	// consumes:
	// - application/json
	// produces:
	// - application/json
	// responses:
	//   '200':
	//     description: healthy components
	//     examples:
	//       application/json: { "database": "healthy" }
	//   '500':
	//     description: unhealthy components
	//     examples:
	//       application/json: { "database": "unhealthy" }
	//     type: json
	router.Handle("/health", h.HealthCheckHandler).Methods("GET")
	router.Handle("/api/v1/swagger.yaml", h.SwaggerHandler).Methods("GET")

	// swagger:operation POST /api/v1/jwt/issue JWT issue
	//
	// Returns a JWT signed token to be used for the next 5 minutes
	// ---
	// consumes:
	// - application/json
	// produces:
	// - application/json
	// parameters:
	// - name: content-type
	//   in: headers
	//   description: application/json
	//   required: true
	// - name: body
	//   in: body
	//   description: tags to filter by
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/JWTUser"
	// responses:
	//   '201':
	//     description: returned JWT token
	//     examples:
	//       application/json: { "type": "bearer", "refresh": "<REFRESH_TOKEN>", "token": "<JWT_TOKEN>" }
	//   '401':
	//     description: invalid credentials
	//     examples:
	//       application/json: { "message": "invalid credentials for user 'vsantos'" }
	//     type: json
	router.Handle("/api/v1/jwt/issue", m.JSON(h.CreateJWTTokenHandler)).Methods("POST")
	router.Handle("/api/v1/jwt/refresh", m.JSON(h.CreateJWTTokenHandler)).Methods("POST")

	router.Handle("/api/v1/users", m.JSON(m.Auth(h.CreateUserHandler))).Methods("POST")
	router.Handle("/api/v1/users", m.JSON(m.Auth(h.GetUsersHandler))).Methods("GET")
	router.Handle("/api/v1/users/{id}", m.JSON(m.Auth(h.GetUserHandler))).Methods("GET")
	router.Handle("/api/v1/users/{id}", m.JSON(m.Auth(h.DeleteUserHandler))).Methods("DELETE")

	router.Handle("/api/v1/cards", m.JSON(m.Auth(h.CreateCardHandler))).Methods("POST")
	router.Handle("/api/v1/cards", m.JSON(m.Auth(h.GetAllCardsHandler))).Methods("GET")
	router.Handle("/api/v1/cards/{id}", m.JSON(m.Auth(h.DeleteCardHandler))).Methods("DELETE")
	router.Handle("/api/v1/cards/{owner_id}", m.JSON(m.Auth(h.GetCardsHandler))).Methods("GET")

	router.Handle("/api/v1/balance", m.JSON(m.Auth(h.CreateBalanceHandler))).Methods("POST")
	router.Handle("/api/v1/balance/{owner_id}", m.JSON(m.Auth(h.GetBalanceHandler))).Methods("GET")

	router.Handle("/api/v1/spends", m.JSON(m.Auth(h.CreateSpendHandler))).Methods("POST")
	router.Handle("/api/v1/spends/{owner_id}", m.JSON(m.Auth(h.GetSpendsHandler))).Methods("GET")
}

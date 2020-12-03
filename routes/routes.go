package routes

import (
	"budget-tracker/handlers"

	"github.com/gorilla/mux"
)

// InitRoutes will initiate all routes
func InitRoutes(router *mux.Router) {
	m := handlers.GetMiddlewares()
	h := handlers.GetHandlers()

	// swagger:operation GET /health Utils get
	//
	// Returns the API can be considered operational
	// ---
	// produces:
	// - application/json
	// responses:
	//   '200':
	//     description: healthy components
	//     examples:
	//       application/json: { "database": "healthy" }
	//     type: json
	//   '500':
	//     description: unhealthy components
	//     examples:
	//       application/json: { "database": "unhealthy" }
	//     type: json
	router.Handle("/health", h.HealthCheckHandler).Methods("GET")

	// swagger:operation GET /api/v1/swagger.yaml Utils get
	//
	// Returns the swagger yaml file to be get by swagger-ui or similar
	// ---
	// produces:
	// - application/yaml
	// responses:
	//   '200':
	//     description: found swagger document
	//   '404':
	//     description: could not find swagger document
	//     type: json
	router.Handle("/api/v1/swagger.yaml", h.SwaggerHandler).Methods("GET")

	// swagger:operation POST /api/v1/jwt/issue Authentication issue
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
	//   description: credentials
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/JWTUser"
	// responses:
	//   '201':
	//     description: returned JWT token
	//     examples:
	//       application/json: { "type": "bearer", "refresh": "<REFRESH_TOKEN>", "token": "<JWT_TOKEN>" }
	//     type: json
	//   '400':
	//     description: bad request (missing one of params)
	//     examples:
	//       application/json: { "message": "empty required payload attributes" }
	//     type: json
	//   '401':
	//     description: invalid credentials
	//     examples:
	//       application/json: { "message": "invalid credentials for user 'vsantos'" }
	//     type: json
	router.Handle("/api/v1/jwt/issue", m.JSON(h.CreateJWTTokenHandler)).Methods("POST")
	router.Handle("/api/v1/jwt/refresh", m.JSON(h.CreateJWTTokenHandler)).Methods("POST")

	// swagger:operation POST /api/v1/users Users create
	//
	// Creates an user
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
	//   description: user payload
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/User"
	// responses:
	//   '201':
	//     description: returned user
	//     examples:
	//       application/json: { "message": "created user '<USER_LOGIN>'", "id": "<USER_ID>" }
	//     type: json
	//   '500':
	//     description: internal server error
	//     examples:
	//       application/json: { "message": "could not create user", "details": "<ERROR_DETAILS>" }
	//     type: json
	router.Handle("/api/v1/users", m.JSON(m.Auth(h.CreateUserHandler))).Methods("POST")

	// swagger:operation GET /api/v1/users Users list
	//
	// List all users
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
	// responses:
	//   '200':
	//     description: users response
	//     schema:
	//       type: array
	//       items:
	//         "$ref": "#/definitions/SanitizedUser"
	//     type: json
	//   '500':
	//     description: internal server error
	//     examples:
	//       application/json: { "message": "<ERROR_DETAILS>" }
	//     type: json
	router.Handle("/api/v1/users", m.JSON(m.Auth(h.GetUsersHandler))).Methods("GET")

	// swagger:operation GET /api/v1/users/{id} Users get
	//
	// List a single user
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
	// - name: user_d
	//   in: user_id
	//   description: user id
	//   required: true
	// responses:
	//   '200':
	//     description: user response
	//     schema:
	//       type: json
	//       items:
	//         "$ref": "#/definitions/SanitizedUser"
	//   '500':
	//     description: internal server error
	//     examples:
	//       application/json: { "message": "<ERROR_DETAILS>" }
	//     type: json
	router.Handle("/api/v1/users/{id}", m.JSON(m.Auth(h.GetUserHandler))).Methods("GET")

	// swagger:operation DELETE /api/v1/users/{id} Users delete
	//
	// Delete a single user
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
	// - name: id
	//   in: id
	//   description: user id
	//   required: true
	// responses:
	//   '201':
	//     description: deleted user
	//     examples:
	//       application/json: { "message": "deleted user '<USER_ID:>'" }
	//     type: json
	//   '500':
	//     description: internal server error
	//     examples:
	//       application/json: { "message": "could not delete user", "details": "<ERROR_DETAILS>" }
	//     type: json
	router.Handle("/api/v1/users/{id}", m.JSON(m.Auth(h.DeleteUserHandler))).Methods("DELETE")

	// swagger:operation POST /api/v1/cards Cards create
	//
	// Creates a single card
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
	//   description: cards payload
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/CreditCard"
	// responses:
	//   '201':
	//     description: deleted user
	//     examples:
	//       application/json: { "message": "created card '<CARD_ALIAS>'", "id": "<CARD_ID>" }
	//     type: json
	//   '500':
	//     description: internal server error
	//     examples:
	//       application/json: { "message": "could not create card", "details": "given network '<CARD_NETWORK>' is not a valid one" }
	//     type: json
	//   '500':
	//     description: internal server error
	//     examples:
	//       application/json: { "message": "could not create card", "details": "<ERROR_DETAILS>" }
	//     type: json
	router.Handle("/api/v1/cards", m.JSON(m.Auth(h.CreateCardHandler))).Methods("POST")

	// swagger:operation GET /api/v1/cards Cards list
	//
	// List all cards from platform
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: content-type
	//   in: headers
	//   description: application/json
	//   required: true
	// responses:
	//   '200':
	//     description: card response
	//     schema:
	//       type: array
	//       items:
	//         "$ref": "#/definitions/CreditCard"
	//   '500':
	//     description: internal server error
	//     examples:
	//       application/json: {"message": "<ERROR_DETAILS>"}
	//     type: json
	router.Handle("/api/v1/cards", m.JSON(m.Auth(h.GetAllCardsHandler))).Methods("GET")

	// swagger:operation DELETE /api/v1/cards/{id} Cards delete
	//
	// Deletes a single card
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
	// - name: card_id
	//   in: card_id
	//   description: card id
	//   required: true
	// responses:
	//   '201':
	//     description: deleted card
	//     examples:
	//       application/json: { "message": "deleted card '<CARD_ID>'" }
	//     type: json
	//   '500':
	//     description: internal server error
	//     examples:
	//       application/json: { "message": "could not delete card", "details": "<ERROR_DETAILS>" }
	router.Handle("/api/v1/cards/{id}", m.JSON(m.Auth(h.DeleteCardHandler))).Methods("DELETE")

	// swagger:operation GET /api/v1/cards/{owner_id} Cards list
	//
	// List all cards from a given owner
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: content-type
	//   in: headers
	//   description: application/json
	//   required: true
	// - name: owner_id
	//   in: owner_id
	//   description: owner id
	//   required: true
	// responses:
	//   '200':
	//     description: card response
	//     schema:
	//       type: array
	//       items:
	//         "$ref": "#/definitions/CreditCard"
	//   '500':
	//     description: internal server error
	//     examples:
	//       application/json: {"message": "<ERROR_DETAILS>"}
	//     type: json
	router.Handle("/api/v1/cards/{owner_id}", m.JSON(m.Auth(h.GetCardsHandler))).Methods("GET")

	// swagger:operation POST /api/v1/balance Balance create
	//
	// Creates a single balance for a given owner
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
	// - name: owner_id
	//   in: owner_id
	//   description: balance owner_id
	//   required: true
	// responses:
	//   '201':
	//     description: deleted user
	//     examples:
	//       application/json: { "message": "created balance", "id": "<BALANCE_ID>" }
	//     type: json
	//   '400':
	//     description: bad request
	//     examples:
	//       application/json: {"message": "could not create balance", "details": "balances must have an 'owner_id'"}
	//     type: json
	//   '500':
	//     description: internal server error
	//     examples:
	//       application/json: { "message": "could not delete balance", "details": "<ERROR_DETAILS>" }
	//     type: json
	router.Handle("/api/v1/balance", m.JSON(m.Auth(h.CreateBalanceHandler))).Methods("POST")

	// swagger:operation GET /api/v1/balance/{owner_id} Balance list
	//
	// List all balances from a given owner or a single one given a month and year as query params
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: content-type
	//   in: headers
	//   description: application/json
	//   required: true
	// - name: month
	//   in: query
	//   description: month
	// - name: year
	//   in: query
	//   description: year
	// responses:
	//   '200':
	//     description: balance response
	//     schema:
	//       type: array
	//       items:
	//         "$ref": "#/definitions/Balance"
	//   '404':
	//     description: balance not found
	//     examples:
	//       application/json: []
	//     type: json
	//   '500':
	//     description: internal server error
	//     examples:
	//       application/json: {"message": "<ERROR_DETAILS>"}
	//     type: json
	router.Handle("/api/v1/balance/{owner_id}", m.JSON(m.Auth(h.GetBalanceHandler))).Methods("GET")

	// swagger:operation POST /api/v1/spends Spends create
	//
	// Creates a single spend for a given owner
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
	// - name: owner_id
	//   in: owner_id
	//   description: spend owner_id
	//   required: true
	// responses:
	//   '201':
	//     description: deleted user
	//     examples:
	//       application/json: { "message": "created spend to user '<OWNER_ID>'", "id": "<SPEND_ID>"}
	//     type: json
	//   '400':
	//     description: bad request
	//     examples:
	//       application/json: {"message": "could not create spend", "details": "missing owner ID"}
	//     type: json
	//   '500':
	//     description: internal server error
	//     examples:
	//       application/json: { "message": "could not create spend", "details": "<ERROR_DETAILS>" }
	//     type: json
	router.Handle("/api/v1/spends", m.JSON(m.Auth(h.CreateSpendHandler))).Methods("POST")

	// swagger:operation GET /api/v1/spends/{owner_id} Spends list
	//
	// Get all spends for a given owner id
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: content-type
	//   in: headers
	//   description: application/json
	//   required: true
	// - name: owner_id
	//   in: owner_id
	//   description: owner id
	// responses:
	//   '200':
	//     description: spends response
	//     schema:
	//       type: array
	//       items:
	//         "$ref": "#/definitions/Spend"
	//   '500':
	//     description: internal server error
	//     examples:
	//       application/json: {"message": "<ERROR_DETAILS>"}
	//     type: json
	router.Handle("/api/v1/spends/{owner_id}", m.JSON(m.Auth(h.GetSpendsHandler))).Methods("GET")
}

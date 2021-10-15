package controllers

import (
	"budget-tracker-api/models"
	"budget-tracker-api/repository"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// validateCardNetwork will validate if a card's network is an allowed one
func validateCardNetwork(network string) bool {
	networks := []string{"visa", "mastercard", "elo", "vr", "ticket"}
	for _, n := range networks {
		if network == n {
			return true
		}
	}
	return false
}

// CreateCardEndpoint will create a single card to an user
func CreateCardEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	response.Header().Set("Access-Control-Allow-Origin", "*")

	var card repository.CreditCard

	_ = json.NewDecoder(request.Body).Decode(&card)

	validNetwork := validateCardNetwork(card.Network)
	if !validNetwork {
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte(`{"message": "could not create card", "details": "given network '` + card.Network + `' is not a valid one"}`))
		return
	}

	result, err := models.CreateCard(request.Context(), card)
	if err != nil {
		if strings.Contains(err.Error(), "card already exists") {
			response.WriteHeader(http.StatusConflict)
			response.Write([]byte(`{"message": "could not create card", "details": "` + err.Error() + `"}`))
			return
		}

		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "could not create card", "details": "` + err.Error() + `"}`))
		return
	}

	response.WriteHeader(http.StatusCreated)
	response.Write([]byte(`{"message": "created card '` + card.Alias + `'", "id": "` + result + `"}`))
}

// GetAllCardsEndpoint will return all cards from database
func GetAllCardsEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	response.Header().Set("Access-Control-Allow-Origin", "*")

	cards, err := models.GetAllCards(request.Context())
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}

	if len(cards) == 0 {
		response.Write([]byte(`[]`))
		return
	}

	json.NewEncoder(response).Encode(cards)
}

// GetCardsEndpoint will return all cards from a given user
func GetCardsEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	response.Header().Set("Access-Control-Allow-Origin", "*")

	params := mux.Vars(request)

	cards, err := models.GetCards(request.Context(), params["owner_id"])
	if err != nil {
		if strings.Contains(err.Error(), "could not find any cards") {
			response.WriteHeader(http.StatusNotFound)
			response.Write([]byte(`{"message": "` + err.Error() + `", "owner_id": "` + params["owner_id"] + `"}`))
			return
		}

		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}

	if len(cards) == 0 {
		response.Write([]byte(`[]`))
		return
	}

	json.NewEncoder(response).Encode(cards)
}

// DeleteCardEndpoint deletes a card given an ID
func DeleteCardEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	response.Header().Set("Access-Control-Allow-Origin", "*")

	params := mux.Vars(request)

	err := models.DeleteCard(request.Context(), params["id"])
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "could not delete card", "details": "` + err.Error() + `"}`))
		return
	}

	response.WriteHeader(http.StatusOK)
	response.Write([]byte(`{"message": "deleted card '` + params["id"] + `'"}`))
}

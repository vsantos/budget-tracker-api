package routes

import (
	"budget-tracker/models"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// CreateCardEndpoint will create a single card to an user
func CreateCardEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")

	var card models.CreditCard

	_ = json.NewDecoder(request.Body).Decode(&card)

	result, err := models.CreateCard(card)
	if err != nil {
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

	cards, err := models.GetAllCards()
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}

	json.NewEncoder(response).Encode(cards)
}

// GetCardsEndpoint will return all cards from a given user
func GetCardsEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	params := mux.Vars(request)

	cards, err := models.GetCards(params["owner_id"])
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}

	json.NewEncoder(response).Encode(cards)
}

// DeleteCardEndpoint deletes a card given an ID
func DeleteCardEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")

	params := mux.Vars(request)

	err := models.DeleteCard(params["id"])
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "could not delete card", "details": "` + err.Error() + `"}`))
		return
	}

	response.WriteHeader(http.StatusCreated)
	response.Write([]byte(`{"message": "deleted card '` + params["id"] + `'"}`))
}

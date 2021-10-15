package controllers

import (
	"budget-tracker-api/models"
	"budget-tracker-api/repository"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// CreateSpendEndpoint will create a spend and add to the current month balance
func CreateSpendEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	response.Header().Add("backend", "budget-tracker")

	var spend repository.Spend

	_ = json.NewDecoder(request.Body).Decode(&spend)

	if spend.OwnerID.Hex() == "000000000000000000000000" {
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte(`{"message": "could not create spend", "details": "missing owner ID"}`))
		return
	}

	result, err := models.CreateSpend(request.Context(), spend)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "could not create spend", "details": "` + err.Error() + `"}`))
		return
	}

	// add spend to balance

	response.WriteHeader(http.StatusCreated)
	response.Write([]byte(`{"message": "created spend", "owner_id": "` + spend.OwnerID.Hex() + `", "id": "` + result + `"}`))
}

// GetSpendsEndpoint will return all spends from an user
func GetSpendsEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	response.Header().Add("backend", "budget-tracker")

	params := mux.Vars(request)

	spends, err := models.GetSpends(request.Context(), params["owner_id"])
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}

	if len(spends) == 0 {
		response.Write([]byte(`[]`))
		return
	}

	json.NewEncoder(response).Encode(spends)
}

package routes

import (
	"budget-tracker/models"
	"encoding/json"
	"strconv"

	"github.com/gorilla/mux"

	"net/http"
)

// CreateBalanceEndpoint will create a balance to an user
func CreateBalanceEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")

	var balance models.Balance

	_ = json.NewDecoder(request.Body).Decode(&balance)

	if balance.OwnerID.Hex() == "000000000000000000000000" || balance.OwnerID.Hex() == "" {
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte(`{"message": "could not create balance", "details": "balances must have an 'owner_id'"}`))
		return
	}

	result, err := models.CreateBalance(balance)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "could not create balance", "details": "` + err.Error() + `"}`))
		return
	}

	response.WriteHeader(http.StatusCreated)
	response.Write([]byte(`{"message": "created balance", "id": "` + result + `"}`))
}

// GetBalanceEndpoint will return a balance from a given user
func GetBalanceEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")

	v := request.URL.Query()

	params := mux.Vars(request)
	month := v.Get("month")
	year := v.Get("year")

	// in case of a non-existent URL parameters
	if month == "" || year == "" {
		balances, err := models.GetAllBalances(params["owner_id"])
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{"message": "` + err.Error() + `"}`))
			return
		}

		if len(balances) == 0 {
			response.WriteHeader(http.StatusNotFound)
			response.Write([]byte(`[]`))
			return
		}

		json.NewEncoder(response).Encode(balances)
	}

	// in case of a existent URL parameters
	if month != "" && year != "" {
		imonth, _ := strconv.ParseInt(month, 10, 64)
		iyear, _ := strconv.ParseInt(year, 10, 64)

		balance, err := models.GetBalance(params["owner_id"], imonth, iyear)
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{"message": "` + err.Error() + `"}`))
			return
		}

		if balance.ID.Hex() == "" {
			response.WriteHeader(http.StatusNotFound)
			response.Write([]byte(`{}`))
			return
		}

		json.NewEncoder(response).Encode(balance)
	}
}

package controllers

import (
	"budget-tracker-api/crypt"
	"budget-tracker-api/models"
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	"time"

	"github.com/dgrijalva/jwt-go"
)

var mySigninKey = []byte("myhellokey")

// GenerateJWTAccessToken will generate a JWT access token
func GenerateJWTAccessToken(sub string, login string) (string, error) {
	accessToken := jwt.New(jwt.SigningMethodHS256)
	claims := accessToken.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["sub"] = sub
	claims["name"] = login
	claims["exp"] = time.Now().Add(5 * time.Minute).Unix()
	claims["iat"] = time.Now().Unix()

	at, err := accessToken.SignedString(mySigninKey)
	if err != nil {
		return "", err
	}
	return at, nil
}

// GenerateJWTRefreshToken will generate a new refresh token
func GenerateJWTRefreshToken(sub string) (string, error) {
	refreshToken := jwt.New(jwt.SigningMethodHS256)
	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	rtClaims["sub"] = 1
	rtClaims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	rt, err := refreshToken.SignedString(mySigninKey)
	if err != nil {
		return "", err
	}

	return rt, nil
}

// CreateJWTTokenEndpoint creates a token based on user credentials
func CreateJWTTokenEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	response.Header().Set("Access-Control-Allow-Origin", "*")

	var jwtUser models.JWTUser

	_ = json.NewDecoder(request.Body).Decode(&jwtUser)

	if jwtUser.Login == "" || jwtUser.Password == "" {
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte(`{"message": "empty required payload attributes"}`))
		return
	}

	dbUser, err := models.GetUserByFilter(request.Context(), "login", jwtUser.Login)
	if err != nil {
		fmt.Println(err)
		response.WriteHeader(http.StatusUnauthorized)
		response.Write([]byte(`{"message": "invalid credentials for user '` + jwtUser.Login + `'"}`))
		return
	}

	if dbUser.Login == jwtUser.Login {
		// validates password
		match := crypt.CheckPasswordHash(jwtUser.Password, dbUser.SaltedPassword)
		if match {
			AccessToken, err := GenerateJWTAccessToken(dbUser.ID.Hex(), dbUser.Login)
			if err != nil {
				response.WriteHeader(http.StatusInternalServerError)
				response.Write([]byte(`{"message": "could not create access token", "details": "` + err.Error() + `"}`))
				return
			}

			RefreshToken, err := GenerateJWTRefreshToken(dbUser.ID.Hex())
			if err != nil {
				response.WriteHeader(http.StatusInternalServerError)
				response.Write([]byte(`{"message": "could not create refresh token", "details": "` + err.Error() + `"}`))
				return
			}

			log.Infof("created token for user '%s'", jwtUser.Login)
			response.WriteHeader(http.StatusCreated)

			var jwtResponse models.JWTResponse
			jwtResponse.Type = "bearer"
			jwtResponse.RefreshToken = RefreshToken
			jwtResponse.AccessToken = AccessToken
			jwtResponse.Details.ID = dbUser.ID
			jwtResponse.Details.Login = dbUser.Login
			jwtResponse.Details.Firstname = dbUser.Firstname
			jwtResponse.Details.Lastname = dbUser.Lastname
			jwtResponse.Details.Email = dbUser.Email

			jwtResponseJSON, err := json.Marshal(jwtResponse)
			response.Write(jwtResponseJSON)
			if err != nil {
				log.Errorf("could not marshal JWT response for user '%s'", jwtUser.Login)
				response.Write([]byte(`{"message": "could not create access token", "details": "could not marshal JWT response"}`))
			}
			return
		}
	}

	response.WriteHeader(http.StatusUnauthorized)
	response.Write([]byte(`{"message": "invalid credentials for user '` + dbUser.Login + `'"}`))
	return
}

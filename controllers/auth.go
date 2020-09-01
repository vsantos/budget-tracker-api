package controllers

import (
	"budget-tracker/crypt"
	"budget-tracker/models"
	"encoding/json"
	"net/http"
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

	var jwtUser models.JWTUser

	_ = json.NewDecoder(request.Body).Decode(&jwtUser)

	dbUser, err := models.GetUserByFilter("login", jwtUser.Login)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "could not create token", vim"details": "` + err.Error() + `"}`))
		return
	}

	if dbUser.Login == "" || dbUser.SaltedPassword == "" {
		response.WriteHeader(http.StatusNotFound)
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

			response.WriteHeader(http.StatusCreated)
			response.Write([]byte(`{"type": "bearer", "refresh": "` + RefreshToken + `", "token": "` + AccessToken + `"}`))
			return
		}
	}

	response.WriteHeader(http.StatusUnauthorized)
	response.Write([]byte(`{"message": "invalid credentials for user '` + dbUser.Login + `'"}`))
	return
}

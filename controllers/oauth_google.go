package controllers

import (
	"budget-tracker-api/models"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	googleOauthConfig *oauth2.Config
	oauthStateString  = "pseudo-random"
)

func init() {
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:5000/api/v1/oauth/google/callback",
		ClientID:     "963775022969-2o3d2ep205jr5o38iks3qoid9hjm764h.apps.googleusercontent.com",
		ClientSecret: "IXKMusc6paucEwxGHe1zO4fO",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}
}

// GoogleCallbackEndpoint will handle callbacks for google's oauth2
func GoogleCallbackEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	response.Header().Set("Access-Control-Allow-Origin", "*")

	content, err := getGoogleUserInfo(request.FormValue("state"), request.FormValue("code"))
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "error calling callback"}`))
		logrus.Error(err.Error())
		return
	}

	dbUser, err := models.GetUserByFilter(request.Context(), "email", content.Email)
	if err != nil {
		response.WriteHeader(http.StatusUnauthorized)
		response.Write([]byte(`{"message": "invalid credentials for email '` + content.Email + `'"}`))
		return
	}

	// If we do not find an oauth2 user, create it
	// TODO

	// If we find an aouth2 user, just login
	if dbUser.Email == content.Email {
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

		log.Infof("created token for oauth2 email '%s'", content.Email)
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
			log.Errorf("could not marshal JWT response for oauth2 email '%s'", content.Email)
			response.Write([]byte(`{"message": "could not create access token", "details": "could not marshal JWT response"}`))
		}
		return
	}

	response.WriteHeader(http.StatusCreated)
	json.NewEncoder(response).Encode(content)
}

func getGoogleUserInfo(state string, code string) (*models.Oauth2GoogleResponse, error) {
	var oauth2Response models.Oauth2GoogleResponse
	if state != oauthStateString {
		return nil, fmt.Errorf("invalid oauth state")
	}

	token, err := googleOauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, fmt.Errorf("code exchange failed: %s", err.Error())
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}

	defer response.Body.Close()
	err = json.NewDecoder(response.Body).Decode(&oauth2Response)
	if err != nil {
		return nil, err
	}

	return &oauth2Response, nil
}

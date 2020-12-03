package handlers

import (
	"fmt"
	"mime"
	"net/http"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
)

var mySigninKey = []byte("myhellokey")

// Middlewares defines middlewares to intercept handlers
type Middlewares struct {
	Auth func(http.Handler) http.Handler
	JSON func(http.Handler) http.Handler
}

// GetMiddlewares will return all middlewares handlers initialized
func GetMiddlewares() (m Middlewares) {
	m.JSON = RequireContentTypeJSON
	m.Auth = RequireTokenAuthentication
	return m
}

// RequireContentTypeJSON enforces JSON content-type from requests
func RequireContentTypeJSON(h http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		contentType := request.Header.Get("Content-Type")

		if contentType == "" {
			http.Error(response, "Empty Content-Type header", http.StatusBadRequest)
			return
		}
		if contentType != "" {
			mt, _, err := mime.ParseMediaType(contentType)
			if err != nil {
				http.Error(response, "Malformed Content-Type header", http.StatusBadRequest)
				return
			}

			if mt != "application/json" {
				http.Error(response, "Content-Type header must be application/json", http.StatusUnsupportedMediaType)
				return
			}
		}

		h.ServeHTTP(response, request)
	})
}

// RequireTokenAuthentication enforces authentication token from requests
func RequireTokenAuthentication(h http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request.Header["Authorization"] == nil {
			response.WriteHeader(http.StatusBadRequest)
			response.Write([]byte(`{"message": "missing 'Authorization' header"}`))
			return
		}

		if request.Header["Authorization"] != nil {
			jwtString := strings.Split(request.Header["Authorization"][0], "Bearer ")
			if len(jwtString) <= 1 {
				response.WriteHeader(http.StatusUnauthorized)
				response.Write([]byte(`{"message": "could not parse token", "details": "possible mistyped bearer token"}`))
				return
			}
			token, err := jwt.Parse(jwtString[1], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("could not decode token")
				}
				return mySigninKey, nil
			})

			if err != nil {
				response.WriteHeader(http.StatusUnauthorized)
				response.Write([]byte(`{"message": "could not authenticate", "details": "` + err.Error() + `"}`))
				return
			}

			if !token.Valid {
				response.WriteHeader(http.StatusInternalServerError)
				response.Write([]byte(`{"message": "token not valid"}`))
			}
		}

		h.ServeHTTP(response, request)
	})
}

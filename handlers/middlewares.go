package handlers

import (
	"fmt"
	"mime"
	"net/http"
)

// Middlewares defines middlewares to intercept handlers
type Middlewares struct {
	Auth func(http.Handler) http.Handler
	JSON func(http.Handler) http.Handler
}

// GetMiddlewares will return all middlewares handlers initialized
func GetMiddlewares() (m Middlewares) {
	// m.Auth = handlers.Authentication
	m.JSON = ContentTypeJson
	m.Auth = IsAuthenticated
	return m
}

// ContentTypeJson enforces JSON content-type from requests
func ContentTypeJson(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-Type")

		fmt.Println(contentType)
		if contentType == "" {
			http.Error(w, "Empty Content-Type header", http.StatusBadRequest)
			return
		}
		if contentType != "" {
			mt, _, err := mime.ParseMediaType(contentType)
			if err != nil {
				http.Error(w, "Malformed Content-Type header", http.StatusBadRequest)
				return
			}

			if mt != "application/json" {
				http.Error(w, "Content-Type header must be application/json", http.StatusUnsupportedMediaType)
				return
			}
		}

		h.ServeHTTP(w, r)
	})
}

// IsAuthenticated enforces authentication token from requests
func IsAuthenticated(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("auth")

		h.ServeHTTP(w, r)
	})
}

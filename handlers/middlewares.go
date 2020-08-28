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
	m.JSON = RequireContentTypeJSON
	m.Auth = RequireTokenAuthentication
	return m
}

// RequireContentTypeJSON enforces JSON content-type from requests
func RequireContentTypeJSON(h http.Handler) http.Handler {
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

// RequireTokenAuthentication enforces authentication token from requests
func RequireTokenAuthentication(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("auth")

		h.ServeHTTP(w, r)
	})
}

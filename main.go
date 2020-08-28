package main

import (
	"budget-tracker/routes"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

const (
	port = ":8080"
)

var mySigninKey = []byte("myhellokey")

func init() {
	log.SetFormatter(&log.JSONFormatter{})
}

// HomePage b
func HomePage(w http.ResponseWriter, r *http.Request) {
	validToken, err := GenerateJWT()
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	fmt.Fprintf(w, validToken)
	return
}

// GenerateJWT i
func GenerateJWT() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["sub"] = "vsantos"
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(settings.Get().JWTExpirationDelta)).Unix()

	tokenString, err := token.SignedString(mySigninKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func main() {
	log.Infoln("Started Application at port", port)

	router := mux.NewRouter()
	router.HandleFunc("/api/v1/auth/token", HomePage).Methods("POST")
	routes.InitRoutes(router)

	err := http.ListenAndServe(port, router)
	if err != nil {
		log.Errorln(err)
	}
}

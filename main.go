package main

import (
	"budget-tracker/routes"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

const (
	port = ":8080"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
}

func main() {
	log.Infoln("Started Application at port", port)

	router := mux.NewRouter()
	routes.InitRoutes(router)

	err := http.ListenAndServe(port, router)
	if err != nil {
		log.Errorln(err)
	}
}

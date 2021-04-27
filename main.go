// Package classification Budget-tracker API.
//
// the purpose of this application is to provide an application
// that is using plain go code to define an API
//
// This should demonstrate all the possible comment annotations
// that are available to turn go code into a fully compliant swagger 2.0 spec
//
// Terms Of Service:
//
// there are no TOS at this moment, use at your own risk we take no responsibility
//
//     Schemes: http, https
//     Host: budget-tracker:8080
//     BasePath:
//     Version: 0.0.2
//     License: MIT http://opensource.org/licenses/MIT
//     Contact: Victor Santos<vsantos.py@gmail.com> https://github.com/vsantos
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
// swagger:meta
package main

import (
	"budget-tracker/observability"
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

	p, err := observability.JaegerTracerProvider("budget-tracker-api", "http://jaeger:14268/api/traces")
	if err != nil {
		log.Fatal(err)
	}
	observability.InitTrace(p)

	router := mux.NewRouter()
	routes.InitRoutes(router)

	err = http.ListenAndServe(port, router)
	if err != nil {
		log.Errorln(err)
	}
}

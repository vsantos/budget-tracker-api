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
//     Host: budget-tracker:5000
//     BasePath:
//     Version: 0.0.4
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
	"budget-tracker-api/observability"
	"budget-tracker-api/routes"
	"budget-tracker-api/server"
	"budget-tracker-api/services"
	"crypto/tls"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

const (
	port      = ":5000"
	service   = "budget-tracker-api"
	jaegerURL = "http://localhost:14268/api/traces"
	zipkinURL = "http://jaeger:9411/api/v2/spans"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
}

func main() {
	c := observability.ProvidersConfig{
		ServiceName: "budget-tracker-api",
		JaegerURL:   jaegerURL,
		ZipkinURL:   zipkinURL,
	}

	p, err := c.InitTracerProviders()
	if err != nil {
		log.Errorln(err)
	}

	// Change provider exporter if needed. Ex: `p.Stdout`
	observability.InitGlobalTrace(p.Jaeger)
	observability.InitMetrics()

	router := mux.NewRouter()
	routes.InitRoutes(service, router)

	hc := server.HTTPConfig{
		Port:   port,
		Router: router,
		TLSConfig: &tls.Config{
			// In the absence of `NextProtos`, HTTP/1.1 protocol will be enabled
			NextProtos: []string{"h2"},
		},
		CertFile: "config/tls/server.crt",
		KeyFile:  "config/tls/server.key",
	}

	services.MongoClient, err = services.InitDatabaseWithURI(services.MongodbURI)
	if err != nil {
		log.Fatalln(err)
	}

	// repo := repository.NewUserRepository(&repository.UserRepositoryMongoDB{
	// 	Client: services.MongoClient,
	// 	Config: services.MongoCfg{
	// 		URI:       services.MongodbURI,
	// 		Database:  "budget-tracker",
	// 		Colletion: "users",
	// 	},
	// })

	// r, _ := repo.GetAll(context.TODO())
	// fmt.Println(r)

	// In case of 'h2' (HTTP/2) the serverTLS must be set as `true`
	err = hc.InitHTTPServer(false)
	if err != nil {
		log.Fatalln(err)
	}
}

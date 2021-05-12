package server

import (
	"crypto/tls"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

// HTTPConfig will define server configuration for both HTTP/1 and HTTP/2 protocols
type HTTPConfig struct {
	Port      string
	Router    *mux.Router
	TLSConfig *tls.Config
	CertFile  string
	KeyFile   string
}

// InitHTTPServer will init a HTTP/1 (h2) server with optional TLS enforcement
func (c HTTPConfig) InitHTTPServer(serveTLS bool) (err error) {
	c.TLSConfig.NextProtos = []string{"http/1.1"}

	srv := &http.Server{
		Addr:      c.Port,
		Handler:   c.Router,
		TLSConfig: c.TLSConfig,
	}

	if serveTLS {
		log.Infoln(fmt.Sprintf("Started %s Application at port %s with TLS enabled", srv.TLSConfig.NextProtos[0], c.Port))
		err = srv.ListenAndServeTLS(c.CertFile, c.KeyFile)
		if err != nil {
			return err
		}
	}

	if !serveTLS {
		log.Infoln(fmt.Sprintf("Started %s Application at port %s", c.TLSConfig.NextProtos[0], c.Port))
		err = srv.ListenAndServe()
		if err != nil {
			return err
		}
	}

	return nil
}

// InitHTTP2Server will init a HTTP/2 (h2) server with TLS enforcement
func (c HTTPConfig) InitHTTP2Server() (err error) {
	c.TLSConfig.NextProtos = []string{"h2"}

	srv := &http.Server{
		Addr:      c.Port,
		Handler:   c.Router,
		TLSConfig: c.TLSConfig,
	}

	log.Infoln(fmt.Sprintf("Started %s Application at port %s", c.TLSConfig.NextProtos[0], c.Port))
	err = srv.ListenAndServeTLS(c.CertFile, c.KeyFile)
	if err != nil {
		return err
	}

	return nil
}

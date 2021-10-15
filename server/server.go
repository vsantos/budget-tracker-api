package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

// WaitGracefulShutdown is blocking code to wait a SIGNAL to graceful shutdown the server
func WaitGracefulShutdown(srv *http.Server, timeoutSeconds time.Duration) (err error) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	os := <-quit

	log.Printf("Waiting server to shutdown due to signal '%+v'", os)
	ctx, cancel := context.WithTimeout(context.Background(), timeoutSeconds*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}

// InitHTTPServer will init a HTTP/1 (h2) server with optional TLS enforcement
func (c HTTPConfig) InitHTTPServer(serveTLS bool) (err error) {
	srv := &http.Server{
		WriteTimeout: 1 * time.Second,
		ReadTimeout:  1 * time.Second,
		Addr:         c.Port,
		Handler:      c.Router,
		TLSConfig:    c.TLSConfig,
	}

	if serveTLS {
		log.Infoln(fmt.Sprintf("Started %s Application at port %s with TLS enabled", srv.TLSConfig.NextProtos[0], c.Port))
		go func() {
			if err := srv.ListenAndServeTLS(c.CertFile, c.KeyFile); err != nil && err != http.ErrServerClosed {
				log.Panicln("Server error: ", err)
			}
		}()
	}

	if !serveTLS {
		log.Infoln(fmt.Sprintf("Started %s Application at port %s with TLS disabled", srv.TLSConfig.NextProtos[0], c.Port))
		go func() {
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Panicln("Server error: ", err)
			}
		}()
	}

	if err := WaitGracefulShutdown(srv, 15); err != nil {
		return err
	}

	log.Infoln("Server finished")
	return nil
}

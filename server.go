package main

import (
	"github.com/gorilla/mux"

	"fmt"
	"log"
	"net/http"
	"time"
)

func Server(config *Config) *http.Server {

	client := makeClient()
	client.namespace = config.namespace

	r := mux.NewRouter()
	r.HandleFunc("/healthz", healthHandler)

	r.HandleFunc("/full/{fqdn}", client.httpHandler(FqdnHandler))
	r.HandleFunc("/full/{fqdn}/{key}", client.httpHandler(FqdnKeyHandler))
	r.HandleFunc("/full/{fqdn}/{key}/{value}", client.httpHandler(FqdnKeyValueHandler))

	srv := &http.Server{
		Addr: fmt.Sprintf("%s:%d", config.host, config.port),

		WriteTimeout: time.Second * 5,
		ReadTimeout:  time.Second * 5,
		IdleTimeout:  time.Second * 10,

		Handler: r,
	}

	return srv
}

type Application struct {
	config *Config
	server *http.Server
}

func (app *Application) Serve() {
	log.Printf("Serving on %s:%d", app.config.host, app.config.port)
	app.server.ListenAndServe()
}

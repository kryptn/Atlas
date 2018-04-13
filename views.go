package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	ok := struct{ ok bool }{ok: true}
	body, _ := json.Marshal(ok)

	w.Header().Set("content-type", "application/json")
	w.Write(body)
}

type ClientHttpHandlerFunc func(w http.ResponseWriter, r *http.Request, client *Client)

func FqdnHandler(w http.ResponseWriter, r *http.Request, client *Client) {
	vars := mux.Vars(r)

	validFilter := ValidAtlasFilter()
	fqdnFilter := FqdnFilter(vars["fqdn"])
	annotations := client.allAnnotations().Filter(validFilter, fqdnFilter)
	annotations.httpHandlerFunc(w, r)
}

func FqdnKeyHandler(w http.ResponseWriter, r *http.Request, client *Client) {
	vars := mux.Vars(r)

	validFilter := ValidAtlasFilter()
	fqdnFilter := FqdnFilter(vars["fqdn"])
	keyFilter := KeyFilter(vars["key"])
	annotations := client.allAnnotations().Filter(validFilter, fqdnFilter, keyFilter)

	annotations.httpHandlerFunc(w, r)

}

func FqdnKeyValueHandler(w http.ResponseWriter, r *http.Request, client *Client) {
	vars := mux.Vars(r)

	validFilter := ValidAtlasFilter()
	fqdnFilter := FqdnFilter(vars["fqdn"])
	keyValueFilter := KeyValueFilter(vars["key"], vars["value"])
	annotations := client.allAnnotations().Filter(validFilter, fqdnFilter, keyValueFilter)

	annotations.httpHandlerFunc(w, r)
}

func (sas *ServiceAnnotations) httpHandlerFunc(w http.ResponseWriter, r *http.Request) {
	body, _ := json.Marshal(sas.Annotations)
	w.Header().Set("content-type", "application/json")
	w.Write(body)
}

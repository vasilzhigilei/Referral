package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	// Declare a new router
	r := mux.NewRouter()

	// index page handler
	r.HandleFunc("/", indexHandler).Methods("GET")

}

func indexHandler(w http.ResponseWriter, r *http.Request){

}
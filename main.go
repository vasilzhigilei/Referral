package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	// Declare a new router
	r := mux.NewRouter()

	// index page handler
	r.HandleFunc("/", indexHandler).Methods("GET")

	// referral variable paths
	r.HandleFunc("/referral/{service}", serviceHandler)

	// file directory for file serving
	staticFileDirectory := http.Dir("./static/")
	// the prefix is the routing address, the address the user goes to
	staticFileHandler := http.StripPrefix("/static/", http.FileServer(staticFileDirectory))

	// keep PathPrefix empty
	r.PathPrefix("/").Handler(staticFileHandler).Methods("GET")

	http.ListenAndServe(":8000", r)
}

func indexHandler(w http.ResponseWriter, r *http.Request){

}

func serviceHandler(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Category: %v\n", vars["category"])
}
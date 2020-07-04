package main

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
	"net/http"
	"os"
)

var cache redis.Conn
var db *Database

func initCache(){
	conn, err := redis.DialURL(os.Getenv("REDIS_URL"))
	checkErr(err) // check error

	// assign connection to package level 'cache' variable
	cache = conn
}

func initDB() *Database {
	db = NewDatabase(os.Getenv("DATABASE_URL"))
	// err := db.GenerateTable()
	//checkErr(err)
	return db
}

func main() {
	var err error // declare error variable err to avoid :=
	initCache() // initialize redis cache for session/user pairs

	// Declare a new router
	r := mux.NewRouter()

	// index page handler
	r.HandleFunc("/", indexHandler).Methods("GET")

	// categories variable paths
	r.HandleFunc("/categories/{category}", categoryHandler)

	// referral variable paths
	r.HandleFunc("/referrals/{service}", serviceHandler)

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
	// meant to open new tab with referral link. Isn't a separate page, more of an API call
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Service: %v\n", vars["service"])
}

func categoryHandler(w http.ResponseWriter, r *http.Request){
	// meant to route to a category page, listing relevant randomly generated links
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Category: %v\n", vars["category"])
}

/**
Check error func
*/
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
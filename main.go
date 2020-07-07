package main

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
	"html/template"
	"image/color"
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
	//db = NewDatabase(os.Getenv("DATABASE_URL"))
	db = NewDatabase("postgres://postgres:password@localhost:5433/referralshare")
	err := db.GenerateTable()
	checkErr(err)
	return db
}

var indexTemplate *template.Template

func initTemplates() {
	indexTemplate = template.Must(template.ParseFiles("templates/index.html"))
}

type Service struct {
	image string
	background color.RGBA
	backgroundAlt color.RGBA
	description string
}

func main() {
	//var err error // declare error variable err to avoid :=
	initCache() // initialize redis cache for session/user pairs
	db = initDB() // initialize postgres database

	initTemplates()

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
	// main page, meant to open categories, featured categories
	indexTemplate.Execute(w, "")
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
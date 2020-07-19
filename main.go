package main

import (
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
	"html/template"
	"image/color"
	"math/rand"
	"net/http"
	"os"
	"time"
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

var urllists = make(map[string][]string)

func initURLLists() {
	services := []string{"Sofi_money", "Sofi_invest", "Robinhood", "Amazon", "Airbnb", "Grubhub", "Doordash", "Uber"}
	for i := 0; i < len(services); i++ {
		urllists[services[i]] = db.GetServiceURLs(services[i])
	}
}

var indexTemplate *template.Template
var profileTemplate *template.Template

func initTemplates() {
	indexTemplate = template.Must(template.ParseFiles("templates/index.html"))
	profileTemplate = template.Must(template.ParseFiles("templates/profile.html"))
}

type User struct {
	Email string
	Sofi_money string
	Sofi_money_clicks int
	Sofi_invest string
	Sofi_invest_clicks int
	Robinhood string
	Robinhood_clicks int
	Amazon string
	Amazon_clicks int
	Airbnb string
	Airbnb_clicks int
	Grubhub string
	Grubhub_clicks int
	Doordash string
	Doordash_clicks int
	Uber string
	Uber_clicks int
}

type Service struct {
	InternalName string
	Image string
	BackgroundColor color.RGBA
	BorderColor color.RGBA
	ExternalName string
	Clicks int
	Description string
}

func main() {
	//var err error // declare error variable err to avoid :=

	// initialize random generator
	rand.Seed(time.Now().Unix())

	initCache() // initialize redis cache for session/user pairs
	db = initDB() // initialize postgres database
	initURLLists() // initialize urllists map from postgres database on startup

	initTemplates()

	// Declare a new router
	r := mux.NewRouter()

	// index page handler
	r.HandleFunc("/", indexHandler).Methods("GET")

	// categories variable paths
	//r.HandleFunc("/categories/{category}", categoryHandler)

	// profile page handler
	r.HandleFunc("/profile", profileHandler)
	// update form handler
	r.HandleFunc("/updateuser", updateHandler)

	// login/logout management
	r.HandleFunc("/auth/login", loginHandler)
	r.HandleFunc("/auth/callback", callbackHandler)
	r.HandleFunc("/auth/logout", logoutHandler)

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

/**
Check error func
*/
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
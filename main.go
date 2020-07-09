package main

import (
	"encoding/base64"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
	"html/template"
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
	services := []string{"sofi_money", "sofi_invest", "robinhood", "amazon", "airbnb", "grubhub", "doordash", "uber"}
	for i := 0; i < len(services); i++ {
		urllists[services[i]] = db.GetServiceURLs(services[i])
	}
}

var indexTemplate *template.Template

func initTemplates() {
	indexTemplate = template.Must(template.ParseFiles("templates/index.html"))
}

type User struct {
	Email string
	sofi_money string
	sofi_money_clicks int
	sofi_invest string
	sofi_invest_clicks int
	robihood string
	robinhood_clicks int
	amazon string
	amazon_clicks int
	airbnb string
	airbnb_clicks int
	grubhub string
	grubhub_clicks int
	doordash string
	doordash_clicks int
	uber string
	uber_clicks int
}

func main() {
	//var err error // declare error variable err to avoid :=

	// initialize random generator
	rand.Seed(time.Now().Unix())

	//initCache() // initialize redis cache for session/user pairs
	db = initDB() // initialize postgres database
	initURLLists() // initialize urllists map from postgres database on startup

	initTemplates()

	// Declare a new router
	r := mux.NewRouter()

	// index page handler
	r.HandleFunc("/", indexHandler).Methods("GET")

	// categories variable paths
	//r.HandleFunc("/categories/{category}", categoryHandler)

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
	// Opens random referral link for given service
	vars := mux.Vars(r)
	listoflinks := urllists[vars["service"]] // get array of referral links for a given service

	// randomly select a link from the listoflinks string array
	http.Redirect(w, r, listoflinks[rand.Intn(len(listoflinks))], http.StatusTemporaryRedirect)
}

/**
Struct to accept unmarshaling of Google user data
Can be expanded to accept a large variety of additional user information on Google login
Currently only need email address
*/
type GoogleUser struct {
	Email string `json:"email"`
}

// global authentication variable
var authconf = &oauth2.Config {
	RedirectURL: "http://localhost:8000/callback",
	ClientID: os.Getenv("GOOGLE_CLIENT_ID_REFERRALSHARE"),
	ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET_REFERRALSHARE"),
	Scopes: []string{"https://www.googleapis.com/auth/userinfo.email"},
	Endpoint: google.Endpoint,
}

/**
Generates new session with 1 year expiration time
*/
func generateStateOauthCookie(w http.ResponseWriter) string {
	var expiration = time.Now().Add(365 * 24 * time.Hour)

	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}
	http.SetCookie(w, &cookie)

	return state
}

/**
Login handler
Generates random session id, and then redirects client to Google's authentication service
*/
func loginHandler(w http.ResponseWriter, r *http.Request) {
	oauthStateString := generateStateOauthCookie(w)
	url := authconf.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

/**
Logout handler
Handles user logout; Deletes session from Redis cache
*/
func logoutHandler(w http.ResponseWriter, r * http.Request) {
	c, err := r.Cookie("oauthstate")
	checkErr(err)
	_, err = cache.Do("DEL", c.Value)
	checkErr(err)
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

/**
Check error func
*/
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
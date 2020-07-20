package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"time"
)

func indexHandler(w http.ResponseWriter, r *http.Request){
	// main page, meant to open categories, featured categories
	indexTemplate.Execute(w, "")
}

func profileHandler(w http.ResponseWriter, r *http.Request){
	// profile page, if not logged in, auto send to login, if logged in, serve profile template
	c, err := r.Cookie("oauthstate")
	if err != nil {
		// If the session token is not present in cache, set to not logged in
		// For any other type of error, return a bad request status
		if err == http.ErrNoCookie {
			// If the cookie is not set, send to login page
			http.Redirect(w, r, "/auth/login", http.StatusTemporaryRedirect)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response, err := cache.Do("GET", c.Value)
	checkErr(err)
	if response == nil {
		// if session doesn't exist in cache, send to login page
		http.Redirect(w, r, "/auth/login", http.StatusTemporaryRedirect)
		return
	}else {
		fmt.Println(fmt.Sprintf("%s", response), "has loaded profile.html")
		user := db.GetUser(fmt.Sprintf("%s", response))

		profileTemplate.Execute(w, user)
	}
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int){
	w.WriteHeader(status)
	fmt.Fprint(w, status)
}

func serviceHandler(w http.ResponseWriter, r *http.Request){
	// Opens random referral link for given service
	vars := mux.Vars(r)
	listofpairs := urllists[vars["service"]] // get array of referral links for a given service

	if len(listofpairs) == 0 {
		errorHandler(w, r, http.StatusNotFound)
		return
	}

	randompair := listofpairs[rand.Intn(len(listofpairs))]


	// randomly select a link from the listoflinks string array
	http.Redirect(w, r, randompair.URL, http.StatusTemporaryRedirect)
}

func updateHandler(w http.ResponseWriter, r *http.Request){
	c, err := r.Cookie("oauthstate")
	checkErr(err)
	response, err := cache.Do("GET", c.Value)
	checkErr(err)
	if response != nil {
		r.ParseForm()
		user := User{
			Email:              fmt.Sprintf("%s", response),
			Sofi_money:         r.FormValue("Sofi_money"),
			Sofi_invest:		r.FormValue("Sofi_invest"),
			Robinhood:          r.FormValue("Robinhood"),
			Amazon:             r.FormValue("Amazon"),
			Airbnb:             r.FormValue("Airbnb"),
			Grubhub:            r.FormValue("Grubhub"),
			Doordash:           r.FormValue("Doordash"),
			Uber:               r.FormValue("Uber"),
		}

		found, err := regexp.MatchString("^$|(^(https:\\/\\/www\\.)?sofi\\.com\\/invite\\/money\\/\\?gcp=[0-9a-z-]+(\\/)?$)",
			user.Sofi_money)
		if !found{
			http.Error(w, "SoFi Money URL invalid", http.StatusBadRequest)
			return
		}
		found, err = regexp.MatchString("^$|(^(https:\\/\\/www\\.)?sofi\\.com\\/share\\/invest\\/[0-9]+(\\/)?$)",
			user.Sofi_invest)
		if !found{
			http.Error(w, "SoFi Invest URL invalid", http.StatusBadRequest)
			return
		}

		err = db.UpdateUser(&user)
		checkErr(err)
	}
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
	RedirectURL: "http://localhost:8000/auth/callback",
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
	cookie := http.Cookie{Name: "oauthstate", Value: state, Path: "/", Expires: expiration}
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
Callback handler for login
Redirected to by Google's authentication service
Receives session ID and email address, sets session/email pair in cache,
and adds user to Postgres user DB if user doesn't already exist
Redirects to index.html
*/
func callbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")
	token, _ := authconf.Exchange(oauth2.NoContext, code)

	if !token.Valid(){
		fmt.Fprintln(w, "Retrieved invalid token")
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	checkErr(err)

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	checkErr(err)

	var user *GoogleUser
	err = json.Unmarshal(contents, &user)
	checkErr(err)

	state, err := r.Cookie("oauthstate")
	checkErr(err)
	_, err = cache.Do("SETEX", state.Value, 365 * 24 * 60 * 60, user.Email)
	checkErr(err)

	// insert user into postgresql, auto does check if already exists
	err = db.InsertUser(user.Email)
	checkErr(err)

	http.Redirect(w, r, "/profile", http.StatusTemporaryRedirect)
}

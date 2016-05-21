package main

import (
	"net/http"

	"github.com/ChimeraCoder/anaconda"
)

var twitterapi *anaconda.TwitterApi

func init() {
	// Init a single Twitter API
	anaconda.SetConsumerKey(API_KEY)
	anaconda.SetConsumerSecret(API_SECRET)
	twitterapi = anaconda.NewTwitterApi(ACCESS_TOKEN_KEY, ACCESS_TOKEN_SECRET)

	// Create the routes and call handlers
	http.HandleFunc("/", homehandler)
	http.HandleFunc("/task", taskhandler)
	http.HandleFunc("/followers", followerhandler)
	http.HandleFunc("/favicon.ico", nocontenthandler)
	http.HandleFunc("/robots.txt", nocontenthandler)
}

func homehandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "http://www.thedivisionlatino.com/", http.StatusFound)
}

func nocontenthandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
	return
}

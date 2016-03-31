package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/ChimeraCoder/anaconda"
	"google.golang.org/appengine"
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
	http.HandleFunc("/image", imagehandler)
}

func homehandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "http://www.thedivisionlatino.com/", http.StatusFound)
}

func taskhandler(w http.ResponseWriter, r *http.Request) {
	// Create a Google AppEngine Context
	ctx := appengine.NewContext(r)

	// Check if the request was made from Google App Engine's cron
	// service, if not, send unauthorized
	if r.Header.Get("X-Appengine-Cron") != "true" {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	// Get all streamers that are streaming "The Division"
	// and matches our requirements
	streams, err := GetStreamerInfo()

	// Check if the former was possible
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create a cache context
	cache := Instance(ctx)

	// Create a counter
	var count int

	// Create a holder for the last stream
	var holder Stream
	var twerror error

	// Iterate over each one of them finding which one is NOT
	// already in tweet cache
	for _, s := range streams {
		// Check if is in cache
		if cache.IsInStore(s) {
			continue
		}

		// Tweet the stream
		if err := Tweet(s, ctx); err != nil {
			twerror = err
			break
		}

		// Save it to the store
		cache.SaveToStore(s)

		// Increment the counter
		count++

		// Send the given stream to the holder
		holder = *s

		// Break so we can tweet one at a time
		break
	}

	// Check if there was an error
	if twerror != nil {
		http.Error(w, twerror.Error(), http.StatusInternalServerError)
		return
	}

	// If the StreamID has zero-value, that means something is wrong
	if holder.StreamID == 0 {
		fmt.Fprintf(w, "No streams tweeted: none of the streamers had the minimum requirements.")
		return
	}

	fmt.Fprintf(w, "Tweeted stream from %s with %v viewers.", holder.DisplayName, holder.Viewers)
}

func imagehandler(w http.ResponseWriter, r *http.Request) {
	// Get the $GET parameter "url" which contains the image URL
	imageURL := strings.TrimSpace(r.FormValue("url"))

	// Check if it's empty
	if imageURL == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Try converting the image
	image, err := ConvertImage(imageURL)

	// Check if there was an error
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, `<html><body><img src="data:image/png;base64,%s" /></body></html>`, image)
}

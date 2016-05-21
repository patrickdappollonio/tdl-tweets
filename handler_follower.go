package main

import (
	"fmt"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

func followerhandler(w http.ResponseWriter, r *http.Request) {
	// Create a Google AppEngine Context
	ctx := appengine.NewContext(r)

	// Check if the request was made from Google App Engine's cron
	// service, if not, send unauthorized
	// if r.Header.Get("X-Appengine-Cron") != "true" {
	// 	log.Errorf(ctx, "No X-Appengine-Cron header found. Unauthorized!")
	// 	http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	// 	return
	// }

	// Get all follower IDs
	ids, err := GetFollowerIDs(ctx)

	// Check if there was an error processing the list
	if err != nil {
		log.Errorf(ctx, "Error while getting followers IDs: %s", err.Error())
		return
	}

	// Print the user IDs on screen
	for _, v := range ids {
		fmt.Fprintln(w, "Follower ID:", v)
	}
}

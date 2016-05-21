package main

import (
	"fmt"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

func taskhandler(w http.ResponseWriter, r *http.Request) {
	// Create a Google AppEngine Context
	ctx := appengine.NewContext(r)

	// Check if the request was made from Google App Engine's cron
	// service, if not, send unauthorized
	if r.Header.Get("X-Appengine-Cron") != "true" {
		log.Errorf(ctx, "No X-Appengine-Cron header found. Unauthorized!")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	// Get all streamers that are streaming "The Division"
	// and matches our requirements
	streams, err := GetStreamerInfo(ctx)

	// Check if the former was possible
	if err != nil {
		log.Errorf(ctx, err.Error())
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
		if cache.IsStreamInStore(s) {
			continue
		}

		// Tweet the stream
		if err := Tweet(s, ctx); err != nil {
			twerror = err
			break
		}

		// Save it to the store
		cache.SaveStreamToStore(s)

		// Increment the counter
		count++

		// Send the given stream to the holder
		holder = *s

		// Break so we can tweet one at a time
		break
	}

	// Check if there was an error
	if twerror != nil {
		log.Errorf(ctx, err.Error())
		http.Error(w, twerror.Error(), http.StatusInternalServerError)
		return
	}

	// If the StreamID has zero-value, that means something is wrong
	if holder.StreamID == 0 {
		msg := "No streams tweeted: none of the streamers had the minimum requirements."
		log.Debugf(ctx, msg)
		fmt.Fprintf(w, msg)
		return
	}

	msg := fmt.Sprintf("Tweeted stream from %s with %v viewers.", holder.DisplayName, holder.Viewers)
	log.Debugf(ctx, msg)
	fmt.Fprintf(w, msg)
}

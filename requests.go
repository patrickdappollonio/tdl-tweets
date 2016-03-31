package main

import (
	"net/http"

	"golang.org/x/net/context"
	"google.golang.org/appengine/urlfetch"
)

func Request(ctx context.Context, url string) (*http.Response, error) {
	client := urlfetch.Client(ctx)
	return client.Get(url)
}

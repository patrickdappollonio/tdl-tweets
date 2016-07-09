package main

import (
	"fmt"
	"net/http"
	"net/url"

	"golang.org/x/net/context"
	"google.golang.org/appengine/urlfetch"
)

const TWITCH_API_HOST = "api.twitch.tv"

func Request(ctx context.Context, weburl string) (*http.Response, error) {
	u, err := url.Parse(weburl)

	if err != nil {
		return nil, fmt.Errorf("unable to parse url: %v", weburl)
	}

	req, err := http.NewRequest("GET", u.String(), nil)

	if err != nil {
		return nil, fmt.Errorf("unable to construct an HTTP request: %v", err.Error())
	}

	if u.Host == TWITCH_API_HOST {
		req.Header.Add("Client-ID", TWITCH_CLIENT_ID)
	}

	return urlfetch.Client(ctx).Do(req)
}

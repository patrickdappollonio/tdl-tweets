package main

import (
	"encoding/json"

	"golang.org/x/net/context"
)

const (
	TheDivisionURL       = `https://streams.twitch.tv/kraken/streams?limit=10&offset=0&game=Tom+Clancy's+The+Division&broadcaster_language=es&on_site=1`
	MinNumberOfViewers   = 5
	MinNumberOfFollowers = 5
)

// Stream contains all needed values for a given stream
// on "The Division"
type Stream struct {
	StreamID    int64
	Viewers     int64
	Followers   int64
	PreviewURL  string
	DisplayName string
	URL         string
}

// GetStreamerInfo returns a list of "The Division" streams in Spanish
func GetStreamerInfo(ctx context.Context) ([]*Stream, error) {
	// Perform a request to the given Twitch API URL
	res, err := Request(ctx, TheDivisionURL)

	// Check if that was possible
	if err != nil {
		return nil, err
	}

	// Close the body
	defer res.Body.Close()

	// Create a response struct to parse it later
	var response struct {
		Streams []struct {
			StreamID int64 `json:"_id"`
			Viewers  int64 `json:"viewers"`
			Preview  struct {
				URL string `json:"large"`
			} `json:"preview"`
			Channel struct {
				DisplayName string `json:"display_name"`
				Followers   int64  `json:"followers"`
				URL         string `json:"url"`
			} `json:"channel"`
		} `json:"streams"`
	}

	// Decode and check for error
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&response); err != nil {
		return nil, err
	}

	// Iterate and convert what we need
	var streamlist []*Stream
	for _, stream := range response.Streams {
		// Sorry, but we can't tweet new accounts or with little to no
		// followers
		if stream.Channel.Followers < MinNumberOfFollowers {
			continue
		}

		// Sorry, but we can't tweet accounts that have almost no
		// viewers because quality isn't assured
		if stream.Viewers < MinNumberOfViewers {
			continue
		}

		// Append the new stream to the slice
		streamlist = append(streamlist, &Stream{
			StreamID:    stream.StreamID,
			Viewers:     stream.Viewers,
			Followers:   stream.Channel.Followers,
			PreviewURL:  stream.Preview.URL,
			DisplayName: stream.Channel.DisplayName,
			URL:         stream.Channel.URL,
		})
	}

	return streamlist, nil
}

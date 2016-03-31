package main

import (
	"bytes"
	"encoding/base64"

	"github.com/franela/goreq"
)

// Convert converts an image URL to a base64 string.
// Useful to POST to any endpoint
func ConvertImage(image string) (string, error) {
	// Perform the request with plain go request
	resp, err := goreq.Request{Uri: image}.Do()

	// Check if there was an error
	if err != nil {
		return "", err
	}

	// Defer the body close
	defer resp.Body.Close()

	// Read the page body
	var b bytes.Buffer
	_, err = b.ReadFrom(resp.Body)

	// Check if it was possible
	if err != nil {
		return "", err
	}

	// Encode and return
	return base64.StdEncoding.EncodeToString(b.Bytes()), nil
}

package main

import (
	"encoding/base64"
	"net/url"
)

// isValidUrl tests a string to determine if it is a well-structured url or not.
func isValidUrl(toTest string) bool {
	_, err := url.ParseRequestURI(toTest)
	if err != nil {
		return false
	}

	u, err := url.Parse(toTest)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}

	return true
}

func basicAuth(username, password string) string {

	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))

}

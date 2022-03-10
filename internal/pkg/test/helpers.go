package test

import (
	"encoding/json"
	"net/url"
	"testing"
)

// StrToURL - convert string slice to urls slice
func StrToURL(t *testing.T, strURLs []string) []*url.URL {
	t.Helper()
	urls := []*url.URL{}
	for _, v := range strURLs {
		parsedURL, err := url.Parse(v)
		if err != nil {
			t.Error(err)
		}
		urls = append(urls, parsedURL)
	}
	return urls
}

// ToJSON - json test helper method
func ToJSON(t *testing.T, in interface{}) string {
	t.Helper()
	b, err := json.Marshal(in)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

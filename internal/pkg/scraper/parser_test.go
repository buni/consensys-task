package scraper

import (
	"net/url"
	"os"
	"testing"

	"golang.org/x/net/html"
)

func TestParseHTMLLinks(t *testing.T) {
	tests := []struct {
		name         string
		wantExternal uint
		wantInternal uint
		wantErr      bool
		url          string
		setup        func(t *testing.T, url string) (*url.URL, *html.Node)
	}{
		{
			name:         "successfully parse html",
			url:          "http://localhost.com/",
			wantExternal: 2,
			wantInternal: 6,
			setup: func(t *testing.T, baseURL string) (*url.URL, *html.Node) {
				parsedBaseURL, err := url.Parse(baseURL)
				if err != nil {
					t.Error(err)
				}

				f, err := os.Open("testdata/good_links.html")
				if err != nil {
					t.Error(err)
				}

				defer f.Close()
				document, err := html.Parse(f)
				if err != nil {
					t.Error(err)
				}

				return parsedBaseURL, document
			},
		},
		{
			name:         "successfully parse html with bad links",
			url:          "http://localhost/",
			wantExternal: 1,
			wantInternal: 1,
			setup: func(t *testing.T, baseURL string) (*url.URL, *html.Node) {
				parsedBaseURL, err := url.Parse(baseURL)
				if err != nil {
					t.Error(err)
				}
				f, err := os.Open("testdata/bad_links.html")
				if err != nil {
					t.Error(err)
				}
				defer f.Close()
				document, err := html.Parse(f)
				if err != nil {
					t.Error(err)
				}

				return parsedBaseURL, document
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testURL, document := tt.setup(t, tt.url)
			gotExternal, gotInternal, err := ParseHTMLLinks(testURL, document)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseHTMLLinks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotExternal != tt.wantExternal {
				t.Errorf("ParseHTMLLinks() gotExternal = %v, want %v", gotExternal, tt.wantExternal)
			}
			if gotInternal != tt.wantInternal {
				t.Errorf("ParseHTMLLinks() gotInternal = %v, want %v", gotInternal, tt.wantInternal)
			}
		})
	}
}

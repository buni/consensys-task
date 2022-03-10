package urls

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/url"
)

var ErrInvalidURL = errors.New("invalid url")

// ParseURLs - parse a file/readcloser containing line delimited URLs
// each line has to have a valid url other wise an error is returned
// empty lines are also not permitted
func ParseURLs(r io.ReadCloser) ([]*url.URL, error) {
	urls := make([]*url.URL, 0, 64)

	scanner := bufio.NewScanner(r)
	line := 1

	for scanner.Scan() {
		parsedURL, err := url.Parse(scanner.Text())
		if err != nil {
			return nil, fmt.Errorf("failed to parse url on line %v: %w", line, err)
		}

		if parsedURL.Scheme == "" || parsedURL.Host == "" { // url.Parse doesn't always return an error so some extra checks are needed
			return nil, fmt.Errorf("invalid url on line %v %s %w", line, parsedURL, ErrInvalidURL)
		}

		urls = append(urls, parsedURL)

		line++
	}

	return urls, nil
}

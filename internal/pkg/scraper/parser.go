package scraper

import (
	"log"
	"net/url"

	"golang.org/x/net/html"
)

// ParseHTMLLinks extracts external & internal links from html document
func ParseHTMLLinks(page *url.URL, document *html.Node) (external, internal uint, err error) {
	var f func(*html.Node)

	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" { // ignore link/img/script tags
			for _, attr := range n.Attr {
				if attr.Key == "href" && attr.Val != "" { // links with empty href value are not counted as internal (or at all)
					hrefURL, err := url.Parse(attr.Val)
					if err != nil { // its acceptable to skip this error and continue proccessing the next node
						log.Println("malformed href value: ", attr.Val, err)
						break // assuming there is only one href attribute per node, we can break from the loop
					}

					switch {
					case hrefURL.Hostname() == page.Hostname(): // hostnames match (sub domains are treated as external links)
						internal++
					case hrefURL.Hostname() == "" && hrefURL.Path != "": // if the host is not set but path is set the link most likely is internal
						internal++
					default: // everything else is external
						external++
					}

				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(document)
	return
}

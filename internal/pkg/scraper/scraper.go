package scraper

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sync"

	"github.com/hashicorp/go-cleanhttp"
	"golang.org/x/net/html"
)

var (
	ErrBadStatusCode       = errors.New("bad status code")
	ErrCloseTimeout        = errors.New("close took longer than deadline")
	ErrBadConcurrencyValue = errors.New("bad concurrency value")
)

type Scraper struct {
	httpClient        *http.Client
	wg                *sync.WaitGroup
	produceConcurency int
}

//go:generate mockgen -source=scraper.go -destination=mock/scraper_mocks.go -package mock

// ScraperService ...
type ScraperService interface {
	ScrapePages(ctx context.Context, urls []*url.URL, reqOptions ...ScrapeRequestOption) []Result
	Close(ctx context.Context) error
}

type ScraperOption func(s *Scraper) error

// WithConcurrency sets the maximum concurency for the scrape pipeline
func WithConcurrency(concurency int) ScraperOption {
	return func(s *Scraper) error {
		if concurency <= 0 {
			return ErrBadConcurrencyValue
		}
		s.produceConcurency = concurency
		return nil
	}
}

// ScrapeRequestOption modify http request used for scrape
type ScrapeRequestOption func(r *http.Request) error

// NewScraper ...
func NewScraper(options ...ScraperOption) (*Scraper, error) {
	scraper := &Scraper{}
	scraper.wg = &sync.WaitGroup{}
	scraper.httpClient = cleanhttp.DefaultClient()
	scraper.produceConcurency = 1000

	for _, option := range options {
		err := option(scraper)
		if err != nil {
			return nil, fmt.Errorf("failed to apply parser option %w", err)
		}
	}

	return scraper, nil
}

// ScrapePages - scrapes the provided urls
func (p *Scraper) ScrapePages(ctx context.Context, urls []*url.URL, reqOptions ...ScrapeRequestOption) []Result {
	wg := &sync.WaitGroup{}
	resultsChan := make(chan Result)
	urlsChan := make(chan *url.URL)

	p.wg.Add(1)
	defer p.wg.Done()

	results := make([]Result, 0, len(urls))
	p.enqueueURLs(ctx, wg, urlsChan, urls...)
	p.produceResults(ctx, wg, urlsChan, resultsChan)

	for result := range resultsChan {
		results = append(results, result)
	}

	return results
}

func (p *Scraper) enqueueURLs(ctx context.Context, wg *sync.WaitGroup, urlsChan chan *url.URL, urls ...*url.URL) {
	wg.Add(1)
	go func() {
		defer wg.Done()

		for _, v := range urls {
			select {
			case <-ctx.Done():
				close(urlsChan)
				return
			case urlsChan <- v:
			}
		}

		close(urlsChan)
	}()
}

func (p *Scraper) produceResults(ctx context.Context, wg *sync.WaitGroup, urlsChan chan *url.URL, resultsChan chan Result, reqOptions ...ScrapeRequestOption) {
	for i := 0; i < p.produceConcurency; i++ { // TODO: configure concurency
		wg.Add(1)
		go func() {
			defer wg.Done()
			for u := range urlsChan {
				select {
				case <-ctx.Done():
					return
				case resultsChan <- p.scrapePage(ctx, u, reqOptions...):
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(resultsChan)
	}()
}

func (p *Scraper) scrapePage(ctx context.Context, page *url.URL, reqOptions ...ScrapeRequestOption) Result { // TODO: add options to above methods
	result := Result{PageURL: page.String()}

	req, err := http.NewRequestWithContext(ctx, "GET", page.String(), nil)
	if err != nil {
		result.Success = false
		result.Error = err
		return result
	}
	// TODO: add content type html header
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:97.0) Gecko/20100101 Firefox/97.0")

	for _, option := range reqOptions {
		err = option(req)
		if err != nil {
			result.Error = err
			return result
		}
	}

	resp, err := p.httpClient.Do(req)
	if err != nil {
		result.Error = err
		return result
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 { // TODO:
		result.Error = ErrBadStatusCode
		return result
	}

	document, err := html.Parse(resp.Body)
	if err != nil {
		result.Error = err
		return result
	}

	external, internal, err := ParseHTMLLinks(page, document)
	if err != nil {
		result.Error = err
		result.ExternalLinksCount = external
		result.InternalLinksCount = internal
		return result
	}

	result.ExternalLinksCount = external
	result.InternalLinksCount = internal
	result.Success = true

	return result
}

// Close - waits for all pipelines to end
// passing a context with deadline/cancel, and fullfiling cancel conditions
// will make the method exit early
func (s *Scraper) Close(ctx context.Context) error {
	done := make(chan bool)
	go func() {
		s.wg.Wait()
		close(done)
	}()

	for {
		select {
		case <-ctx.Done():
			return ErrCloseTimeout
		case <-done:
			return nil
		}
	}
}

package scraper

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/buni/scraper/internal/pkg/test"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func testServerHelper(t *testing.T, filePath string) string {
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		body, err := os.ReadFile(filePath)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(body)
	})
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatal(err)
	}
	srv := &http.Server{Handler: r}

	go func() {
		srv.Serve(listener)
	}()

	t.Cleanup(func() {
		srv.Shutdown(context.Background())
	})

	return strings.Replace(listener.Addr().String(), "127.0.0.1", "localhost", -1)
}

func TestScraper_ScrapePages(t *testing.T) {
	tests := []struct {
		name    string
		ctx     context.Context
		setup   func(t *testing.T) (s *Scraper, hostURLs []*url.URL, opts []ScrapeRequestOption)
		wantLen int
		wantErr bool
	}{
		{
			name: "successfully scrape pages",
			ctx:  context.Background(),
			setup: func(t *testing.T) (s *Scraper, hostURLs []*url.URL, opts []ScrapeRequestOption) {
				s, err := NewScraper()
				assert.NoError(t, err)

				host := "http://" + testServerHelper(t, "testdata/good_links_serve.html") + "/"
				opts = []ScrapeRequestOption{}
				hostURL, err := url.Parse(host)
				hostURLs = append(hostURLs, hostURL, hostURL)
				assert.NoError(t, err)

				return
			},
			wantLen: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, host, opts := tt.setup(t)
			got := s.ScrapePages(tt.ctx, host, opts...)
			t.Log(got)
			assert.Len(t, got, tt.wantLen)
		})
	}
}

func TestScraper_enqueueURLs(t *testing.T) {
	t.Run("successfully enqueue urls", func(t *testing.T) {
		wg := &sync.WaitGroup{}
		urlsChan := make(chan *url.URL)
		urls := test.StrToURL(t, []string{
			"https://locahost/",
			"https://locahost/",
			"http://locahost/",
			"http://locahost/",
		})

		s, err := NewScraper()
		assert.NoError(t, err)
		s.enqueueURLs(context.Background(), wg, urlsChan, urls...)
		consumedURLs := 0
		for range urlsChan {
			consumedURLs++
		}
		wg.Wait() // makes sure that Done was called
		assert.Equal(t, len(urls), consumedURLs)
	})
	t.Run("cancel enqueue urls", func(t *testing.T) {
		wg := &sync.WaitGroup{}
		urlsChan := make(chan *url.URL)
		urls := test.StrToURL(t, []string{
			"https://locahost/",
			"https://locahost/",
			"http://locahost/",
			"http://locahost/",
		})

		s, err := NewScraper()
		assert.NoError(t, err)
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, time.Millisecond)
		defer cancel()
		s.enqueueURLs(ctx, wg, urlsChan, urls...)

		wg.Wait() // makes sure that Done was called
	})
}

func TestScraper_produceResults(t *testing.T) {
	t.Run("successfully produce results", func(t *testing.T) {
		wg := &sync.WaitGroup{}
		urlsChan := make(chan *url.URL)
		resultsChan := make(chan Result)
		urls := test.StrToURL(t, []string{
			"https://localhost/",
			"https://localhost/",
			"http://localhost/",
			"http://localhost/",
		})

		s, err := NewScraper()
		assert.NoError(t, err)
		s.enqueueURLs(context.Background(), wg, urlsChan, urls...)
		s.httpClient.Timeout = time.Millisecond * 100
		s.produceConcurency = 10
		s.produceResults(context.Background(), wg, urlsChan, resultsChan)
		producedResults := 0
		for range resultsChan {
			producedResults++
		}
		wg.Wait() // makes sure that Done was called
		assert.Equal(t, len(urls), producedResults)
	})
	t.Run("cancel produce results", func(t *testing.T) {
		wg := &sync.WaitGroup{}
		urlsChan := make(chan *url.URL)
		resultsChan := make(chan Result)
		urls := test.StrToURL(t, []string{
			"https://localhost/",
			"https://localhost/",
			"http://localhost/",
			"http://localhost/",
		})

		s, err := NewScraper()
		assert.NoError(t, err)
		s.enqueueURLs(context.Background(), wg, urlsChan, urls...)
		s.httpClient.Timeout = time.Millisecond * 100
		s.produceConcurency = 10
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, time.Millisecond)
		defer cancel()
		s.produceResults(ctx, wg, urlsChan, resultsChan)

		wg.Wait() // makes sure that Done was called
	})
}

func TestScraper_scrapePage(t *testing.T) {
	tests := []struct {
		name    string
		ctx     context.Context
		setup   func(t *testing.T) (s *Scraper, hostURL *url.URL, opts []ScrapeRequestOption)
		want    Result
		wantErr bool
	}{
		{
			name: "successfully scrape page",
			ctx:  context.Background(),
			setup: func(t *testing.T) (s *Scraper, hostURL *url.URL, opts []ScrapeRequestOption) {
				s, err := NewScraper()
				assert.NoError(t, err)

				host := "http://" + testServerHelper(t, "testdata/good_links_serve.html") + "/"
				opts = []ScrapeRequestOption{}
				hostURL, err = url.Parse(host)
				assert.NoError(t, err)

				return
			},
			want: Result{
				InternalLinksCount: 6,
				ExternalLinksCount: 2,
				Success:            true,
			},
		},
		{
			name: "fail to scrape bad url",
			ctx:  context.Background(),
			setup: func(t *testing.T) (s *Scraper, hostURL *url.URL, opts []ScrapeRequestOption) {
				s, err := NewScraper()
				assert.NoError(t, err)
				opts = []ScrapeRequestOption{}
				hostURL = &url.URL{Scheme: "bad", Path: `\\\\\\\\\\\\\\\\\`}
				return
			},
			want: Result{
				Success: false,
			},
			wantErr: true,
		},
		{
			name: "bad scrape option",
			ctx:  context.Background(),
			setup: func(t *testing.T) (s *Scraper, hostURL *url.URL, opts []ScrapeRequestOption) {
				s, err := NewScraper()
				assert.NoError(t, err)
				opts = []ScrapeRequestOption{
					func(r *http.Request) error {
						return errors.New("some error")
					},
				}
				hostURL = &url.URL{}
				return
			},
			want: Result{
				Success: false,
			},
			wantErr: true,
		},
		{
			name: "fail to scrape bad status code",
			ctx:  context.Background(),
			setup: func(t *testing.T) (s *Scraper, hostURL *url.URL, opts []ScrapeRequestOption) {
				s, err := NewScraper()
				assert.NoError(t, err)
				host := "http://" + testServerHelper(t, "testdata/non_existant_file") + "/"
				opts = []ScrapeRequestOption{}
				hostURL, err = url.Parse(host)
				assert.NoError(t, err)

				return
			},
			want: Result{
				Success: false,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, host, opts := tt.setup(t)

			tt.want.PageURL = host.String()
			got := s.scrapePage(tt.ctx, host, opts...)
			if tt.wantErr {
				t.Log(tt.wantErr, got.Error)
				assert.Error(t, got.Error)
			} else {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestScraper_Close(t *testing.T) {
	t.Run("successfully close scraper", func(t *testing.T) {
		s, err := NewScraper()
		assert.NoError(t, err)
		err = s.Close(context.Background())
		assert.NoError(t, err)
	})
	t.Run("successfully close scraper with work", func(t *testing.T) {
		s, err := NewScraper()
		assert.NoError(t, err)
		s.wg.Add(1)
		go func() {
			time.Sleep(time.Millisecond * 10)
			s.wg.Done()
		}()
		err = s.Close(context.Background())
		assert.NoError(t, err)
		s.wg.Wait()
	})
	t.Run("successfully close scraper before context timeout", func(t *testing.T) {
		s, err := NewScraper()
		assert.NoError(t, err)
		s.wg.Add(1)
		go func() {
			time.Sleep(time.Millisecond * 10)
			s.wg.Done()
		}()
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, time.Millisecond*20)
		defer cancel()
		err = s.Close(ctx)
		assert.NoError(t, err)
		s.wg.Wait()
	})
	t.Run("fail to close scraper before context timeout", func(t *testing.T) {
		s, err := NewScraper()
		assert.NoError(t, err)
		s.wg.Add(1)
		go func() {
			time.Sleep(time.Millisecond * 100)
			s.wg.Done()
		}()
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, time.Millisecond*20)
		defer cancel()
		err = s.Close(ctx)
		assert.Error(t, err)
		s.wg.Wait()
	})
}

func TestWithConcurrency(t *testing.T) {
	type args struct {
		concurency int
	}
	tests := []struct {
		name string
		args args
		want ScraperOption
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithConcurrency(tt.args.concurency); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithConcurrency() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewScraper(t *testing.T) {
	type args struct {
		options []ScraperOption
	}
	tests := []struct {
		name    string
		args    args
		want    *Scraper
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewScraper(tt.args.options...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewScraper() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewScraper() = %v, want %v", got, tt.want)
			}
		})
	}
}

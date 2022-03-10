package handler

import (
	"bytes"
	"context"
	"io"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/buni/scraper/internal/api/links"
	"github.com/buni/scraper/internal/api/links/mock"
	"github.com/buni/scraper/internal/api/links/repository"
	"github.com/buni/scraper/internal/api/links/service"
	"github.com/kinbiko/jsonassert"

	"github.com/buni/scraper/internal/pkg/scraper"
	"github.com/buni/scraper/internal/pkg/test"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestHandler_EnqueueLinksJobIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E tests")
	}
	t.Parallel()
	r := chi.NewRouter()
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

	repo := repository.NewInMemoryRepository()
	scraperSvc, _ := scraper.NewScraper()
	svc := service.NewService(repo, scraperSvc)

	h := NewHandler(svc)
	h.RegisterRoutes(r)
	tests := []struct {
		name         string
		statusCode   int
		responseBody links.Response
		request      []string
		setup        func(*mock.MockService)
	}{
		{
			name:       "successfully enqueue job",
			statusCode: http.StatusAccepted,
			request: []string{
				"https://localhost",
				"http://localhost",
			},
			responseBody: links.Response{
				Data: links.EnqueueLinksJobResponse{
					JobID: "<<PRESENCE>>",
				},
			},
			setup: func(ms *mock.MockService) {
				ms.EXPECT().EnqueueLinksJob(gomock.Any(), gomock.Any()).Return(
					links.Job{
						ID: uuid.Nil.String(),
					}, nil,
				)
			},
		},
		{
			name: "error no request body",

			statusCode: http.StatusBadRequest,
			responseBody: links.Response{
				Errors: []string{
					links.ErrEmptyJobRequest.Error(),
				},
			},
			setup: func(ms *mock.MockService) {
			},
		},
		{
			name: "error no request body",
			request: []string{
				"j�^>���p",
				"http://localhost",
			},
			statusCode: http.StatusBadRequest,
			responseBody: links.Response{
				Errors: []string{
					"failed to parse url on line 1: parse \"\\x02j�^>�\\x04��p\": net/url: invalid control character in URL",
				},
			},
			setup: func(ms *mock.MockService) {
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := &bytes.Buffer{}
			for _, v := range tt.request {
				body.WriteString(v + "\n")
			}
			req, err := http.NewRequest("POST", "http://"+listener.Addr().String()+"/links", body)
			assert.NoError(t, err)
			defer req.Body.Close()

			resp, err := http.DefaultClient.Do(req)
			assert.NoError(t, err)

			defer resp.Body.Close()

			assert.Equal(t, tt.statusCode, resp.StatusCode)
			jsonBody, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)
			ja := jsonassert.New(t)
			ja.Assertf(string(jsonBody), test.ToJSON(t, tt.responseBody))
		})
	}
}

func TestHandler_GetJobStatusIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E tests")
	}
	t.Parallel()
	r := chi.NewRouter()
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatal(err)
	}
	srv := &http.Server{Handler: r}

	repo := repository.NewInMemoryRepository()
	scraperSvc, _ := scraper.NewScraper()
	svc := service.NewService(repo, scraperSvc)
	h := NewHandler(svc)
	h.RegisterRoutes(r)

	go func() {
		srv.Serve(listener)
	}()

	t.Cleanup(func() {
		srv.Shutdown(context.Background())
	})
	chi.Walk(r, func(method, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		t.Log(method, route)
		return nil
	})
	epoch := time.Unix(0, 0)

	tests := []struct {
		name         string
		statusCode   int
		responseBody links.Response
		setup        func(t *testing.T, repo links.Repository) string
	}{
		{
			name:       "successfully get job status",
			statusCode: 200,
			responseBody: links.Response{
				Data: links.JobResultsResponse{
					Results: []links.JobResult{
						{
							ID:        uuid.Nil.String(),
							JobID:     uuid.Nil.String(),
							PageURL:   "http://localhost",
							CreatedAt: epoch,
							UpdatedAt: epoch,
						},
					},
				},
			},
			setup: func(t *testing.T, repo links.Repository) string {
				_, err := repo.CreateLinksJob(context.Background(), links.Job{
					ID:         uuid.Nil.String(),
					FinishedAt: &epoch,
				})
				assert.NoError(t, err)
				err = repo.CreateLinksJobResult(context.Background(), []links.JobResult{
					{
						ID:        uuid.Nil.String(),
						JobID:     uuid.Nil.String(),
						PageURL:   "http://localhost",
						CreatedAt: epoch,
						UpdatedAt: epoch,
					},
				})
				assert.NoError(t, err)
				return uuid.Nil.String()
			},
		},
		{
			name:       "job not found",
			statusCode: 404,
			responseBody: links.Response{
				Errors: []string{
					repository.ErrJobNotFound.Error(),
				},
			},
			setup: func(t *testing.T, repo links.Repository) string {
				return uuid.NewString()
			},
		},
		{
			name:         "results not ready",
			statusCode:   202,
			responseBody: links.Response{},
			setup: func(t *testing.T, repo links.Repository) string {
				id := uuid.NewString()
				_, err := repo.CreateLinksJob(context.Background(), links.Job{
					ID:         id,
					FinishedAt: &epoch,
				})
				assert.NoError(t, err)
				return id
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jobID := tt.setup(t, repo)

			req, err := http.NewRequest("GET", "http://"+listener.Addr().String()+"/links/status/"+jobID, nil)
			assert.NoError(t, err)
			resp, err := http.DefaultClient.Do(req)
			assert.NoError(t, err)

			defer resp.Body.Close()

			assert.Equal(t, tt.statusCode, resp.StatusCode)
			jsonBody, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)
			ja := jsonassert.New(t)
			ja.Assertf(string(jsonBody), test.ToJSON(t, tt.responseBody))
		})
	}
}

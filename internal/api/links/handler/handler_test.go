package handler

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/buni/scraper/internal/api/links"
	"github.com/buni/scraper/internal/api/links/mock"
	"github.com/buni/scraper/internal/api/links/repository"
	"github.com/buni/scraper/internal/pkg/test"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestHandler_EnqueueLinksJob(t *testing.T) {
	t.Parallel()
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
					JobID: uuid.Nil.String(),
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
			name: "job  already exists",
			request: []string{
				"https://localhost",
				"http://localhost",
			},
			statusCode: http.StatusConflict,
			responseBody: links.Response{
				Errors: []string{
					repository.ErrJobAlreadyExists.Error(),
				},
			},
			setup: func(ms *mock.MockService) {
				ms.EXPECT().EnqueueLinksJob(gomock.Any(), gomock.Any()).Return(links.Job{}, repository.ErrJobAlreadyExists)
			},
		},
		{
			name: "internal error",
			request: []string{
				"https://localhost",
				"http://localhost",
			},
			statusCode: http.StatusInternalServerError,
			responseBody: links.Response{
				Errors: []string{
					links.ErrInternalServerError.Error(),
				},
			},
			setup: func(ms *mock.MockService) {
				ms.EXPECT().EnqueueLinksJob(gomock.Any(), gomock.Any()).Return(links.Job{}, errors.New("some internal error"))
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
			ctrl := gomock.NewController(t)
			service := mock.NewMockService(ctrl)
			h := NewHandler(service)
			tt.setup(service)
			body := &bytes.Buffer{}
			for _, v := range tt.request {
				body.WriteString(v + "\n")
			}
			req, err := http.NewRequest("GET", "/", body)
			assert.NoError(t, err)
			recorder := httptest.NewRecorder()
			h.EnqueueLinksJob(recorder, req)
			assert.Equal(t, tt.statusCode, recorder.Code)

			assert.JSONEq(t, test.ToJSON(t, tt.responseBody), recorder.Body.String())
		})
	}
}

func TestHandler_GetJobStatus(t *testing.T) {
	t.Parallel()
	epoch := time.Unix(0, 0)
	tests := []struct {
		name         string
		statusCode   int
		responseBody links.Response
		setup        func(*mock.MockService)
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
			setup: func(ms *mock.MockService) {
				ms.EXPECT().GetLinksJobStatus(gomock.Any(), gomock.Any()).Return(
					[]links.JobResult{
						{
							ID:        uuid.Nil.String(),
							JobID:     uuid.Nil.String(),
							PageURL:   "http://localhost",
							CreatedAt: epoch,
							UpdatedAt: epoch,
						},
					}, nil,
				)
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
			setup: func(ms *mock.MockService) {
				ms.EXPECT().GetLinksJobStatus(gomock.Any(), gomock.Any()).Return(nil, repository.ErrJobNotFound)
			},
		},
		{
			name:         "results not ready",
			statusCode:   202,
			responseBody: links.Response{},
			setup: func(ms *mock.MockService) {
				ms.EXPECT().GetLinksJobStatus(gomock.Any(), gomock.Any()).Return(nil, repository.ErrJobResultsNotFound)
			},
		},
		{
			name:       "internal error",
			statusCode: 500,
			responseBody: links.Response{
				Errors: []string{
					links.ErrInternalServerError.Error(),
				},
			},
			setup: func(ms *mock.MockService) {
				ms.EXPECT().GetLinksJobStatus(gomock.Any(), gomock.Any()).Return(nil, errors.New("some internal error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			service := mock.NewMockService(ctrl)
			h := NewHandler(service)
			tt.setup(service)
			req, err := http.NewRequest("GET", "/", nil)
			assert.NoError(t, err)
			recorder := httptest.NewRecorder()
			h.GetJobStatus(recorder, req)
			assert.Equal(t, tt.statusCode, recorder.Code)

			assert.JSONEq(t, test.ToJSON(t, tt.responseBody), recorder.Body.String())
		})
	}
}

package links

import (
	"errors"
	"net/url"
	"time"
)

var (
	ErrInternalServerError = errors.New("internal.server.error")
	ErrEmptyJobRequest     = errors.New("empty job request")
)

// EnqueueLinksJobRequest ...
type EnqueueLinksJobRequest struct {
	JobID string
	URLs  []*url.URL
}

// Response - generic http response structure
type Response struct {
	Errors []string    `json:"errors,omitempty"`
	Data   interface{} `json:"data,omitempty"`
}

// EnqueueLinksJobResponse ...
type EnqueueLinksJobResponse struct {
	JobID string `json:"job_id"`
}

// JobResultsResponse ...
type JobResultsResponse struct {
	Results []JobResult `json:"results"` // FIXME: don't reuse the "model" in the response
}

// GetJobStatusRequest ...
type GetJobStatusRequest struct {
	JobID string
}

// Job model
type Job struct {
	ID         string
	URLs       []*url.URL
	CreatedAt  time.Time
	UpdatedAt  time.Time
	FinishedAt *time.Time
}

// JobResult model
type JobResult struct {
	ID                 string    `json:"id"`
	JobID              string    `json:"job_id"`
	PageURL            string    `json:"page_url"`
	InternalLinksCount uint      `json:"internal_links_count"`
	ExternalLinksCount uint      `json:"external_links_count"`
	Success            bool      `json:"success"`
	Error              error     `json:"error"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

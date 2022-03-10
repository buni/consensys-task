package repository

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/buni/scraper/internal/api/links"
	"github.com/google/uuid"
)

var (
	ErrJobAlreadyExists   = errors.New("job already exists")
	ErrJobNotFound        = errors.New("job not found")
	ErrJobResultsNotFound = errors.New("job results found")
)

type inMemRepository struct {
	jobs       map[string]links.Job
	jobResults map[string][]links.JobResult
	rw         *sync.RWMutex
}

func NewInMemoryRepository() links.Repository {
	return &inMemRepository{jobs: map[string]links.Job{}, jobResults: map[string][]links.JobResult{}, rw: &sync.RWMutex{}}
}

// CreateLinksJob - creates new links job
// if job id exists returns ErrJobAlreadyExists
func (r *inMemRepository) CreateLinksJob(ctx context.Context, job links.Job) (links.Job, error) {
	r.rw.Lock()
	defer r.rw.Unlock()

	if job.CreatedAt.IsZero() {
		job.CreatedAt = time.Now().UTC()
	}

	if job.UpdatedAt.IsZero() {
		job.UpdatedAt = time.Now().UTC()
	}

	if job.ID == "" {
		job.ID = uuid.NewString()
	}

	if _, ok := r.jobs[job.ID]; ok {
		return links.Job{}, ErrJobAlreadyExists
	}

	r.jobs[job.ID] = job

	return job, nil
}

// GetLinksJob - get links job by id
// if it doesn't exists returns  ErrJobNotFound
func (r *inMemRepository) GetLinksJob(ctx context.Context, jobID string) (links.Job, error) {
	r.rw.Lock()
	defer r.rw.Unlock()

	job, ok := r.jobs[jobID]
	if !ok {
		return links.Job{}, ErrJobNotFound
	}

	return job, nil
}

// FinishLinksJob - mark links job as finished
func (r *inMemRepository) FinishLinksJob(ctx context.Context, jobID string) error {
	r.rw.Lock()
	defer r.rw.Unlock()

	job, ok := r.jobs[jobID]
	if !ok {
		return ErrJobNotFound
	}

	finishedAt := time.Now().UTC()

	job.FinishedAt = &finishedAt
	r.jobs[jobID] = job

	return nil
}

// CreateLinksJobResult - create links job result
func (r *inMemRepository) CreateLinksJobResult(ctx context.Context, results []links.JobResult) error {
	r.rw.Lock()
	defer r.rw.Unlock()

	if len(results) == 0 {
		return nil
	}

	r.jobResults[results[0].JobID] = results // assume all results have the same job id and job already exists

	return nil
}

// GetLinksJobResult - get links job by job id
// if it doesn't exists an error is returned
func (r *inMemRepository) GetLinksJobResult(ctx context.Context, jobID string) ([]links.JobResult, error) {
	r.rw.RLock()
	defer r.rw.RUnlock()

	results, ok := r.jobResults[jobID]
	if !ok {
		return nil, ErrJobResultsNotFound
	}

	return results, nil
}

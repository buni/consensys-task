package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/buni/scraper/internal/api/links"
	"github.com/buni/scraper/internal/pkg/scraper"
	"github.com/google/uuid"
)

type service struct {
	scraperClient scraper.ScraperService
	repository    links.Repository
}

func NewService(repository links.Repository, scraperClient scraper.ScraperService) links.Service {
	return &service{scraperClient: scraperClient, repository: repository}
}

// EnqueueLinksJob - create links job and start executing it
func (s *service) EnqueueLinksJob(ctx context.Context, req links.EnqueueLinksJobRequest) (job links.Job, err error) {
	job = links.Job{ID: req.JobID, URLs: req.URLs}

	job, err = s.repository.CreateLinksJob(ctx, job)
	if err != nil {
		return links.Job{}, fmt.Errorf("failed to create links job %w", err)
	}

	go func() { // since we are not using some sort of distributed scheduler, start doing work in a go routine
		err := s.ExecuteLinksJob(context.Background(), job.ID)
		if err != nil {
			log.Println("failed to execute job #", job.ID)
		}
	}()

	return
}

// ExecuteJob - execute links job
func (s *service) ExecuteLinksJob(ctx context.Context, jobID string) error {
	jobResults := []links.JobResult{}

	job, err := s.repository.GetLinksJob(ctx, jobID)
	if err != nil {
		return fmt.Errorf("failed to fetch links job %w", err)
	}

	results := s.scraperClient.ScrapePages(context.Background(), job.URLs)

	for _, result := range results {
		jobResults = append(jobResults, links.JobResult{
			ID:                 uuid.NewString(),
			JobID:              job.ID,
			PageURL:            result.PageURL,
			InternalLinksCount: result.InternalLinksCount,
			ExternalLinksCount: result.ExternalLinksCount,
			Success:            result.Success,
			Error:              result.Error,
			CreatedAt:          time.Now().UTC(),
			UpdatedAt:          time.Now().UTC(),
		})
	}

	err = s.repository.CreateLinksJobResult(ctx, jobResults)
	if err != nil {
		return fmt.Errorf("failed to create links job results %w", err)
	}

	err = s.repository.FinishLinksJob(ctx, job.ID)
	if err != nil {
		return fmt.Errorf("failed to mark links job as finished %w", err)
	}
	return nil
}

// GetJobStatus - get links job status
func (s *service) GetLinksJobStatus(ctx context.Context, req links.GetJobStatusRequest) ([]links.JobResult, error) {
	_, err := s.repository.GetLinksJob(ctx, req.JobID)
	if err != nil {
		return nil, fmt.Errorf("failed to get links job results %w", err)
	}

	results, err := s.repository.GetLinksJobResult(ctx, req.JobID)
	if err != nil {
		return nil, fmt.Errorf("failed to get links job results %w", err)
	}

	return results, nil
}

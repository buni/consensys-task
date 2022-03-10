package links

import "context"

//go:generate mockgen -source=repository.go -destination=mock/repository_mocks.go -package mock

// Repository
type Repository interface {
	CreateLinksJob(ctx context.Context, job Job) (Job, error)
	FinishLinksJob(ctx context.Context, jobID string) error
	CreateLinksJobResult(ctx context.Context, results []JobResult) error
	GetLinksJobResult(ctx context.Context, jobID string) ([]JobResult, error)
	GetLinksJob(ctx context.Context, jobID string) (Job, error)
	
}

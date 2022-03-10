package links

import "context"

//go:generate mockgen -source=service.go -destination=mock/service_mocks.go -package mock

// Service
type Service interface {
	EnqueueLinksJob(ctx context.Context, req EnqueueLinksJobRequest) (job Job, err error)
	GetLinksJobStatus(ctx context.Context, req GetJobStatusRequest) ([]JobResult, error)
	ExecuteLinksJob(ctx context.Context, jobID string) error
}

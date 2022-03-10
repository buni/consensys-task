package repository

import (
	"context"
	"testing"
	"time"

	"github.com/buni/scraper/internal/api/links"
	"github.com/buni/scraper/internal/pkg/test"
	"github.com/stretchr/testify/assert"
)

func Test_inMemRepository_CreateLinksJob(t *testing.T) {
	t.Parallel()
	t.Run("successfully create job", func(t *testing.T) {
		r := NewInMemoryRepository()
		wantJob := links.Job{
			ID:        "test",
			URLs:      test.StrToURL(t, []string{"http://localhost", "https://localhost"}),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		gotJob, err := r.CreateLinksJob(context.Background(), wantJob)
		assert.NoError(t, err)
		assert.Equal(t, wantJob, gotJob)
	})
	t.Run("duplicate job err", func(t *testing.T) {
		r := NewInMemoryRepository()
		wantJob := links.Job{
			ID:        "test",
			URLs:      test.StrToURL(t, []string{"http://localhost", "https://localhost"}),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		gotJob, err := r.CreateLinksJob(context.Background(), wantJob)
		assert.NoError(t, err)
		assert.Equal(t, wantJob, gotJob)
		_, err = r.CreateLinksJob(context.Background(), wantJob)
		assert.Error(t, err)
	})
}

func Test_inMemRepository_GetLinksJob(t *testing.T) {
	t.Parallel()
	t.Run("successfully get job", func(t *testing.T) {
		r := NewInMemoryRepository()
		wantJob := links.Job{
			ID:        "test",
			URLs:      test.StrToURL(t, []string{"http://localhost", "https://localhost"}),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		r.CreateLinksJob(context.Background(), wantJob)
		gotJob, err := r.GetLinksJob(context.Background(), wantJob.ID)
		assert.NoError(t, err)
		assert.Equal(t, wantJob, gotJob)
	})
	t.Run("fail get job", func(t *testing.T) {
		r := NewInMemoryRepository()

		gotJob, err := r.GetLinksJob(context.Background(), "")
		assert.Empty(t, gotJob)
		assert.Error(t, err)
	})
}

func Test_inMemRepository_FinishLinksJob(t *testing.T) {
	t.Parallel()
	t.Run("successfully finish job", func(t *testing.T) {
		r := NewInMemoryRepository()
		wantJob := links.Job{
			ID:        "test",
			URLs:      test.StrToURL(t, []string{"http://localhost", "https://localhost"}),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		r.CreateLinksJob(context.Background(), wantJob)

		err := r.FinishLinksJob(context.Background(), wantJob.ID)
		assert.NoError(t, err)
	})
	t.Run("fail finish - job not found", func(t *testing.T) {
		r := NewInMemoryRepository()

		err := r.FinishLinksJob(context.Background(), "")
		assert.Error(t, err)
	})
}

func Test_inMemRepository_CreateLinksJobResult(t *testing.T) {
	t.Parallel()
	t.Run("successfully create job results", func(t *testing.T) {
		r := NewInMemoryRepository()
		err := r.CreateLinksJobResult(context.Background(), []links.JobResult{{ID: "test", JobID: "test"}})
		assert.NoError(t, err)
	})
}

func Test_inMemRepository_GetLinksJobResult(t *testing.T) {
	t.Parallel()
	t.Run("successfully create job results", func(t *testing.T) {
		r := NewInMemoryRepository()
		wantResults := []links.JobResult{{ID: "test", JobID: "test", PageURL: "test", InternalLinksCount: 1, ExternalLinksCount: 2, Success: true}}
		r.CreateLinksJobResult(context.Background(), wantResults)
		got, err := r.GetLinksJobResult(context.Background(), wantResults[0].JobID)
		assert.NoError(t, err)
		assert.Equal(t, wantResults, got)
	})
	t.Run("results not found", func(t *testing.T) {
		r := NewInMemoryRepository()
		got, err := r.GetLinksJobResult(context.Background(), "")
		assert.Error(t, err)
		assert.Empty(t, got)
	})
}

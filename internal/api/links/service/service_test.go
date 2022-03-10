package service_test

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/buni/scraper/internal/api/links"
	"github.com/buni/scraper/internal/api/links/mock"
	"github.com/buni/scraper/internal/api/links/repository"
	"github.com/buni/scraper/internal/api/links/service"
	"github.com/buni/scraper/internal/pkg/scraper"
	scraperMock "github.com/buni/scraper/internal/pkg/scraper/mock"
	"github.com/buni/scraper/internal/pkg/test"
	"github.com/google/uuid"

	"github.com/golang/mock/gomock"
)

func Test_service_EnqueueLinksJob(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		ctx     context.Context
		req     links.EnqueueLinksJobRequest
		setup   func(t *testing.T, mockRepo *mock.MockRepository, mockScraper *scraperMock.MockScraperService)
		wantJob links.Job
		wantErr bool
	}{
		{
			name: "successfully enqueue links job",
			req: links.EnqueueLinksJobRequest{
				JobID: uuid.Nil.String(),
				URLs: test.StrToURL(t,
					[]string{
						"http://localhost/",
						"http://localhost/page1",
					},
				),
			},
			ctx: context.Background(),
			setup: func(t *testing.T, mockRepo *mock.MockRepository, mockScraper *scraperMock.MockScraperService) {
				mockRepo.EXPECT().CreateLinksJob(gomock.Any(), gomock.Any()).Return(links.Job{
					ID: uuid.Nil.String(),
					URLs: test.StrToURL(t,
						[]string{
							"http://localhost/",
							"http://localhost/page1",
						},
					),
				}, nil)

				mockRepo.EXPECT().GetLinksJob(gomock.Any(), gomock.Any()).Return(links.Job{
					ID: uuid.Nil.String(),
					URLs: test.StrToURL(t,
						[]string{
							"http://localhost/",
							"http://localhost/page1",
						},
					),
				}, nil)

				// []string{"http://localhost/", "http://localhost/page2"}
				mockScraper.EXPECT().ScrapePages(gomock.Any(), test.StrToURL(t,
					[]string{
						"http://localhost/",
						"http://localhost/page1",
					},
				)).Return([]scraper.Result{
					{
						PageURL:            "http://localhost/",
						InternalLinksCount: 1,
						ExternalLinksCount: 1,
						Success:            true,
						Error:              nil,
					},
					{
						PageURL:            "http://localhost/page1",
						InternalLinksCount: 1,
						ExternalLinksCount: 1,
						Success:            true,
						Error:              nil,
					},
				})
				mockRepo.EXPECT().CreateLinksJobResult(gomock.Any(), gomock.Any()).Return(nil)
				mockRepo.EXPECT().FinishLinksJob(gomock.Any(), uuid.Nil.String()).Return(nil)
			},
			wantJob: links.Job{
				ID: uuid.Nil.String(),
				URLs: test.StrToURL(t,
					[]string{
						"http://localhost/",
						"http://localhost/page1",
					},
				),
			},
			wantErr: false,
		},
		{
			name: "fail to create links job",
			ctx:  context.Background(),
			req: links.EnqueueLinksJobRequest{
				JobID: uuid.Nil.String(),
				URLs: test.StrToURL(t,
					[]string{
						"http://localhost/",
						"http://localhost/page1",
					},
				),
			},
			setup: func(t *testing.T, mockRepo *mock.MockRepository, mockScraper *scraperMock.MockScraperService) {
				mockRepo.EXPECT().CreateLinksJob(gomock.Any(), gomock.Any()).Return(links.Job{}, errors.New("some error"))
			},
			wantJob: links.Job{},
			wantErr: true,
		},
		{
			name: "fail to execute links job",
			ctx:  context.Background(),
			req: links.EnqueueLinksJobRequest{
				JobID: uuid.Nil.String(),
				URLs: test.StrToURL(t,
					[]string{
						"http://localhost/",
						"http://localhost/page1",
					},
				),
			},
			setup: func(t *testing.T, mockRepo *mock.MockRepository, mockScraper *scraperMock.MockScraperService) {
				mockRepo.EXPECT().CreateLinksJob(gomock.Any(), gomock.Any()).Return(links.Job{
					ID: uuid.Nil.String(),
					URLs: test.StrToURL(t,
						[]string{
							"http://localhost/",
							"http://localhost/page1",
						},
					),
				}, nil)
				mockRepo.EXPECT().GetLinksJob(gomock.Any(), gomock.Any()).Return(links.Job{}, errors.New("some error"))
			},
			wantJob: links.Job{
				ID: uuid.Nil.String(),
				URLs: test.StrToURL(t,
					[]string{
						"http://localhost/",
						"http://localhost/page1",
					},
				),
			},
			wantErr: false,
		},
		{
			name: "fail to execute links job - create results",
			ctx:  context.Background(),
			req: links.EnqueueLinksJobRequest{
				JobID: uuid.Nil.String(),
				URLs: test.StrToURL(t,
					[]string{
						"http://localhost/",
						"http://localhost/page1",
					},
				),
			},
			setup: func(t *testing.T, mockRepo *mock.MockRepository, mockScraper *scraperMock.MockScraperService) {
				mockRepo.EXPECT().CreateLinksJob(gomock.Any(), gomock.Any()).Return(links.Job{
					ID: uuid.Nil.String(),
					URLs: test.StrToURL(t,
						[]string{
							"http://localhost/",
							"http://localhost/page1",
						},
					),
				}, nil)
				mockRepo.EXPECT().GetLinksJob(gomock.Any(), gomock.Any()).Return(links.Job{}, nil)
				mockScraper.EXPECT().ScrapePages(gomock.Any(), gomock.Any()).Return([]scraper.Result{})
				mockRepo.EXPECT().CreateLinksJobResult(gomock.Any(), gomock.Any()).Return(errors.New("some error"))
			},
			wantJob: links.Job{
				ID: uuid.Nil.String(),
				URLs: test.StrToURL(t,
					[]string{
						"http://localhost/",
						"http://localhost/page1",
					},
				),
			},
			wantErr: false,
		},
		{
			name: "fail to execute links job - finish job",
			ctx:  context.Background(),
			req: links.EnqueueLinksJobRequest{
				JobID: uuid.Nil.String(),
				URLs: test.StrToURL(t,
					[]string{
						"http://localhost/",
						"http://localhost/page1",
					},
				),
			},
			setup: func(t *testing.T, mockRepo *mock.MockRepository, mockScraper *scraperMock.MockScraperService) {
				mockRepo.EXPECT().CreateLinksJob(gomock.Any(), gomock.Any()).Return(links.Job{
					ID: uuid.Nil.String(),
					URLs: test.StrToURL(t,
						[]string{
							"http://localhost/",
							"http://localhost/page1",
						},
					),
				}, nil)
				mockRepo.EXPECT().GetLinksJob(gomock.Any(), gomock.Any()).Return(links.Job{}, nil)
				mockScraper.EXPECT().ScrapePages(gomock.Any(), gomock.Any()).Return([]scraper.Result{})
				mockRepo.EXPECT().CreateLinksJobResult(gomock.Any(), gomock.Any()).Return(nil)
				mockRepo.EXPECT().FinishLinksJob(gomock.Any(), gomock.Any()).Return(errors.New("some error"))
			},
			wantJob: links.Job{
				ID: uuid.Nil.String(),
				URLs: test.StrToURL(t,
					[]string{
						"http://localhost/",
						"http://localhost/page1",
					},
				),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			crtl := gomock.NewController(t)
			repository := mock.NewMockRepository(crtl)
			scraper := scraperMock.NewMockScraperService(crtl)
			tt.setup(t, repository, scraper)
			s := service.NewService(repository, scraper)
			gotJob, err := s.EnqueueLinksJob(tt.ctx, tt.req)
			time.Sleep(time.Second)

			if (err != nil) != tt.wantErr {
				t.Errorf("service.EnqueueLinksJob() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotJob, tt.wantJob) {
				t.Errorf("service.EnqueueLinksJob() = %v, want %v", gotJob, tt.wantJob)
			}
		})
	}
}

func Test_service_GetLinksJobStatus(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		ctx     context.Context
		req     links.GetJobStatusRequest
		setup   func(t *testing.T, mockRepo *mock.MockRepository)
		want    []links.JobResult
		wantErr bool
	}{
		{
			name: "successfully get job status",
			ctx:  context.Background(),
			req:  links.GetJobStatusRequest{JobID: uuid.Nil.String()},
			setup: func(t *testing.T, mockRepo *mock.MockRepository) {
				mockRepo.EXPECT().GetLinksJob(gomock.Any(), uuid.Nil.String()).Return(links.Job{}, nil)
				mockRepo.EXPECT().GetLinksJobResult(gomock.Any(), uuid.Nil.String()).Return([]links.JobResult{
					{
						ID:                 uuid.Nil.String(),
						JobID:              uuid.Nil.String(),
						PageURL:            "http://localhost/",
						InternalLinksCount: 1,
						ExternalLinksCount: 1,
						Success:            true,
						Error:              nil,
					},
					{
						ID:                 uuid.Nil.String(),
						JobID:              uuid.Nil.String(),
						PageURL:            "http://localhost/page2",
						InternalLinksCount: 1,
						ExternalLinksCount: 1,
						Success:            true,
						Error:              nil,
					},
				}, nil)
			},
			want: []links.JobResult{
				{
					ID:                 uuid.Nil.String(),
					JobID:              uuid.Nil.String(),
					PageURL:            "http://localhost/",
					InternalLinksCount: 1,
					ExternalLinksCount: 1,
					Success:            true,
					Error:              nil,
				},
				{
					ID:                 uuid.Nil.String(),
					JobID:              uuid.Nil.String(),
					PageURL:            "http://localhost/page2",
					InternalLinksCount: 1,
					ExternalLinksCount: 1,
					Success:            true,
					Error:              nil,
				},
			},
			wantErr: false,
		},
		{
			name: "job not found error",
			ctx:  context.Background(),
			req:  links.GetJobStatusRequest{JobID: uuid.Nil.String()},
			setup: func(t *testing.T, mockRepo *mock.MockRepository) {
				mockRepo.EXPECT().GetLinksJob(gomock.Any(), uuid.Nil.String()).Return(links.Job{}, repository.ErrJobNotFound)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "job results not found error",
			ctx:  context.Background(),
			req:  links.GetJobStatusRequest{JobID: uuid.Nil.String()},
			setup: func(t *testing.T, mockRepo *mock.MockRepository) {
				mockRepo.EXPECT().GetLinksJob(gomock.Any(), uuid.Nil.String()).Return(links.Job{}, nil)
				mockRepo.EXPECT().GetLinksJobResult(gomock.Any(), uuid.Nil.String()).Return(nil, repository.ErrJobNotFound)
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			crtl := gomock.NewController(t)
			repository := mock.NewMockRepository(crtl)
			tt.setup(t, repository)
			s := service.NewService(repository, nil)

			got, err := s.GetLinksJobStatus(tt.ctx, tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("service.GetJobStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("service.GetJobStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

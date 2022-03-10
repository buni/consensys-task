package handler

import (
	"errors"
	"net/http"

	"github.com/buni/scraper/internal/api/links"
	"github.com/buni/scraper/internal/api/links/repository"
	"github.com/buni/scraper/internal/pkg/urls"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type Handler struct {
	service links.Service
}

func NewHandler(service links.Service) *Handler {
	return &Handler{service: service}
}

// EnqueueLinksJob - handler
func (h *Handler) EnqueueLinksJob(w http.ResponseWriter, r *http.Request) {
	urls, err := urls.ParseURLs(r.Body)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, links.Response{Errors: []string{err.Error()}}) // this is treated sorta like a validation error, so it is fine to return it "naked" in the response
		return
	}

	if len(urls) == 0 {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, links.Response{Errors: []string{links.ErrEmptyJobRequest.Error()}}) // this is treated sorta like a validation error, so it is fine to return it "naked" in the response
		return
	}

	job, err := h.service.EnqueueLinksJob(r.Context(), links.EnqueueLinksJobRequest{URLs: urls})
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrJobAlreadyExists): //
			render.Status(r, http.StatusConflict)
			render.JSON(w, r, links.Response{Errors: []string{repository.ErrJobAlreadyExists.Error()}})
			return
		default:
			render.Status(r, http.StatusInternalServerError) // all other errors are treated as ise, the error message is also generic as to not leak details about the back-end
			render.JSON(w, r, links.Response{Errors: []string{links.ErrInternalServerError.Error()}})
			return
		}
	}

	render.Status(r, http.StatusAccepted)
	render.JSON(w, r, links.Response{Data: links.EnqueueLinksJobResponse{JobID: job.ID}})
}

// GetJobStatus - handler
func (h *Handler) GetJobStatus(w http.ResponseWriter, r *http.Request) {
	jobID := chi.URLParam(r, "jobID")

	results, err := h.service.GetLinksJobStatus(r.Context(), links.GetJobStatusRequest{JobID: jobID})
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrJobNotFound): // job not found
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, links.Response{Errors: []string{repository.ErrJobNotFound.Error()}})
			return
		case errors.Is(err, repository.ErrJobResultsNotFound): // job exists but is not still completed
			render.Status(r, http.StatusAccepted) // in that case we return status 202
			render.JSON(w, r, links.Response{})   // and empty body
			return
		default: // all other errors are treated as ise, the error message is also generic as to not leak details about the back-end
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, links.Response{Errors: []string{links.ErrInternalServerError.Error()}})
			return
		}
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, links.Response{Data: links.JobResultsResponse{Results: results}})
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/links", func(r chi.Router) {
		r.Post("/", h.EnqueueLinksJob)
		r.Get("/status/{jobID}", h.GetJobStatus)
	})
}

package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/buni/scraper/internal/api/links/handler"
	"github.com/buni/scraper/internal/api/links/repository"
	"github.com/buni/scraper/internal/api/links/service"

	"github.com/buni/scraper/internal/pkg/scraper"
	"github.com/go-chi/chi/v5"
)

func main() {
	log.Println("Starting API server")
	ctx := context.Background()
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	r := chi.NewRouter()

	scraperService, err := scraper.NewScraper()
	if err != nil {
		log.Fatalln(err)
	}

	jobsRepository := repository.NewInMemoryRepository()
	jobsService := service.NewService(jobsRepository, scraperService)
	jobsHandler := handler.NewHandler(jobsService)
	r.Route("/api/v1/", func(r chi.Router) {
		jobsHandler.RegisterRoutes(r)
	})

	srv := &http.Server{Handler: r, Addr: ":8080"}
	go func() {
		err := srv.ListenAndServe()
		log.Println(err)
	}()
	<-sig
	ctx, _ = context.WithTimeout(ctx, time.Second*30)
	srv.Shutdown(ctx)
	scraperService.Close(ctx)
}

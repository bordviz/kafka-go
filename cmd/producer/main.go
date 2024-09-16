package main

import (
	"context"
	"fmt"
	"go-kafka/internal/config"
	"go-kafka/internal/handlers/message"
	"go-kafka/internal/lib/logger/sl"
	"go-kafka/internal/logger"
	"go-kafka/internal/services/producer"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg := config.MustLoad()

	log := logger.New(cfg.Env)

	producerService, err := producer.NewProducerService(cfg.Kafka, log)
	if err != nil {
		log.Error("failed to create producer", sl.Err(err))
		os.Exit(1)
	}

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Route("/message", message.AddMessageHandler(log, cfg.Kafka, producerService))

	srv := http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.HTTPServer.Port),
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	go func() {
		log.Info("starting server", slog.String("addr", fmt.Sprintf("::%d", cfg.HTTPServer.Port)))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("failed to listen and serve", sl.Err(err))
			os.Exit(1)
		}
	}()

	cancelCh := make(chan os.Signal, 1)
	signal.Notify(cancelCh, syscall.SIGINT, syscall.SIGTERM)
	stopSignal := <-cancelCh
	srv.Shutdown(context.Background())
	producerService.StopProducer()
	log.Info("stoppping server", slog.String("signal", stopSignal.String()))
}

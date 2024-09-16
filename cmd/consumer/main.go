package main

import (
	"go-kafka/internal/config"
	"go-kafka/internal/lib/logger/sl"
	"go-kafka/internal/logger"
	"go-kafka/internal/services/consumer"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.MustLoad()
	log := logger.New(cfg.Env)

	consumerService, err := consumer.NewConsumerService(log, cfg.FirstConsumer)
	if err != nil {
		log.Error("failed to create consumer", sl.Err(err))
		os.Exit(1)
	}

	go consumerService.ConsumeMessages()

	cancelCh := make(chan os.Signal, 1)
	signal.Notify(cancelCh, syscall.SIGINT, syscall.SIGTERM)
	stopSignal := <-cancelCh
	consumerService.StopConsumer()
	log.Info("stoppping server", slog.String("signal", stopSignal.String()))
}

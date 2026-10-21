package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/abrshDev/ledger-service/internal/app/ledger/commands"
	"github.com/abrshDev/ledger-service/internal/infrastructure/config"
	"github.com/abrshDev/ledger-service/internal/infrastructure/database/postgres"
	"github.com/abrshDev/ledger-service/internal/infrastructure/kafka"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("configuration error: %v", err)
	}

	db, err := postgres.NewConnection(cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	repo := postgres.NewPostgresLedgerRepository(db)
	handler := commands.NewProcessPaymentEventHandler(repo)

	consumer := kafka.NewConsumer(cfg.KafkaBroker, cfg.KafkaTopic, cfg.KafkaGroup)
	defer consumer.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go consumer.Start(ctx, handler.Handle)

	log.Println("ledger-service consuming payment events...")

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh
	log.Println("shutting down...")
	cancel()
}
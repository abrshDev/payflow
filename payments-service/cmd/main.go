package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/abrshDev/payment-service/internal/app/payment/commands"
	"github.com/abrshDev/payment-service/internal/app/payment/queries"
	httpdelivery "github.com/abrshDev/payment-service/internal/delivery/http"
	"github.com/abrshDev/payment-service/internal/delivery/http/handlers"
	"github.com/abrshDev/payment-service/internal/infrastructure/config"
	"github.com/abrshDev/payment-service/internal/infrastructure/database/postgres"
	"github.com/abrshDev/payment-service/internal/infrastructure/kafka"
	"github.com/abrshDev/payment-service/internal/infrastructure/outbox"
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

	repo := postgres.NewPostgresPaymentRepository(db)

	createHandler := commands.NewCreatePaymentHandler(repo)
	authorizeHandler := commands.NewAuthorizePaymentHandler(repo)
	captureHandler := commands.NewCapturePaymentHandler(repo)
	refundHandler := commands.NewRefundPaymentHandler(repo)
	getHandler := queries.NewGetPaymentHandler(repo)

	paymentHandler := handlers.NewPaymentHandler(
		createHandler, authorizeHandler, captureHandler, refundHandler, getHandler,
	)

	router := httpdelivery.NewRouter(paymentHandler)

	producer := kafka.NewProducer(cfg.KafkaBroker, cfg.KafkaTopic)
	defer producer.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	relay := outbox.NewRelay(db, producer)
	go relay.Start(ctx)

	go func() {
		log.Printf("payments-service listening on :%s", cfg.Port)
		if err := router.Run(":" + cfg.Port); err != nil {
			log.Fatalf("server failed: %v", err)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	<-sigCh
	log.Println("shutting down...")
	cancel()
}

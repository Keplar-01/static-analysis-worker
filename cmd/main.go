package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/diploma/worker-static-analyzer/internal/analyzer"
	"github.com/diploma/worker-static-analyzer/internal/config"
	"github.com/diploma/worker-static-analyzer/internal/kafka"
	"github.com/diploma/worker-static-analyzer/internal/storage"
	"github.com/diploma/worker-static-analyzer/internal/usecase"
)

func main() {
	log.Println("[worker-static] starting...")

	cfg := config.Load()

	// --- Infrastructure ---

	minioClient, err := storage.NewMinIOClient(cfg.MinioEndpoint, cfg.MinioAccessKey, cfg.MinioSecretKey)
	if err != nil {
		log.Fatalf("[worker-static] minio connection failed: %v", err)
	}

	producer := kafka.NewProducer(cfg.KafkaBrokers)
	defer producer.Close()

	// --- Business logic ---

	llvmAnalyzer := analyzer.New(cfg.AnalyzerBinary)
	analysisUC := usecase.NewAnalysisUseCase(llvmAnalyzer, minioClient, producer)

	// --- Kafka consumer ---

	consumer := kafka.NewConsumer(cfg.KafkaBrokers, analysisUC)
	defer consumer.Close()

	// --- Graceful shutdown ---

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		log.Println("[worker-static] shutting down...")
		cancel()
	}()

	consumer.Listen(ctx)
}

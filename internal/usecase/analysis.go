package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/diploma/worker-static-analyzer/internal/analyzer"
	"github.com/diploma/worker-static-analyzer/internal/kafka"
	"github.com/diploma/worker-static-analyzer/internal/model"
	"github.com/diploma/worker-static-analyzer/internal/storage"
)

type AnalysisUseCase struct {
	analyzer *analyzer.Analyzer
	minio    *storage.MinIOClient
	producer *kafka.Producer
}

func NewAnalysisUseCase(
	a *analyzer.Analyzer,
	m *storage.MinIOClient,
	p *kafka.Producer,
) *AnalysisUseCase {
	return &AnalysisUseCase{
		analyzer: a,
		minio:    m,
		producer: p,
	}
}

func (uc *AnalysisUseCase) HandleStartEvent(ctx context.Context, event model.StartEvent) {
	log.Printf("[usecase] processing task %s, file: %s", event.TaskID, event.FileS3Path)

	if err := uc.process(ctx, event); err != nil {
		log.Printf("[usecase] task %s failed: %v", event.TaskID, err)
		uc.sendCompleted(ctx, event.TaskID, "error", "", err.Error())
		return
	}

	log.Printf("[usecase] task %s completed successfully", event.TaskID)
}

func (uc *AnalysisUseCase) process(ctx context.Context, event model.StartEvent) error {
	workDir, err := os.MkdirTemp("", "static-analysis-*")
	if err != nil {
		return fmt.Errorf("create temp dir: %w", err)
	}
	defer os.RemoveAll(workDir)

	sourceFile, err := uc.minio.DownloadSource(ctx, event.FileS3Path, workDir)
	if err != nil {
		return fmt.Errorf("download source: %w", err)
	}

	patterns, err := uc.analyzer.Run(ctx, sourceFile, workDir)
	if err != nil {
		return fmt.Errorf("run analyzer: %w", err)
	}

	log.Printf("[usecase] analysis done: %d patterns found", len(patterns))

	outputJSON, err := json.MarshalIndent(patterns, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal patterns: %w", err)
	}

	artifactPath, err := uc.minio.UploadArtifact(ctx, event.TaskID, outputJSON)
	if err != nil {
		return fmt.Errorf("upload artifact: %w", err)
	}

	uc.sendCompleted(ctx, event.TaskID, "success", artifactPath, "")
	return nil
}

func (uc *AnalysisUseCase) sendCompleted(ctx context.Context, taskID, status, artifactPath, errMsg string) {
	event := model.CompletedEvent{
		TaskID:         taskID,
		Status:         status,
		ArtifactS3Path: artifactPath,
		Error:          errMsg,
	}

	if err := uc.producer.Publish(ctx, kafka.TopicStaticCompleted, taskID, event); err != nil {
		log.Printf("[usecase] failed to send completed event: %v", err)
	}
}

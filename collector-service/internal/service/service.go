package service

import (
	"collector-service/internal/models"
	"collector-service/internal/repository"
	"context"
)

//go:generate mockery --name=LogServiceMock --output=../mocks --outpkg=mock
type LogService interface {
	StoreBatch(ctx context.Context, batch []models.TelemetryRecord) error
	Stop() error
}

type logService struct {
	repo repository.Repository
}

func (l *logService) StoreBatch(ctx context.Context, batch []models.TelemetryRecord) error {
	return l.repo.InsertBatch(ctx, batch)
}

func NewLogService(repo repository.Repository) *logService {
	return &logService{repo: repo}
}

// closing connections to db
func (l *logService) Stop() error {
	return l.repo.Stop()
}

package repository

import (
	"context"
	"github.com/google/uuid"
	"go_ping_kube/internal/domain/models"
)

//go:generate mockgen -build_flags=-mod=mod -destination mock/ping_repository_mock.go go_ping_kube/internal/infrastructure/repository IPingRepository

type IPingRepository interface {
	Save(ctx context.Context, ping *models.PingData) error
	Get(ctx context.Context, id uuid.UUID) (*models.PingData, error)
	All(ctx context.Context) ([]*models.PingData, error)
}

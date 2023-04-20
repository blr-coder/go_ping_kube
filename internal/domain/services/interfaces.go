package services

import (
	"context"
	"github.com/google/uuid"
	"go_ping_kube/internal/domain/models"
)

//go:generate mockgen -build_flags=-mod=mod -destination mock/ping_service_mock.go go_ping_kube/internal/domain/services IPingService

type IPingService interface {
	Add(ctx context.Context, add *models.CreatePingData) (*models.PingData, error)
	Get(ctx context.Context, uuid uuid.UUID) (*models.PingData, error)
	All(ctx context.Context) ([]*models.PingData, error)
}

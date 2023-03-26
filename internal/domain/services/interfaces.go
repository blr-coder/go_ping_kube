package services

import (
	"context"
	"github.com/google/uuid"
	"go_ping_kube/internal/domain/models"
)

type IPingService interface {
	Add(ctx context.Context, add *models.CreatePingData) (*models.PingData, error)
	Get(ctx context.Context, uuid uuid.UUID) (*models.PingData, error)
	All(ctx context.Context) ([]*models.PingData, error)
}

package repository

import (
	"context"
	"github.com/google/uuid"
	"go_ping_kube/internal/domain/models"
)

type IPingRepository interface {
	Save(ctx context.Context, ping *models.PingData) error
	Get(ctx context.Context, id uuid.UUID) (*models.PingData, error)
	All(ctx context.Context) ([]*models.PingData, error)
}

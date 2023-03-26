package webapi

import (
	"context"
	"go_ping_kube/internal/infrastructure/webapi/event"
)

type EventSender interface {
	Send(ctx context.Context, event event.CreateEvent) error
}

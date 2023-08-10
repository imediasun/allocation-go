package service

import (
	"context"
	"gitlab.hotel.tools/backend-team/allocation-go/internal/domain/model"
)

type AllocatorService interface {
	AllocateAll(ctx context.Context, reservationIDs []int32, userID *int32) ([]model.AllocateResult, error)
	AutoAllocate(ctx context.Context, reservationID int, isNotify bool)
}

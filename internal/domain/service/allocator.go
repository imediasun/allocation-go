package service

import (
	"context"
)

type AllocatorService interface {
	AllocateAll(ctx context.Context, reservationIDs []int32, userID *int32) ([]byte, error)
	AutoAllocate(ctx context.Context, userID *int32, reservationID int32, isNotify bool)
}

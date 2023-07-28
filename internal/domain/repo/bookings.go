package repo

import (
	"context"
	"gitlab.hotel.tools/backend-team/allocation-go/internal/common/adapter/db"
	"gitlab.hotel.tools/backend-team/allocation-go/internal/domain/model"
)

type BookingRepoFactory interface {
	Create(ctx context.Context, db db.DB) (BookingRepo, error)
}

type BookingRepo interface {
	Create(context.Context, *model.Booking) error
	Delete(context.Context, int32) error
}

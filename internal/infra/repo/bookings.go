package repo

import (
	"context"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"gitlab.hotel.tools/backend-team/allocation-go/internal/common/adapter/db"
	"gitlab.hotel.tools/backend-team/allocation-go/internal/common/adapter/log"
	"gitlab.hotel.tools/backend-team/allocation-go/internal/domain/adapter"
	"gitlab.hotel.tools/backend-team/allocation-go/internal/domain/model"
	"gitlab.hotel.tools/backend-team/allocation-go/internal/domain/repo"
	"gitlab.hotel.tools/backend-team/allocation-go/internal/infra/repo/entity"

	"go.uber.org/zap"
)

type bookingRepoFactory struct {
	logger                log.Logger
	bookingAdapterFactory adapter.BookingAdapterFactory
}

func NewBookingRepoFactory(logger log.Logger, bookingAdapterFactory adapter.BookingAdapterFactory) repo.BookingRepoFactory {
	return &bookingRepoFactory{
		logger:                logger,
		bookingAdapterFactory: bookingAdapterFactory,
	}
}

func (f *bookingRepoFactory) Create(ctx context.Context, db db.DB) (repo.BookingRepo, error) {
	return &bookingRepo{
		logger:                f.logger.WithComponent(ctx, "BookingRepo"),
		db:                    db,
		bookingAdapterFactory: f.bookingAdapterFactory,
	}, nil
}

type bookingRepo struct {
	logger                log.Logger
	db                    db.DB
	bookingAdapterFactory adapter.BookingAdapterFactory
}

func (r bookingRepo) Create(ctx context.Context, booking *model.Booking) error {
	ent, err := r.bookingAdapterFactory.Create(ctx).ToEntity(ctx, booking)
	if err != nil {
		r.logger.WithMethod(ctx, "Create").Error("failed to convert booking to entity", zap.Error(err))
		return err
	}

	if err := ent.Insert(ctx, r.db, boil.Blacklist(entity.BookingColumns.ID)); err != nil {
		r.logger.WithMethod(ctx, "Create").Error("failed to insert booking", zap.Error(err))
	}

	booking.ID = ent.ID

	return nil

}

func (r bookingRepo) Delete(ctx context.Context, bookingID int32) error {
	ent := entity.Booking{ID: bookingID}

	if _, err := ent.Delete(ctx, r.db); err != nil {
		r.logger.WithMethod(ctx, "Delete").Error("failed to delete booking", zap.Error(err))
	}

	return nil

}

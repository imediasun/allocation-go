package adapter

import (
	"context"
	"gitlab.hotel.tools/backend-team/allocation-go/internal/domain/model"
	"gitlab.hotel.tools/backend-team/allocation-go/internal/infra/repo/entity"
)

type BookingAdapterFactory interface {
	Create(context.Context) BookingAdapter
}

type BookingAdapter interface {
	FromEntity(context.Context, *entity.Booking) (*model.Booking, error)
	FromEntities(context.Context, entity.BookingSlice) (model.Bookings, error)
	ToEntity(context.Context, *model.Booking) (*entity.Booking, error)
	ToEntities(context.Context, model.Bookings) (entity.BookingSlice, error)
}

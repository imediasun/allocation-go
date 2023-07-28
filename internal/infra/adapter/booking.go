package adapter

import (
	"context"
	"github.com/volatiletech/null/v8"
	"gitlab.hotel.tools/backend-team/allocation-go/internal/common/adapter/log"
	"gitlab.hotel.tools/backend-team/allocation-go/internal/domain/adapter"
	"gitlab.hotel.tools/backend-team/allocation-go/internal/domain/model"
	"gitlab.hotel.tools/backend-team/allocation-go/internal/infra/repo/entity"
	"time"
)

type BookingAdapterFactory struct {
	logger log.Logger
}

func NewBookingAdapterFactory(logger log.Logger) adapter.BookingAdapterFactory {
	return &BookingAdapterFactory{
		logger: logger,
	}
}

type BookingAdapter struct {
	logger log.Logger
}

func (f *BookingAdapterFactory) Create(ctx context.Context) adapter.BookingAdapter {
	return &BookingAdapter{
		logger: f.logger.WithComponent(ctx, "BookingAdapter"),
	}
}
func (a *BookingAdapter) FromEntity(ctx context.Context, booking *entity.Booking) (*model.Booking, error) {
	res := &model.Booking{
		ID:                booking.ID,
		Price:             0,
		Currency:          "",
		AgentID:           0,
		Date:              null.String{},
		ProviderReference: null.String{},
		Channel:           null.String{},
		Status:            null.String{},
		IsManual:          false,
		Client:            null.String{},
		CancellationDate:  null.Time{},
		PaymentOption:     model.PaymentOptionFromEntity(booking.PaymentOption),
		Segment:           null.String{},
		Source:            null.String{},
		FOCT:              false,
		MetaGroupID:       null.Uint32{},
		ClientID:          null.Uint32{},
		IsVirtualCC:       false,
		//Color:                  null.String{},
		ChannelCommissionType:  nil,
		ChannelCommissionValue: null.Float32{},
		ReleazeTime:            null.Int32{},
		RequestedCheckInTime:   "",
		RequestedCheckOutTime:  "",
	}

	return res, nil
}

func (a *BookingAdapter) FromEntities(ctx context.Context, logs entity.BookingSlice) (model.Bookings, error) {
	result := make(model.Bookings, len(logs))
	for i, log := range logs {
		var err error
		result[i], err = a.FromEntity(ctx, log)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (a *BookingAdapter) ToEntity(ctx context.Context, booking *model.Booking) (*entity.Booking, error) {
	ent := &entity.Booking{
		ID:                booking.ID,
		Price:             booking.Price,
		Currency:          booking.Currency,
		AgentID:           booking.AgentID,
		Date:              time.Time{},
		ProviderReference: null.String{},
		Channel:           null.String{},
		Status:            null.String{},
		IsManual:          false,
		Client:            null.String{},
		CancellationDate:  null.Time{},
		PaymentOption:     booking.PaymentOption.ToEntity(),
		Taxes:             "",
		Segment:           null.String{},
		Source:            null.String{},
		FOCT:              false,
		MetaGroupID:       null.Uint32{},
		ClientID:          null.Uint32{},
		IsVirtualCC:       false,
		//Color:                  null.String{},
		ChannelCommissionType:  entity.BookingsNullChannelCommissionType{},
		ChannelCommissionValue: null.Float32{},
		ReleazeTime:            null.Int32{},
		RequestedCheckInTime:   "",
		RequestedCheckOutTime:  "",
	}

	return ent, nil
}

func (a *BookingAdapter) ToEntities(ctx context.Context, logs model.Bookings) (entity.BookingSlice, error) {
	result := make(entity.BookingSlice, len(logs))
	for i, log := range logs {
		var err error
		result[i], err = a.ToEntity(ctx, log)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

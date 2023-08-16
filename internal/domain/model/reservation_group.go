package model

import (
	"github.com/volatiletech/null/v8"
	"time"
)

type ReservationGroup struct {
	Item           BookingItems
	ID             int32
	BookingID      int32
	PaxNationality null.String
	StartDate      time.Time
	EndDate        time.Time
	ParentID       null.Int64
	Items          []BookingItems
}

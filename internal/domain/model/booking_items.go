package model

import (
	"github.com/volatiletech/null/v8"
)

type BookingItems struct {
	ID        int         `db:"id"`
	Type      string      `db:"Type"`
	VenueID   int32       `db:"VenueID"`
	ProductID null.String `db:"ProductID"`
	Status    string      `db:"Status"`
	Product   Product
}

package model

import "time"

type CurrencyRate struct {
	BookingID int64
	Source    string
	Target    string
	Rate      float64
	Date      time.Time
	Final     bool
}

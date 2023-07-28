package model

type ReservationStatus string

const (
	FAILED              ReservationStatus = "failed"
	CANCELLED           ReservationStatus = "cancelled"
	CONFIRMED           ReservationStatus = "confirmed"
	ABORTED             ReservationStatus = "aborted"
	PARTIALLY_CONFIRMED ReservationStatus = "partially-confirmed"
	ON_REQUEST          ReservationStatus = "on-request"
	REJECTED            ReservationStatus = "rejected"
	QUOTE               ReservationStatus = "quote"
)

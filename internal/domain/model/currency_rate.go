package model

type History struct {
	User   Agent
	Action string
	Before Reservation
	After  Reservation
}

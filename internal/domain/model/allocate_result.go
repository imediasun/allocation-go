package model

type AllocateResult struct {
	Status        string         `json:"status"`
	BookingID     string         `json:"bookingId"`
	GroupID       string         `json:"groupId"`
	ItemID        string         `json:"itemId"`
	AllocatedRoom *AllocatedRoom `json:"allocatedRoom"`
	Reason        string         `json:"reason"`
}

package model

type AllocateResult struct {
	Status        string      `json:"status"`
	BookingID     int         `json:"bookingID"`
	GroupID       int32       `json:"groupID"`
	ItemID        int         `json:"itemID"`
	AllocatedRoom MetaObjects `json:"allocatedRoom"`
	Reason        string      `json:"reason"`
}

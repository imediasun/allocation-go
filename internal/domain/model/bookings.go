package model

import (
	"github.com/volatiletech/null/v8"
)

type Booking struct {
	ID                int32
	Price             float32
	Currency          string
	AgentID           uint32
	Date              null.String
	ProviderReference null.String
	Channel           null.String
	Status            null.String
	IsManual          bool
	Client            null.String
	CancellationDate  null.Time
	PaymentOption     *PaymentOption
	Segment           null.String
	Source            null.String
	FOCT              bool
	MetaGroupID       null.Uint32
	ClientID          null.Uint32
	IsVirtualCC       bool
	//Color                  null.String
	ChannelCommissionType  *ChannelCommissionType
	ChannelCommissionValue null.Float32
	ReleazeTime            null.Int32
	RequestedCheckInTime   string
	RequestedCheckOutTime  string
}

type Bookings []*Booking

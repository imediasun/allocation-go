package model

import (
	"github.com/volatiletech/null/v8"
	"time"
)

type Reservation struct {
ID                  int
Creator             Agent
Price               Money
CreationDate        time.Time
Status              ReservationStatus
ProviderReference   null.String
Channel             null.String
Remark              string
Client              null.String
Manual              bool
PaymentOption       null.String
Groups              []ReservationGroup
CancellationDate    []uint8
StartDate           []uint8
EndDate             []uint8
Segment             null.String
Source              null.String
Logs                interface{} // Replace with actual type
CurrencyRates       []CurrencyRate
Foct                bool
IsCityTaxToProvider bool
MetaGroupID         null.Int64
Customer            *Client
}
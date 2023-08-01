package model

import "github.com/volatiletech/null/v8"

type Client struct {
	ID             null.Int32  `db:"ID"`
	AccountID      null.Int32  `db:"AccountID"`
	Email          null.String `db:"Email"`
	Phone          null.String `db:"Phone"`
	Title          null.String `db:"Title"`
	Gender         null.String `db:"Gender"`
	Nationality    null.String `db:"Nationality"`
	LanguageID     null.Int32  `db:"LanguageID"`
	Identification null.String `db:"Identification"`
	LastName       null.String `db:"LastName"`
	BirthDate      []uint8     `db:"BirthDate"`
	Address        null.String `db:"Address"`
	AdditionalInfo null.String `db:"AdditionalInfo"`
	AgentID        null.Int32  `db:"AgentID"`
	CreatedAt      []uint8     `db:"CreatedAt"`
	Status         null.String `db:"Status"`
}

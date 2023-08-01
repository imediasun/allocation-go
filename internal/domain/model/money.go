package model

type Agent struct {
	ID        int32  `db:"id"`
	Name      string `db:"name"`
	AccountID int32  `db:"accountID"`
}
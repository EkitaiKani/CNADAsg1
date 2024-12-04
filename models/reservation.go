package models

import "database/sql"

// for reservation table
type Reservation struct {
	ReservationId int          `json:"id"`
	UserId        int          `json:"userId"`
	CarId         int          `json:"carId"`
	Date          sql.NullTime `json: "date"`
	Start         sql.NullTime `json:"startTime"`
	End           sql.NullTime `json:"endTime"`
	Status        string       `json: status`
}

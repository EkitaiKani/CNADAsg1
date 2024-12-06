package models

import "database/sql"

// for reservation table
type Reservation struct {
	ReservationId int          `json:"id"`
	UserId        int          `json:"userId"`
	CarId         int          `json:"carId"`
	Start         sql.NullTime `json:"startTime"`
	End           sql.NullTime `json:"endTime"`
	Status        string       `json: status`
	CarDetails    Car          `json:"car_details"`
}

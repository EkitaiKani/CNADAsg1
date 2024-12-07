package models

import "database/sql"

// for Payment table
type Payment struct {
	PaymentId     int          `json:"id"`
	ReservationId int          `json:"resId"`
	UserId        int          `json:"userId"`
	Method        string       `json:"method"`
	Status        string       `json: status`
	Amount        float32      `json: amount`
	TransactionId string       `json:"transactionId"`
	Date          sql.NullTime `json: dateCompleted`
	Reservation   Reservation  `json: reservation`
}

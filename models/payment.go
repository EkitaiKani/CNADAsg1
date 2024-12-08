package models

import "database/sql"

// for Payment table
type Payment struct {
	PaymentId     int          `json:"id"`
	ReservationId int          `json:"resId"`
	UserId        int          `json:"userId"`
	Method        string       `json:"method"`
	Status        string       `json: status`
	TotalAmount   float32      `json: amount`
	TransactionId string       `json:"transactionId"`
	Date          sql.NullTime `json: dateCompleted`
	Reservation   Reservation  `json: reservation`
	Discount      float32      `json: discount`
	AmtPayable    float32      `json: amtPay`
}

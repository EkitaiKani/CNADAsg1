package models

import "database/sql"

// for Payment table
type Payment struct {
	PaymentId     string       `json:"id"`
	RentalId      string       `json:"rentalId"`
	UserId        string       `json:"userId"`
	Method        int          `json:"method"`
	Status        string       `json: status`
	Amount        float32      `json: amount`
	TransactionId string       `json:"transactionId"`
	Date          sql.NullTime `json: dateCompleted`
}

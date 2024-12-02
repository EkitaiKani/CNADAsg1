package rental

import "database/sql"

// for payment table
type payment struct {
	paymentId     string  `json:"id"`
	rentalId      string  `json:"rentalId"`
	userId        string  `json:"userId"`
	method        int     `json:"method"`
	status        string  `json: status`
	amount        float32 `json: amount`
	transactionId string  `json:"transactionId"`
	date          sql.NullTime  `json: dateCompleted`
}

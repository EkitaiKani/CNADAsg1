package rental

import "database/sql"

// for reservation table
type reservation struct {
	reservationId string 	`json:"id"`
	userId        int 		`json:"userId"`
	carId         int    	`json:"carId"`
	start         sql.NullTime `json:"startTime"`
	end           sql.NullTime `json:"endTime"`
	status        string 	`json: status`
}
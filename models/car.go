package car

import "database/sql"

// for car table
type car struct {
	carId         string       `json:"id"`
	carModel      string       `json:"model"`
	liscencePlate string       `json:"plate"`
	status        string       `json: status`
	currLoc       string       `json:"loc"`
	charge        int          `json:"charge"`
	cleanliness   string       `json: cleanliness`
	lastServiced  sql.NullTime `json: dateServiced`
	addedDate     sql.NullTime `json: dateAdded`
}
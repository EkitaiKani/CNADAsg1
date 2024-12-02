package models

import "database/sql"

// for car table
type Car struct {
	CarId         string       `json:"id"`
	CarModel      string       `json:"model"`
	LiscencePlate string       `json:"plate"`
	Status        string       `json: status`
	CurrLoc       string       `json:"loc"`
	Charge        int          `json:"charge"`
	Cleanliness   string       `json: cleanliness`
	LastServiced  sql.NullTime `json: dateServiced`
	AddedDate     sql.NullTime `json: dateAdded`
}

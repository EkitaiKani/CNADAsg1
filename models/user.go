package models

import "database/sql"

// for user table
type User struct {
	UserId       int          `json:"id"`
	UserEmail    string       `json:"name"`
	//PhoneNo      int          `json:"phone"`
	UserName     string       `json: username`
	HashPassword string       `json:"pw"`
	MemberTier   string       `json:"tier"`
	FirstName    string       `json: firstname`
	LastName     string       `json: lastname`
	DateofBirth  sql.NullTime `json: dob`
}

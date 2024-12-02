package user

import "database/sql"

// for user table
type user struct {
	userId       string    `json:"id"`
	userEmail    string    `json:"name"`
	phoneNo      int       `json:"phone"`
	userName     string    `json: username`
	hashPassword string    `json:"pw"`
	memberTier   string    `json:"tier"`
	name         string    `json: name`
	dateofBirth  sql.NullTime `json: dob`
	isVerified   bool      `json: verified`
}
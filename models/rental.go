package models

// for rental table
type Rental struct {
	RentalId    string `json:"id"`
	UserId      string `json:"userId"`
	CarId       int    `json:"carId"`
	Start       string `json:"startTime"`
	End         int    `json:"endTime"`
	TotalAmount int    `json:"amount"`
	Status      string `json: status`
}

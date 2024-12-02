package rental

// for rental table
type rental struct {
	rentalId    string `json:"id"`
	userId      string `json:"userId"`
	carId       int    `json:"carId"`
	start       string `json:"startTime"`
	end         int    `json:"endTime"`
	totalAmount int    `json:"amount"`
	status      string `json: status`
}
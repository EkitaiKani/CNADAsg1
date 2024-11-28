package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

// for user table
type User struct {
	userId       string    `json:"id"`
	userEmail    string    `json:"name"`
	phoneNo      int       `json:"phone"`
	userName     string    `json: username`
	hashPassword int       `json:"pw"`
	memberTier   int       `json:"tier"`
	name         string    `json: name`
	dateofBirth  time.Time `json: dob`
	isVerified   bool      `json: verified`
}

// for car table
type car struct {
	carId         string    `json:"id"`
	carModel      string    `json:"model"`
	liscencePlate string    `json:"plate"`
	status        string    `json: status`
	currLoc       string    `json:"loc"`
	charge        int       `json:"charge"`
	cleanliness   string    `json: cleanliness`
	lastServiced  time.Time `json: dateServiced`
	addedDate     time.Time `json: dateAdded`
}

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

// for reservation table
type reservation struct {
	reservationId string `json:"id"`
	userId        string `json:"userId"`
	carId         int    `json:"carId"`
	start         string `json:"startTime"`
	end           int    `json:"endTime"`
	status        string `json: status`
}

// for reservation table
type payment struct {
	paymentId     string  `json:"id"`
	rentalId      string  `json:"rentalId"`
	userId        string  `json:"userId"`
	method        int     `json:"method"`
	status        string  `json: status`
	amount        float32 `json: amount`
	transactionId int     `json:"transactionId"`
	date          string  `json: dateCompleted`
}

var db *sql.DB

// Set env varibles in CLI
// set DB_USER=root // create new user in DB
// set DB_PASSWORD=XXXXXXXXXXXX
// set DB_HOST=127.0.0.1:3306
// set DB_NAME=ICTcourses

func main() {

	r := mux.NewRouter()

	// Basic route handler
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, world!")
	})
}

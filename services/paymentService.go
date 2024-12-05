package services

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

type PaymentService struct {
	DB *sql.DB
}

// fetches Car details for the Cars page
func (s *PaymentService) CalculatePayment(id int) (*float32, error) {

	// define start and end
	var start sql.NullTime
	var end sql.NullTime
	var carId int

	// get total time and carId from reservation
	query := "SELECT start_datetime, end_datetime, car_id FROM reservations WHERE reservation_id = ?"
	err := s.DB.QueryRow(query, id).Scan(&start, &end, &carId)
	if err != nil {
		log.Println("Row scan error:", err)
		return nil, err
	}

	// Set the end time to the current time if it's nil
	if !end.Valid {
		end.Time = time.Now() // set end to the current time
		end.Valid = true      // mark it as valid
	}

	// get total duration
	var duration time.Duration
	if start.Valid && end.Valid {
		// Calculate the difference between the two times
		duration = end.Time.Sub(start.Time)
	} else {
		// If the dates are invalid, return an error
		return nil, fmt.Errorf("invalid start or end datetime")
	}

	// get rate for car
	var rate int
	query = "SELECT rate FROM cars WHERE car_id = ?"
	err = s.DB.QueryRow(query, carId).Scan(&rate)
	if err != nil {
		log.Println("Row scan error:", err)
		return nil, err
	}

	// Calculate the total payment, converting rate to rate per minute
	totalPayment := float32(rate) / 60 * float32(duration.Minutes())

	// Return the total payment as a pointer to float32
	return &totalPayment, nil
}


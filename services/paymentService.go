package services

import (
	"CNADASG1/models"
	"database/sql"
	"fmt"
	"log"
	"math"
	"time"
)

type PaymentService struct {
	DB *sql.DB
}

func (s *PaymentService) CalculatePayment(p models.Payment) (*models.Payment, error) {

	// Define start, end, carId, and userId
	var startStr, endStr string
	var carId, userId int

	//log.Print(id)

	// Get total time and carId from reservation
	query := "SELECT user_id, start_datetime, end_datetime, car_id FROM reservations WHERE reservation_id = ?"
	err := s.DB.QueryRow(query, p.ReservationId).Scan(&userId, &startStr, &endStr, &carId)
	if err != nil {
		log.Printf("Error retrieving reservation: %v", err)
		return nil, err
	}

	var start, end sql.NullTime
	if startStr != "" {
		start.Time, err = time.Parse("2006-01-02 15:04:05", startStr)
		if err != nil {
			log.Printf("Error parsing start time: %v", err)
			return nil, err
		}
	}

	if endStr != "" {
		end.Time, err = time.Parse("2006-01-02 15:04:05", endStr)
		if err != nil {
			log.Printf("Error parsing end time: %v", err)
			return nil, err
		}
	}

	// Calculate the duration between start and end times
	duration := end.Time.Sub(start.Time)
	if duration < 0 {
		// Check if the duration is negative (start time after end time)
		log.Println("End time is before start time")
		return nil, fmt.Errorf("end time is before start time")
	}

	//log.Print(duration)

	// Get rate for car
	var rate int
	query = "SELECT rate FROM cars WHERE car_id = ?"
	err = s.DB.QueryRow(query, carId).Scan(&rate)
	if err != nil {
		log.Printf("Error retrieving car rate: %v", err)
		return nil, err
	}

	// Calculate the total payment, converting rate to rate per minute
	totalPayment := float32(rate) / 60 * float32(duration.Minutes())

	// Get user membership tier for discount
	var tier string
	query = "SELECT membership_tier FROM users WHERE user_id = ?"
	err = s.DB.QueryRow(query, userId).Scan(&tier)
	if err != nil {
		log.Printf("Error retrieving membership tier: %v", err)
		return nil, err
	}

	// Get discount based on membership tier
	var discount int
	query = "SELECT discount_percentage FROM MembershipTiers WHERE tier_name = ?"
	err = s.DB.QueryRow(query, tier).Scan(&discount)
	if err != nil {
		log.Printf("Error retrieving discount percentage: %v", err)
		return nil, err
	}

	p.TotalAmount = totalPayment

	// Apply discount if applicable
	discountAmt := totalPayment * float32(discount) / 100 // Divide by 100 to get the percentage
	totalPayment = totalPayment - discountAmt

	// Round up to 2 decimal places
	roundedPayment := math.Round(float64(totalPayment)*100) / 100
	p.AmtPayable = float32(roundedPayment)
	p.Discount = discountAmt

	// Return the total payment as a pointer to float32
	return &p, nil
}

func (s *PaymentService) CreatePayment(pay *models.Payment) (*models.Payment, error) {

	// Prepare the SQL INSERT statement
	query := "INSERT INTO payments (reservation_id, user_id, amount, transaction_id) VALUES (?, ?, ?, ?)"
	result, err := s.DB.Exec(query, pay.ReservationId, pay.UserId, pay.AmtPayable, pay.TransactionId)
	if err != nil {
		log.Println("Database insert error:", err)
		return nil, err
	}

	// Get the last inserted ID
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		log.Println("Error getting last insert ID:", err)
		return nil, err
	}

	// Set the UserId in the user model
	pay.PaymentId = int(lastInsertID)

	return pay, nil
}

// get payment details
func (s *PaymentService) GetPayment(id int) (*models.Payment, error) {
	// create var to store Payment
	var p models.Payment

	p.PaymentId = id

	// Get Payment record
	query := "SELECT reservation_id, amount, transaction_id, payment_date FROM payments WHERE payment_id = ?"
	err := s.DB.QueryRow(query, id).Scan(&p.ReservationId, &p.AmtPayable, &p.TransactionId, &p.Date)
	if err != nil {
		log.Printf("Query error: %v", err)
		return nil, err
	}

	// create res var
	var r models.Reservation
	var start, end sql.NullString

	// get reservation details
	query = "SELECT car_id, start_datetime, end_datetime FROM reservations WHERE reservation_id = ?"
	err = s.DB.QueryRow(query, p.ReservationId).Scan(&r.CarId, &start, &end)
	if err != nil {
		log.Println("Row scan error:", err)
		return nil, err
	}
	// Parse start datetime
	if start.Valid {
		parsedTime, err := time.Parse("2006-01-02 15:04:05", start.String)
		if err != nil {
			log.Printf("Start datetime parse error: %v", err)
			return nil, err
		}
		r.Start = sql.NullTime{Time: parsedTime, Valid: true}
	}

	// Parse end datetime
	if end.Valid {
		parsedTime, err := time.Parse("2006-01-02 15:04:05", end.String)
		if err != nil {
			log.Printf("End datetime parse error: %v", err)
			return nil, err
		}
		r.End = sql.NullTime{Time: parsedTime, Valid: true}
	}

	// get car details
	// Create a new instance of models.Car
	u := &models.Car{CarId: r.CarId}
	var lastServiced string

	// Get car details from the database
	query = "SELECT car_model, license_plate, status, current_location, charge_level, cleanliness_status, last_serviced, rate FROM cars WHERE car_id = ?"
	err = s.DB.QueryRow(query, r.CarId).Scan(
		&u.CarModel,      // car_model
		&u.LiscencePlate, // license_plate (fixed typo)
		&u.Status,        // status
		&u.CurrLoc,       // current_location
		&u.Charge,        // charge_level
		&u.Cleanliness,   // cleanliness_status
		&lastServiced,    // last_serviced
		&u.Rate,          // rate
	)

	// Handle scan error
	if err != nil {
		if err == sql.ErrNoRows {
			// No rows found for the given car ID
			log.Println("No car found with ID:", r.CarId)
			return nil, nil
		}
		// Other query error
		log.Println("Row scan error:", err)
		return nil, err
	}

	// Parse the last serviced date
	if lastServiced != "" {
		parsedTime, err := time.Parse("2006-01-02 15:04:05", lastServiced)
		if err != nil {
			log.Printf("Last serviced datetime parse error: %v", err)
			return nil, err
		}
		u.LastServiced = sql.NullTime{Time: parsedTime, Valid: true}
	} else {
		u.LastServiced = sql.NullTime{Valid: false} // Handle NULL case
	}

	// Attach car details to reservation
	r.CarDetails = u

	// attach res details to payment
	p.Reservation = r

	return &p, nil
}

func (s *PaymentService) MakePayment(pay models.Payment) (*models.Payment, error) {

	// Ensure we are using the correct time format, in UTC if needed
	paymentDate := time.Now()

	// Prepare the SQL UPDATE statement
	query := "UPDATE payments SET Status = ?, payment_method = ?, payment_date = ? WHERE reservation_id = ?"
	_, err := s.DB.Exec(query, pay.Status, pay.Method, paymentDate)
	if err != nil {
		log.Println("Database update error:", err)
		return nil, err
	}

	return &pay, nil
}

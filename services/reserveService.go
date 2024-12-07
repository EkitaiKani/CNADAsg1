package services

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"CNADASG1/models"
)

type ReserveService struct {
	DB *sql.DB
}

type AvailableSlot struct {
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

// fetches Reserve details for the Reserves page
func (s *ReserveService) GetCarReservations(id int) (map[int]models.Reservation, error) {
	// create dictionary to store Reservations
	resList := make(map[int]models.Reservation)

	// Get Reservations
	query := "SELECT reservation_id, start_datetime, end_datetime FROM reservations WHERE car_id = ? AND status NOT IN ('Cancelled', 'Completed')"
	rows, err := s.DB.Query(query, id)
	if err != nil {
		log.Printf("Query error: %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		r := models.Reservation{}
		var start, end sql.NullString

		if err := rows.Scan(&r.ReservationId, &start, &end); err != nil {
			log.Printf("Row scan error: %v", err)
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

		// Add Reservation to the map
		resList[r.ReservationId] = r
	}

	// Check for errors during iteration
	if err := rows.Err(); err != nil {
		log.Printf("Rows iteration error: %v", err)
		return nil, err
	}

	return resList, nil
}

func (s *ReserveService) GetAvailableTimeSlots(carId int, year int, month int, day int) ([]AvailableSlot, error) {

	// Define available time range for the entire day
	availableStart := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	availableEnd := time.Date(year, time.Month(month), day, 23, 59, 0, 0, time.UTC)

	// Prepare date for query
	queryDate := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)

	// Get all reservations for the car on this particular day
	query := `SELECT start_datetime, end_datetime 
              FROM reservations 
              WHERE car_id = ? AND DATE(start_datetime) = ?
              ORDER BY start_datetime`

	rows, err := s.DB.Query(query, carId, queryDate.Format("2006-01-02"))
	if err != nil {
		return nil, fmt.Errorf("failed to query reservations: %v", err)
	}
	defer rows.Close()

	// Collect reserved time slots
	var reservedSlots []struct {
		start time.Time
		end   time.Time
	}

	for rows.Next() {
		var startByte, endByte []byte // Scan into []byte first
		if err := rows.Scan(&startByte, &endByte); err != nil {
			return nil, fmt.Errorf("error scanning reservation: %v", err)
		}

		// Parse the byte slices into time.Time
		start, err := time.Parse("2006-01-02 15:04:05", string(startByte))
		if err != nil {
			return nil, fmt.Errorf("error parsing start time: %v", err)
		}
		end, err := time.Parse("2006-01-02 15:04:05", string(endByte))
		if err != nil {
			return nil, fmt.Errorf("error parsing end time: %v", err)
		}

		reservedSlots = append(reservedSlots, struct {
			start time.Time
			end   time.Time
		}{start, end})
	}

	// Check for any errors encountered during iteration
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error reading reservation rows: %v", err)
	}

	// Generate available time slots with detailed breakdown
	return generateAvailableSlots(availableStart, availableEnd, reservedSlots), nil
}

func generateAvailableSlots(availableStart, availableEnd time.Time, reservedSlots []struct {
	start time.Time
	end   time.Time
}) []AvailableSlot {
	var availableSlots []AvailableSlot
	currentStart := availableStart

	// If no reservations, return the entire day as available
	if len(reservedSlots) == 0 {
		return []AvailableSlot{
			{
				StartTime: availableStart.Format("15:04"),
				EndTime:   availableEnd.Format("15:04"),
			},
		}
	}

	// Sort reserved slots by start time (assumed to be done in the query)
	for _, reserved := range reservedSlots {
		// Add available slot before the reservation
		if currentStart.Before(reserved.start) {
			availableSlots = append(availableSlots, AvailableSlot{
				StartTime: currentStart.Format("15:04"),
				EndTime:   reserved.start.Format("15:04"),
			})
		}

		// Move current start to the end of this reserved slot
		currentStart = max(currentStart, reserved.end)
	}

	// Add final slot if there's remaining time after last reservation
	if currentStart.Before(availableEnd) {
		availableSlots = append(availableSlots, AvailableSlot{
			StartTime: currentStart.Format("15:04"),
			EndTime:   availableEnd.Format("15:04"),
		})
	}

	return availableSlots
}

// Helper function to get max of two times
func max(a, b time.Time) time.Time {
	if a.After(b) {
		return a
	}
	return b
}

// fetches Reserve details for the profile page
func (s *ReserveService) GetAllReservations(id int) (map[int]models.Reservation, error) {
	// create dictionary to store Reservations
	resList := make(map[int]models.Reservation)

	// Get Reservations
	query := "SELECT reservation_id, start_datetime, end_datetime FROM reservations WHERE user_id = ?"
	rows, err := s.DB.Query(query, id)
	if err != nil {
		log.Printf("Query error: %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		r := models.Reservation{}
		var start, end sql.NullString

		if err := rows.Scan(&r.ReservationId, &start, &end); err != nil {
			log.Printf("Row scan error: %v", err)
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

		// Add Reservation to the map
		resList[r.ReservationId] = r
	}

	// Check for errors during iteration
	if err := rows.Err(); err != nil {
		log.Printf("Rows iteration error: %v", err)
		return nil, err
	}

	return resList, nil
}

// fetches Reserve details for the Reserves page
func (s *ReserveService) GetUserReservations(id int) (map[int]models.Reservation, error) {
	// create a dictionary to store Reservations
	resList := make(map[int]models.Reservation)

	// Get Reservations query
	query := "SELECT reservation_id, car_id, start_datetime, end_datetime, status FROM reservations WHERE user_id = ? AND status NOT IN ('Cancelled', 'Completed') ORDER BY start_datetime DESC"
	rows, err := s.DB.Query(query, id)
	if err != nil {
		log.Printf("Query error: %v", err)
		return nil, err
	}
	defer rows.Close()

	// Iterate over each row
	for rows.Next() {
		r := models.Reservation{}
		var start, end sql.NullString

		// Scan the reservation data
		if err := rows.Scan(&r.ReservationId, &r.CarId, &start, &end, &r.Status); err != nil {
			log.Printf("Row scan error: %v", err)
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

		// Create a new Car object to store car details
		c := &models.Car{}

		// Get car details based on the car_id
		carQuery := "SELECT car_model, license_plate, current_location, charge_level, rate FROM cars WHERE car_id = ?"
		err := s.DB.QueryRow(carQuery, r.CarId).Scan(&c.CarModel, &c.LiscencePlate, &c.CurrLoc, &c.Charge, &c.Rate)
		if err != nil {
			log.Println("Row scan error for car details:", err)
			return nil, err
		}

		// Attach car details to the reservation
		c.CarId = r.CarId
		r.CarDetails = c

		// Add Reservation to the map
		resList[r.ReservationId] = r
	}

	// Check for any errors during rows iteration
	if err := rows.Err(); err != nil {
		log.Printf("Rows iteration error: %v", err)
		return nil, err
	}

	// Return the map containing reservations
	return resList, nil
}

// create a reservation
func (s *ReserveService) CreateReservation(res *models.Reservation) (*models.Reservation, error) {

	// Prepare the SQL INSERT statement
	query := "INSERT INTO reservations (user_id, car_id, start_datetime, end_datetime, status) VALUES (?, ?, ?, ?, ?)"
	result, err := s.DB.Exec(query, res.UserId, res.CarId, res.Start, res.End, res.Status)
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
	res.ReservationId = int(lastInsertID)

	return res, nil
}

// update an existing reservation
func (s *ReserveService) UpdateReservationStatus(res *models.Reservation) (*models.Reservation, error) {

	// Prepare the SQL INSERT statement
	query := "UPDATE reservations SET Status = ? WHERE reservation_id = ?"
	_, err := s.DB.Exec(query, res.Status, res.ReservationId)
	if err != nil {
		log.Println("Database update error:", err)
		return nil, err
	}

	return res, nil
}

// GetReserveDetails retrieves details of a reservation by id
func (s *ReserveService) GetReservationDetails(id int) (*models.Reservation, error) {
	// Declare variables
	var r models.Reservation
	var start, end sql.NullString

	// Query to fetch reservation details
	query := "SELECT reservation_id, car_id, start_datetime, end_datetime, status FROM reservations WHERE reservation_id = ?"
	err := s.DB.QueryRow(query, id).Scan(&r.ReservationId, &r.CarId, &start, &end, &r.Status)

	// Handle query error (e.g., no rows found)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("No reservation found with ID:", id)
			return nil, nil // or return an error depending on your requirement
		}
		log.Println("Query error:", err)
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

	return &r, nil
}

// fetches Reserve details for the Reserves page
func (s *ReserveService) GetCompletedReservations(id int) (map[int]models.Reservation, error) {
	// create a dictionary to store Reservations
	resList := make(map[int]models.Reservation)

	// Get Reservations query
	query := "SELECT reservation_id, car_id, start_datetime, end_datetime, status FROM reservations WHERE user_id = ? AND status = 'Completed' ORDER BY start_datetime DESC"
	rows, err := s.DB.Query(query, id)
	if err != nil {
		log.Printf("Query error: %v", err)
		return nil, err
	}
	defer rows.Close()

	// Iterate over each row
	for rows.Next() {
		r := models.Reservation{}
		var start, end sql.NullString

		// Scan the reservation data
		if err := rows.Scan(&r.ReservationId, &r.CarId, &start, &end, &r.Status); err != nil {
			log.Printf("Row scan error: %v", err)
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

		// Create a new Car object to store car details
		c := &models.Car{}

		// Get car details based on the car_id
		carQuery := "SELECT car_model, license_plate, current_location, charge_level, rate FROM cars WHERE car_id = ?"
		err := s.DB.QueryRow(carQuery, r.CarId).Scan(&c.CarModel, &c.LiscencePlate, &c.CurrLoc, &c.Charge, &c.Rate)
		if err != nil {
			log.Println("Row scan error for car details:", err)
			return nil, err
		}

		// Attach car details to the reservation
		c.CarId = r.CarId
		r.CarDetails = c

		// Add Reservation to the map
		resList[r.ReservationId] = r
	}

	// Check for any errors during rows iteration
	if err := rows.Err(); err != nil {
		log.Printf("Rows iteration error: %v", err)
		return nil, err
	}

	// Return the map containing reservations
	return resList, nil
}

func (s *ReserveService) EndReservation(res *models.Reservation) (*models.Reservation, error) {

	// Ensure we are using the correct time format, in UTC if needed
	endTime := time.Now()

	// Prepare the SQL UPDATE statement
	query := "UPDATE reservations SET Status = ?, end_datetime = ? WHERE reservation_id = ?"
	_, err := s.DB.Exec(query, res.Status, endTime, res.ReservationId)
	if err != nil {
		log.Println("Database update error:", err)
		return nil, err
	}

	return res, nil
}

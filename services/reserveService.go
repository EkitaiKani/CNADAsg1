package services

import (
	"database/sql"
	"log"
	"time"

	"CNADASG1/models"
)

type ReserveService struct {
	DB *sql.DB
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


// fetches Reserve details for the Reserves page
func (s *ReserveService) GetUserReservations(id int) (map[int]models.Reservation, error) {
	// create dictionary to store Reservations
	resList := make(map[int]models.Reservation)

	// Get Reservations
	query := "SELECT reservation_id, start_datetime, end_datetime FROM reservations WHERE user_id = ? AND status NOT IN ('Cancelled', 'Completed')"
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


// create a new reservation
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
func (s *ReserveService) UpdateReservation(res *models.Reservation) (*models.Reservation, error) {

	// Prepare the SQL INSERT statement
	query := "UPDATE reservation SET Status = ?"
	_, err := s.DB.Exec(query, res.Status)
	if err != nil {
		log.Println("Database update error:", err)
		return nil, err
	}

	return res, nil
}

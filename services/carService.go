package services

import (
	"database/sql"
	"log"
	"time"

	"CNADASG1/models"
)

type CarService struct {
	DB *sql.DB
}

// fetches Car details for the Cars page
func (s *CarService) GetCars() (map[int]models.Car, error) {

	// create dictionary to store cars
	carList := make(map[int]models.Car)

	// Get cars
	query := "SELECT car_id, car_model FROM cars"
	rows, err := s.DB.Query(query)
	if err != nil {
		log.Println("Row scan error:", err)
		return nil, err
	}

	for rows.Next() {
		var carId int
		var carModel string

		if err := rows.Scan(&carId, &carModel); err != nil {
			log.Println("Row scan error:", err)
			return nil, err
		}

		// Add car to the map
		carList[carId] = models.Car{
			CarId:    carId,
			CarModel: carModel,
		}
	}

	// Check for errors during iteration
	if err := rows.Err(); err != nil {
		log.Println("Rows iteration error:", err)
		return nil, err
	}

	return carList, nil
}

// GetCarDetails retrieves details of a car by its ID
func (s *CarService) GetCarDetails(id int) (*models.Car, error) {
	// Create a new instance of models.Car
	u := &models.Car{}

	// Set car ID
	u.CarId = id

	// set variables for dates
	var lastServiced string

	// Get car details from the database
	query := "SELECT car_model, current_location, charge_level, cleanliness_status, last_serviced, rate FROM cars WHERE car_id = ?"
	err := s.DB.QueryRow(query, id).Scan(
		&u.CarModel,    // car_model
		&u.CurrLoc,     // current_location
		&u.Charge,      // charge_level
		&u.Cleanliness, // cleanliness_status
		&lastServiced,  // last_serviced (sql.NullTime)
		&u.Rate,        // rate
	)

	// Parse the date into a time.Time if it's not NULL
	if lastServiced != "" {
		parsedTime, err := time.Parse("2006-01-02 15:04:05", lastServiced)
		if err != nil {
			log.Fatal(err)
		}
		u.LastServiced = sql.NullTime{Time: parsedTime, Valid: true}
	} else {
		u.LastServiced = sql.NullTime{Valid: false} // Handle NULL case
	}

	// Handle scan error
	if err != nil {
		if err == sql.ErrNoRows {
			// No rows found for the given ID
			log.Println("No car found with ID:", id)
		} else {
			// Any other error
			log.Println("Row scan error:", err)
		}
		return nil, err
	}

	return u, nil
}

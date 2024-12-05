package apis

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"CNADASG1/services"

	"github.com/gorilla/mux"
)

type CarAPI struct {
	Service *services.CarService
}

func (h *CarAPI) Cars(w http.ResponseWriter, r *http.Request) {

	// get car details
	carList, err := h.Service.GetCars()
	jsonBody := make(map[string]interface{})

	if err != nil {
		// Log the actual error
		jsonBody = map[string]interface{}{
			"message": "Error getting cars, please try again",
			"error":   true,
		}
		log.Print("Internal server error:", err)
	} else if len(carList) == 0 {
		// if there are no cars
		jsonBody = map[string]interface{}{
			"message": "There are no available cars",
			"error":   false,
		}
	} else {
		// Render cars
		jsonBody = map[string]interface{}{
			"carList": carList,
			"error":   false,
		}
	}

	// Secure HTTP headers
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("X-XSS-Protection", "1; mode=block")

	// Encode the response and handle errors
	if err := json.NewEncoder(w).Encode(jsonBody); err != nil {
		log.Println("JSON encoding error:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

}

func (h *CarAPI) CarDetails(w http.ResponseWriter, r *http.Request) {
	// Get car id from URL
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	// Handle error if converting id fails
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Get car details
	car, err := h.Service.GetCarDetails(id)
	jsonBody := make(map[string]interface{})

	// log.Print(car.CarModel)

	if err != nil {
		jsonBody = map[string]interface{}{
			"message": "Error getting car details, please try again",
			"error":   true,
		}
		log.Print("Internal server error:", err)
	} else {
		// Render car details
		jsonBody = map[string]interface{}{
			"car":     car,
			"error":   false,
			"message": "",
		}

	}

	// Secure HTTP headers
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("X-XSS-Protection", "1; mode=block")

	// Encode the response and handle errors
	if err := json.NewEncoder(w).Encode(jsonBody); err != nil {
		log.Println("JSON encoding error:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

}

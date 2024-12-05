package apis

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"CNADASG1/models"
	"CNADASG1/services"

	"github.com/gorilla/mux"
)

type ReserveAPI struct {
	Service *services.ReserveService
}

func (h *ReserveAPI) CarReservations(w http.ResponseWriter, r *http.Request) {

	// Get id from URL
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// get car's reservation details
	resList, err := h.Service.GetCarReservations(id)
	jsonBody := make(map[string]interface{})

	if err != nil {
		// Log the actual error
		jsonBody = map[string]interface{}{
			"message": "Error getting cars, please try again",
			"error":   true,
		}
		log.Print("Internal server error:", err)

	} else if len(resList) == 0 {
		// if there are no cars
		jsonBody = map[string]interface{}{
			"message": "There are no reservations for this car.",
			"error":   false,
		}

	} else { // Render cars
		jsonBody = map[string]interface{}{
			"resList": resList,
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

func (h *ReserveAPI) UserReservations(w http.ResponseWriter, r *http.Request) {
	// Get id from URL
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// get car's reservation details
	resList, err := h.Service.GetUserReservations(id)
	jsonBody := make(map[string]interface{})

	if err != nil {
		// Log the actual error
		jsonBody = map[string]interface{}{
			"message": "Error getting cars, please try again",
			"error":   true,
		}
		log.Print("Internal server error:", err)

	} else if len(resList) == 0 {
		// if there are no cars
		jsonBody = map[string]interface{}{
			"message": "You have not made any reservations.",
			"error":   false,
		}

	} else { // Render cars
		jsonBody = map[string]interface{}{
			"resList": resList,
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

func (h *ReserveAPI) CreateReservation(w http.ResponseWriter, r *http.Request) {
	// Decode the incoming JSON request body
	var res *models.Reservation
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&res); err != nil {
		log.Println("JSON decoding error:", err)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Try to create the user using the service
	createdRes, err := h.Service.CreateReservation(res)
	jsonBody := make(map[string]interface{})
	if err != nil {
		jsonBody = map[string]interface{}{
			"message": "Account was not created. Check your input fields and try again.",
			"error":   true,
		}
		log.Print("Internal server error:", err)
	} else {
		// After successful creation, render the user in the template
		jsonBody = map[string]interface{}{
			"res":     createdRes,
			"message": "Account created successfully",
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

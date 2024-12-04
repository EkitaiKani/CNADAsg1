package apis

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"CNADASG1/services"

	"github.com/gorilla/mux"
)

type ReserveAPI struct {
	Templates *template.Template
	Service   *services.ReserveService
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

func (h *ReserveAPI) UpdateReservationStatus(w http.ResponseWriter, r *http.Request) {

}

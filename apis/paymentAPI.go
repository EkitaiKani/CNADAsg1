package apis

import (
	"CNADASG1/models"
	"CNADASG1/services"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type PaymentAPI struct {
	Service *services.PaymentService
}

func (h *PaymentAPI) CreatePayment(w http.ResponseWriter, r *http.Request) {

	// Decode the incoming JSON request body
	var u *models.Payment
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		log.Println("JSON decoding error:", err)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	//log.Print(u)

	var err error
	//calculate amount
	u, err = h.Service.CalculatePayment(*u)
	// log.Print(amount)

	// Try to create the user using the service
	createdPay, err := h.Service.CreatePayment(u)
	jsonBody := make(map[string]interface{})
	if err != nil {
		jsonBody = map[string]interface{}{
			"message": "Payment was not created. Check your input fields and try again.",
			"error":   true,
		}
		log.Print("Internal server error:", err)
	} else {
		// After successful creation, render the user in the template
		jsonBody = map[string]interface{}{
			"pay":     createdPay,
			"message": "Payment created successfully",
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

func (h *PaymentAPI) PaymentDetails(w http.ResponseWriter, r *http.Request) {
	// Get payment id from URL
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	// Handle error if converting id fails
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Get payment details
	pay, err := h.Service.GetPayment(id)
	pay, err = h.Service.CalculatePayment(*pay)

	jsonBody := make(map[string]interface{})

	if err != nil {
		jsonBody = map[string]interface{}{
			"message": "Error getting payment details, please try again",
			"error":   true,
		}
		log.Print("Internal server error:", err)
	} else {
		// Render car details
		jsonBody = map[string]interface{}{
			"payment": pay,
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

func (h *PaymentAPI) CompletePayment(w http.ResponseWriter, r *http.Request) {

	// Decode the incoming JSON request body
	var pay *models.Payment
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&pay); err != nil {
		log.Println("JSON decoding error:", err)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Try to update using the service
	updatedPay, err := h.Service.MakePayment(*pay)
	// log.Print(updatedRes)

	jsonBody := make(map[string]interface{})
	if err != nil {
		jsonBody = map[string]interface{}{
			"message": "Payment was not updated. Please try again.",
			"error":   true,
		}
		log.Print("Internal server error:", err)
	} else {
		// After successful creation, render the user in the template
		jsonBody = map[string]interface{}{
			"res":     updatedPay,
			"message": "Payment updated successfully",
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

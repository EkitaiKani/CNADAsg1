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

type login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserAPI struct {
	Service *services.UserService
}

func (h *UserAPI) RegisterUser(w http.ResponseWriter, r *http.Request) {
	// Get id from URL
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Decode the incoming JSON request body
	var u *models.User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		log.Println("JSON decoding error:", err)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	u.UserId = id

	// Try to create the user using the service
	createdUser, err := h.Service.CreateUser(u)
	jsonBody := make(map[string]interface{})
	if err != nil {
		jsonBody = map[string]interface{}{
			"message": "Account was not created. Check your input fields and try again.",
			"error":   true,
		}
	} else {
		// After successful creation, render the user in the template
		jsonBody = map[string]interface{}{
			"user":    createdUser,
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

func (h *UserAPI) LoginUser(w http.ResponseWriter, r *http.Request) {

	// Decode the incoming JSON request body
	var u login
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		log.Println("JSON decoding error:", err)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	user, err := h.Service.LogInUser(u.Username, u.Password)

	jsonBody := make(map[string]interface{})

	if err != nil {
		// error message
		jsonBody = map[string]interface{}{
			"message": "Username or password is incorrect.",
			"error":   true,
		}
		return

	} else {
		// After successful creation, render the user in the template
		jsonBody = map[string]interface{}{
			"user":    user,
			"message": "Logged in successfully",
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

func (h *UserAPI) UserDetails(w http.ResponseWriter, r *http.Request) {
	// Get id from URL
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// get user details
	user, err := h.Service.GetUserDetails(id)
	// log.Print(user.UserEmail)
	jsonBody := make(map[string]interface{})

	if err != nil {
		// Log the actual error
		jsonBody = map[string]interface{}{
			"message": "Error getting user details, please try again",
			"error":   true,
		}
		return
	}

	// Render user details
	jsonBody = map[string]interface{}{
		"user":  user,
		"error": false,
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

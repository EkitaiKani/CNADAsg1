package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"CNADASG1/models"
	"CNADASG1/templates"
)

type ReserveHandler struct {
	Templates *template.Template
	BaseURL   string
}

// NewResHandler is a constructor function to create a new ResHandler with the API base URL
func NewResHandler(baseURL string) *CarHandler {
	return &CarHandler{
		BaseURL: baseURL, // Store the base URL
	}
}

func (h *ReserveHandler) CarReservations(w http.ResponseWriter, r *http.Request) {
	// get id of car
	id := r.URL.Query().Get("id")

	// get car's reservation details
	var response map[string]interface{}
	url := h.BaseURL + "car/" + id
	client := &http.Client{}

	if req, err := http.NewRequest("GET", url, nil); err == nil {
		if res, err := client.Do(req); err == nil {

			body, err := ioutil.ReadAll(res.Body)

			if err != nil {
				log.Print("An error occured")
			}

			// unmarshal response data
			err = json.Unmarshal(body, &response)

		}
	}

	//log.Print(response)

	// Render cars
	if err := templates.Templates.ExecuteTemplate(w, "reservation.html", response); err != nil {
		http.Error(w, "Template render error", http.StatusInternalServerError)
	}

}

func (h *ReserveHandler) UserReservations(w http.ResponseWriter, r *http.Request) {
	// get id of car
	id := r.FormValue("userid")

	// get car's reservation details
	var response map[string]interface{}
	url := h.BaseURL + "/user/" + id
	client := &http.Client{}

	if req, err := http.NewRequest("GET", url, nil); err == nil {
		if res, err := client.Do(req); err == nil {

			body, err := ioutil.ReadAll(res.Body)

			if err != nil {
				log.Print("An error occured")
			}

			// unmarshal response data
			err = json.Unmarshal(body, &response)

		}
	}
	log.Print(response)

	// Render cars
	if err := templates.Templates.ExecuteTemplate(w, "reservation.html", response); err != nil {
		http.Error(w, "Template render error", http.StatusInternalServerError)
	}

}

func (h *ReserveHandler) PostReservation(w http.ResponseWriter, r *http.Request) {
	// If the method is POST, handle form submission
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	session, err := store.Get(r, "user-session")
	userID, ok := session.Values["user_id"].(string)
	if !ok {
		// If the user_id is not found or has an incorrect type
		// log.Print(session.Values["user_id"]) // For debugging
		http.Error(w, "User not logged in", http.StatusUnauthorized)
		return
	}

	// Retrieve form values by their "name" attribute
	userId, _ := strconv.Atoi(userID)
	carIdStr := r.FormValue("CarId")
	startStr := r.FormValue("Start")
	endStr := r.FormValue("End")
	dateStr := r.FormValue("date")

	//log.Print(startStr)
	//log.Print(dateStr)

	// Convert CarId to int
	carId, err := strconv.Atoi(carIdStr)
	if err != nil {
		// Handle the error if conversion fails
		log.Printf("Error converting CarId: %v", err)
		http.Error(w, "Invalid CarId", http.StatusBadRequest)
		return
	}

	// Parse the date string to a time.Time (we assume it's in UTC)
	baseTime, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return
	}

	// DO NOT REMOVE
	baseTime = baseTime.AddDate(0, 0, 1)

	// Split the time string into hours and minutes and convert to a time.Duration
	hours, minutes := 0, 0
	_, err = fmt.Sscanf(startStr, "%d:%d", &hours, &minutes)
	if err != nil {
		fmt.Println("Error parsing time:", err)
		return
	}

	// Set the hours and minutes from the timeStr into the baseTime
	startTime := time.Date(baseTime.Year(), baseTime.Month(), baseTime.Day(), hours, minutes, 0, 0, time.UTC)

	startDate := sql.NullTime{
		Time:  startTime, // This is your combined date and time
		Valid: true,      // Set it to true if it's a valid date
	}

	_, err = fmt.Sscanf(endStr, "%d:%d", &hours, &minutes)
	if err != nil {
		fmt.Println("Error parsing time:", err)
		return
	}

	// Set the hours and minutes from the timeStr into the baseTime
	endTime := time.Date(baseTime.Year(), baseTime.Month(), baseTime.Day(), hours, minutes, 0, 0, time.UTC)

	endDate := sql.NullTime{
		Time:  endTime, // This is your combined date and time
		Valid: true,    // Set it to true if it's a valid date
	}
	res := &models.Reservation{
		CarId:  carId,
		UserId: userId,
		Start:  endDate,
		End:    startDate,
	}

	var response map[string]interface{}
	url := h.BaseURL + "reservation/"
	client := &http.Client{}
	postBody, _ := json.Marshal(res)
	resBody := bytes.NewBuffer(postBody)

	if req, err := http.NewRequest("POST", url, resBody); err == nil {
		if res, err := client.Do(req); err == nil {
			// You can log the status code here if necessary
			body, err := ioutil.ReadAll(res.Body)

			if err != nil {
				log.Print("An error occured")
			}

			// unmarshal response data
			err = json.Unmarshal(body, &response)

		}
	}

	renderErr := templates.Templates.ExecuteTemplate(w, "reservation.html", response)
	if renderErr != nil {
		http.Error(w, renderErr.Error(), http.StatusInternalServerError)
	}
}

func (h *ReserveHandler) UpdateReservation(w http.ResponseWriter, r *http.Request) {
	// If the method is POST, handle form submission
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	session, err := store.Get(r, "user-session")
	userID, ok := session.Values["user_id"].(int)
	if !ok {
		// If the user_id is not found or has an incorrect type
		// log.Print(session.Values["user_id"]) // For debugging
		http.Error(w, "User not logged in", http.StatusUnauthorized)
		return
	}

	// Retrieve form values by their "name" attribute
	resIdStr := r.FormValue("ResId")
	userId := userID
	carIdStr := r.FormValue("CarId")
	startStr := r.FormValue("Start")
	endStr := r.FormValue("End")

	log.Print(startStr)

	// Convert UserId to int
	resId, err := strconv.Atoi(resIdStr)
	if err != nil {
		// Handle the error if conversion fails
		log.Printf("Error converting UserId: %v", err)
		http.Error(w, "Invalid UserId", http.StatusBadRequest)
		return
	}

	// Convert CarId to int
	carId, err := strconv.Atoi(carIdStr)
	if err != nil {
		// Handle the error if conversion fails
		log.Printf("Error converting CarId: %v", err)
		http.Error(w, "Invalid CarId", http.StatusBadRequest)
		return
	}

	// Try to parse the date string to time.Time
	var start sql.NullTime
	var end sql.NullTime

	if startStr != "" {
		parsedTime, err := time.Parse("2006-01-02 15:04:05", startStr) // Expecting date in "YYYY-MM-DD" format
		if err != nil {
			log.Println("Error parsing start time:", err)
			http.Error(w, "Invalid date format", http.StatusBadRequest)
			return
		}
		start = sql.NullTime{Time: parsedTime, Valid: true} // Mark as valid with parsed time
	} else {
		start = sql.NullTime{Valid: false} // If empty, mark as invalid (representing NULL in SQL)
	}

	if endStr != "" {
		parsedTime, err := time.Parse("2006-01-02 15:04:05", endStr) // Expecting date in "YYYY-MM-DD" format
		if err != nil {
			log.Println("Error parsing end time:", err)
			http.Error(w, "Invalid date format", http.StatusBadRequest)
			return
		}
		end = sql.NullTime{Time: parsedTime, Valid: true} // Mark as valid with parsed time
	} else {
		end = sql.NullTime{Valid: false} // If empty, mark as invalid (representing NULL in SQL)
	}

	res := &models.Reservation{
		ReservationId: resId,

		CarId:  carId,
		UserId: userId,
		Start:  start,
		End:    end,
	}

	var response map[string]interface{}
	url := h.BaseURL + "reservation/" + resIdStr
	client := &http.Client{}
	postBody, _ := json.Marshal(res)
	resBody := bytes.NewBuffer(postBody)

	if req, err := http.NewRequest("PUT", url, resBody); err == nil {
		if res, err := client.Do(req); err == nil {
			// You can log the status code here if necessary
			body, err := ioutil.ReadAll(res.Body)

			if err != nil {
				log.Print("An error occured")
			}

			// unmarshal response data
			err = json.Unmarshal(body, &response)

		}
	}

	renderErr := templates.Templates.ExecuteTemplate(w, "reservation.html", response)
	if renderErr != nil {
		http.Error(w, renderErr.Error(), http.StatusInternalServerError)
	}
}

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

	"github.com/gorilla/mux"
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
	session, err := store.Get(r, "user-session")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Retrieve the user ID as a string from session
	userID, ok := session.Values["user_id"].(string)
	if !ok {
		// If the user_id is not found or has an incorrect type
		// log.Print(session.Values["user_id"]) // For debugging
		http.Error(w, "User not logged in", http.StatusUnauthorized)
		return
	}

	// get reservation details
	var response map[string]interface{}
	url := h.BaseURL + "reservation/user/" + userID
	// log.Print(url)
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

	// Render cars
	if err := templates.Templates.ExecuteTemplate(w, "userReservations.html", response); err != nil {
		http.Error(w, "Template render error", http.StatusInternalServerError)
	}

}

func (h *ReserveHandler) AllReservations(w http.ResponseWriter, r *http.Request) {
	// get id of car
	id := r.FormValue("userid")

	// get car's reservation details
	var response map[string]interface{}
	url := h.BaseURL + "reservation/user/all/" + id
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
	// log.Print(response)

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
		Start:  startDate,
		End:    endDate,
		Status: "Pending",
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

func (h *ReserveHandler) CancelReservation(w http.ResponseWriter, r *http.Request) {

	// If the method is POST, handle form submission
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get res id from URL
	vars := mux.Vars(r)
	resIdStr := vars["id"]
	id, err := strconv.Atoi(resIdStr)

	// Handle error if converting id fails
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	res := &models.Reservation{
		ReservationId: id,
		Status:        "Cancelled",
	}

	var response map[string]interface{}
	url := h.BaseURL + "reservation/update/" + resIdStr
	// log.Print(url)

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

	log.Print(response["error"])
	log.Print(response["message"])

	http.Redirect(w, r, "/reserve/user", http.StatusSeeOther)
}

func (h *ReserveHandler) ReserveNow(w http.ResponseWriter, r *http.Request) {
	// If the method is POST, handle form submission
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
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

	// Get car id from URL
	vars := mux.Vars(r)
	CaridStr := vars["id"]
	id, err := strconv.Atoi(CaridStr)
	// Handle error if converting id fails
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	userId, _ := strconv.Atoi(userID)

	// Set the hours and minutes from the timeStr into the baseTime
	startTime := time.Now()
	startDate := sql.NullTime{
		Time:  startTime, // This is your combined date and time
		Valid: true,      // Set it to true if it's a valid date
	}

	res := &models.Reservation{
		CarId:  id,
		UserId: userId,
		Start:  startDate,
		Status: "Confirmed",
	}

	var resResponse map[string]interface{}
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
			err = json.Unmarshal(body, &resResponse)

		}
	}

	// update the car status
	car := &models.Car{
		Status: "Reserved",
	}

	var carResponse map[string]interface{}
	url = h.BaseURL + "car/" + CaridStr
	postBody, _ = json.Marshal(car)
	resBody = bytes.NewBuffer(postBody)

	if req, err := http.NewRequest("PUT", url, resBody); err == nil {
		if res, err := client.Do(req); err == nil {
			// You can log the status code here if necessary
			body, err := ioutil.ReadAll(res.Body)

			if err != nil {
				log.Print("An error occured")
			}

			// unmarshal response data
			err = json.Unmarshal(body, &carResponse)

		}
	}
	// Redirect to the user reservation page after the operation
	http.Redirect(w, r, "/reserve/user", http.StatusSeeOther)
}

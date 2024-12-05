package handlers

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"CNADASG1/templates"
)

type CarHandler struct {
	Templates *template.Template
	BaseURL   string
}

// NewCarHandler is a constructor function to create a new CarHandler with the API base URL
func NewCarHandler(baseURL string) *CarHandler {
	return &CarHandler{
		BaseURL: baseURL, // Store the base URL
	}
}

func (h *CarHandler) Cars(w http.ResponseWriter, r *http.Request) {

	var response map[string]interface{}
	url := h.BaseURL + "car/"
	client := &http.Client{}

	if req, err := http.NewRequest("GET", url, nil); err == nil {
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

	//log.Print(response)

	// Render cars
	if err := templates.Templates.ExecuteTemplate(w, "cars.html", response); err != nil {
		http.Error(w, "Template render error", http.StatusInternalServerError)
	}

}

func (h *CarHandler) CarDetails(w http.ResponseWriter, r *http.Request) {
	// Get the car ID from the query parameter
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Car ID is required", http.StatusBadRequest)
		return
	}

	var response map[string]interface{}
	url := h.BaseURL + "car/" + id
	client := &http.Client{}

	if req, err := http.NewRequest("GET", url, nil); err == nil {
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

	//log.Print(response)

	// Render car details
	err := templates.Templates.ExecuteTemplate(w, "carDetails.html", response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Handler for listing cars or showing car details
func (h *CarHandler) HandleCars(w http.ResponseWriter, r *http.Request) {
	// Check if 'id' query parameter is present
	id := r.URL.Query().Get("id")

	if id != "" {
		// If 'id' is present, show details of the specific car
		h.CarDetails(w, r)
		return
	}

	// If 'id' is not present, list all cars
	h.Cars(w, r)
}

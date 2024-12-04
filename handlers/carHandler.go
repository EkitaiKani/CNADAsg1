package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"CNADASG1/templates"

	"github.com/gorilla/mux"
)

type CarHandler struct {
	BaseURL string
}

// NewCarHandler is a constructor function to create a new CarHandler with the API base URL
func NewCarHandler(baseURL string) *CarHandler {
	return &CarHandler{
		BaseURL: baseURL + "/Cars/", // Store the base URL
	}
}

func (h *CarHandler) Cars(w http.ResponseWriter, r *http.Request) {

	var response map[string]interface{}
	url := h.BaseURL
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

	// Render cars
	if err := templates.Templates.ExecuteTemplate(w, "cars.html", response); err != nil {
		http.Error(w, "Template render error", http.StatusInternalServerError)
	}

}

func (h *CarHandler) CarDetails(w http.ResponseWriter, r *http.Request) {
	// Get car id from URL
	vars := mux.Vars(r)
	id := vars["id"]

	var response map[string]interface{}
	url := h.BaseURL + id
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

	// Render car details
	err := templates.Templates.ExecuteTemplate(w, "carDetails.html", response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

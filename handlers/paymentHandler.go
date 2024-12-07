package handlers

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"CNADASG1/templates"

	"github.com/gorilla/mux"
)

type PaymentHandler struct {
	Templates *template.Template
	BaseURL   string
}

// NewPaymentHandler is a constructor function to create a new PaymentHandler with the API base URL
func NewPaymentHandler(baseURL string) *PaymentHandler {
	return &PaymentHandler{
		BaseURL: baseURL, // Store the base URL
	}
}

func (h *PaymentHandler) Payment(w http.ResponseWriter, r *http.Request) {

	// Get pay id from URL
	vars := mux.Vars(r)
	payIdStr := vars["id"]

	var response map[string]interface{}
	url := h.BaseURL + "payment/" + payIdStr
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

	// Render Payments
	if err := templates.Templates.ExecuteTemplate(w, "payment.html", response); err != nil {
		http.Error(w, "Template render error", http.StatusInternalServerError)
	}

}

package handlers

import (
	"html/template"
	"net/http"
	"strconv"

	"CNADASG1/services"
	"CNADASG1/templates"

	"github.com/gorilla/mux"
)

type CarHandler struct {
	Templates *template.Template
	Service   *services.CarService
}

func (h *CarHandler) Cars(w http.ResponseWriter, r *http.Request) {

	// get car details
	carList, err := h.Service.GetCars()

	if err != nil {
		// Log the actual error
		renderErr := templates.Templates.ExecuteTemplate(w, "cars.html", map[string]interface{}{
			"message": "Error getting cars, please try again",
			"error":   true,
		})
		if renderErr != nil {
			http.Error(w, "Template render error", http.StatusInternalServerError)
		}
		return
	}

	if len(carList) == 0 {
		// if there are no cars
		if err := templates.Templates.ExecuteTemplate(w, "cars.html", map[string]interface{}{
			"message": "There are no available cars",
			"error":   false,
		}); err != nil {
			http.Error(w, "Template render error", http.StatusInternalServerError)
		}
		return
	}

	// Render cars
	if err := templates.Templates.ExecuteTemplate(w, "cars.html", map[string]interface{}{
		"carList": carList,
		"error":   false,
	}); err != nil {
		http.Error(w, "Template render error", http.StatusInternalServerError)
	}

}

func (h *CarHandler) CarDetails(w http.ResponseWriter, r *http.Request) {
	// Get car id from URL
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	// Handle error if converting id fails
	if err != nil {
		renderErr := templates.Templates.ExecuteTemplate(w, "carDetails.html", map[string]interface{}{
			"message": "Invalid car ID",
			"error":   true,
		})
		if renderErr != nil {
			http.Error(w, renderErr.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Get car details
	car, err := h.Service.GetCarDetails(id)
	// log.Print(car.CarModel)

	if err != nil {
		renderErr := templates.Templates.ExecuteTemplate(w, "carDetails.html", map[string]interface{}{
			"message": "Error getting car details, please try again",
			"error":   true,
		})
		if renderErr != nil {
			http.Error(w, renderErr.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Render car details
	err = templates.Templates.ExecuteTemplate(w, "carDetails.html", map[string]interface{}{
		"car":     car,
		"error":   false,
		"message": "",
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

package main

import (
	"log"
	"net/http"

	"CNADASG1/apis"
	"CNADASG1/config"
	"CNADASG1/services"

	"github.com/gorilla/mux"
)

func main() {

	// Connect to database
	db := config.ConnectDatabase()
	defer db.Close()

	// Initialize services
	userService := &services.UserService{DB: db}
	carService := &services.CarService{DB: db}
	resService := &services.ReserveService{DB: db}

	// Initialize apis
	userAPI := &apis.UserAPI{Service: userService}
	carAPI := &apis.CarAPI{Service: carService}
	resAPI := &apis.ReserveAPI{Service: resService}

	// Create router
	r := mux.NewRouter()

	// User routes
	r.HandleFunc("/api/v1/user/{id}", userAPI.UserDetails).Methods("GET")
	r.HandleFunc("/api/v1/user/", userAPI.LoginUser).Methods("POST")
	r.HandleFunc("/api/v1/user/{id}", userAPI.RegisterUser).Methods("POST")

	// Car routes
	r.HandleFunc("/api/v1/car/", carAPI.Cars).Methods("GET")
	r.HandleFunc("/api/v1/car/{id}", carAPI.CarDetails).Methods("GET")

	// Reservation routs
	r.HandleFunc("/api/v1/reservation/user/{id}", resAPI.UserReservations).Methods("GET")
	r.HandleFunc("/api/v1/reservation/car/{id}", resAPI.CarReservations).Methods("GET")

	// Start server
	log.Println("Server starting on :8081")
	log.Fatal(http.ListenAndServe(":8081", r))
}

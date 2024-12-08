package main

import (
	"log"
	"net/http"

	"CNADASG1/apis"
	"CNADASG1/config"
	"CNADASG1/services"

	"github.com/gorilla/handlers" // Import the handlers package
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
	payService := &services.PaymentService{DB: db}

	// Initialize apis
	userAPI := &apis.UserAPI{Service: userService}
	carAPI := &apis.CarAPI{Service: carService}
	resAPI := &apis.ReserveAPI{Service: resService}
	payAPI := &apis.PaymentAPI{Service: payService}

	// Create router
	r := mux.NewRouter()

	// User routes
	r.HandleFunc("/api/v1/user/{id}", userAPI.UserDetails).Methods("GET")
	r.HandleFunc("/api/v1/user/", userAPI.LoginUser).Methods("POST")
	r.HandleFunc("/api/v1/user/{id}", userAPI.RegisterUser).Methods("POST")

	// Car routes
	r.HandleFunc("/api/v1/car/", carAPI.Cars).Methods("GET")
	r.HandleFunc("/api/v1/car/{id}", carAPI.CarDetails).Methods("GET")
	r.HandleFunc("/api/v1/car/{id}", carAPI.UpdateCarStatus).Methods("PUT")

	// Reservation routes
	r.HandleFunc("/api/v1/reservation/user/all/{id}", resAPI.AllReservations).Methods("GET")
	r.HandleFunc("/api/v1/reservation/user/{id}", resAPI.UserReservations).Methods("GET")
	r.HandleFunc("/api/v1/reservation/car/{id}", resAPI.CarReservations).Methods("GET")
	r.HandleFunc("/api/v1/reservation/available-times", resAPI.GetAvailableTimes).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/v1/reservation/", resAPI.CreateReservation).Methods("POST")
	r.HandleFunc("/api/v1/reservation/update/{id}", resAPI.UpdateStatus).Methods("PUT")
	r.HandleFunc("/api/v1/reservation/{id}", resAPI.ReservationDetails).Methods("GET")
	r.HandleFunc("/api/v1/reservation/completed/{id}", resAPI.CompletedReservations).Methods("GET")
	r.HandleFunc("/api/v1/reservation/end/{id}", resAPI.EndReservation).Methods("PUT")

	// Payment routes
	r.HandleFunc("/api/v1/payment/res/", payAPI.CreatePayment).Methods("POST")
	r.HandleFunc("/api/v1/payment/{id}", payAPI.PaymentDetails).Methods("GET")
	r.HandleFunc("/api/v1/payment/{id}", payAPI.PaymentDetails).Methods("PUT")


	// Apply CORS middleware for JS
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),                             // Allow all origins (use specific ones for security)
		handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"}),        // Allowed HTTP methods
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}), // Allowed headers
	)
	// Wrap the router with CORS handler
	http.Handle("/", corsHandler(r))

	// Start server
	log.Println("Server starting on :8081")
	log.Fatal(http.ListenAndServe(":8081", r))
}

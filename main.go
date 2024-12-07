package main

import (
	"log"
	"net/http"

	"CNADASG1/config"
	"CNADASG1/handlers"
	"CNADASG1/templates"

	"github.com/gorilla/mux"
)

func main() {

	// base url
	baseURL := "http://localhost:8081/api/v1/"

	// load templates
	templates.InitializeTemplates()

	// Connect to database
	db := config.ConnectDatabase()
	defer db.Close()

	// Initialize handlers
	homeHandler := &handlers.HomeHandler{}
	userHandler := &handlers.UserHandler{BaseURL: baseURL}
	carHandler := &handlers.CarHandler{BaseURL: baseURL}
	resHandler := &handlers.ReserveHandler{BaseURL: baseURL}
	paymentHandler := &handlers.PaymentHandler{BaseURL: baseURL}

	// Create router
	r := mux.NewRouter()

	// Wrap routes with the NotFound middleware
	r.Use(handlers.NotFoundMiddleware)

	// Serve static files (CSS, JS, images)
	staticHandler := http.FileServer(http.Dir("./static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", staticHandler))

	// home routes
	r.HandleFunc("/", homeHandler.Home).Methods("GET")
	r.HandleFunc("/login", homeHandler.Login).Methods("GET")
	r.HandleFunc("/register", homeHandler.Register).Methods("GET")

	// User routes
	r.HandleFunc("/profile", handlers.AuthMiddleware(userHandler.UserDetails)).Methods("GET")
	r.HandleFunc("/login", userHandler.LoginUser).Methods("POST")
	r.HandleFunc("/register", userHandler.RegisterUser).Methods("POST")
	//r.HandleFunc("/login", userHandler.GetUser).Methods("GET")
	r.HandleFunc("/logout", userHandler.LogOutUser).Methods("GET")

	// Car routes
	r.HandleFunc("/cars", carHandler.HandleCars).Methods("GET")

	// Reservation routes
	r.HandleFunc("/reserve", resHandler.CarReservations).Methods("GET")
	r.HandleFunc("/reserve/car/{id}", resHandler.PostReservation).Methods("POST")
	r.HandleFunc("/reserve/user", resHandler.UserReservations).Methods("GET")
	r.HandleFunc("/cancel/{id}", resHandler.CancelReservation).Methods("POST")
	r.HandleFunc("/reserve/now/{id}", resHandler.ReserveNow).Methods("POST")

	// Payment routes
	r.HandleFunc("/payment", paymentHandler.Payments).Methods("POST")


	// Start server
	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

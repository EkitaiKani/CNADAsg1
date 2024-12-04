package main

import (
	"log"
	"net/http"

	"CNADASG1/config"
	"CNADASG1/handlers"
	"CNADASG1/services"
	"CNADASG1/templates"

	"github.com/gorilla/mux"
)

func main() {

	// load templates
	templates.InitializeTemplates()

	// Connect to database
	db := config.ConnectDatabase()
	defer db.Close()

	// Initialize services
	userService := &services.UserService{DB: db}
	carService := &services.CarService{DB: db}

	// Initialize handlers
	homeHandler := &handlers.HomeHandler{}
	userHandler := &handlers.UserHandler{Service: userService}
	carHandler := &handlers.CarHandler{Service: carService}

	// Create router
	r := mux.NewRouter()

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
	r.HandleFunc("/cars", carHandler.Cars).Methods("GET")
	r.HandleFunc("/cars/{id}", carHandler.CarDetails).Methods("GET")




	// Start server
	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

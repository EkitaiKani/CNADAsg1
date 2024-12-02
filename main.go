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

	// Initialize handlers
	homeHandler := &handlers.HomeHandler{}
	userHandler := &handlers.UserHandler{Service: userService}

	// Create router
	r := mux.NewRouter()

	// home routes
	r.HandleFunc("/", homeHandler.Home).Methods("GET")
	r.HandleFunc("/login", homeHandler.Login).Methods("GET")

	// User routes
	r.HandleFunc("/login", userHandler.CreateUser).Methods("POST")
	//r.HandleFunc("/login", userHandler.GetUser).Methods("GET")

	// Start server
	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

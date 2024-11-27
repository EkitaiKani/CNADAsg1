package main

import (
    "fmt"
    "html/template"
    "log"
    "net/http"

    "github.com/gorilla/mux"
)

var baseurl string
var tmpl *template.Template

func home(w http.ResponseWriter, r *http.Request) {

    err := tmpl.ExecuteTemplate(w, "home.html", nil)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}

func main() {
    baseurl = "http://localhost:8080/api/v1/"
    
    // Parse templates
    tmpl = template.Must(template.ParseGlob("templates/*.html"))
    
    // Create a new router
    r := mux.NewRouter()
    
    // Serve static files (CSS, JS, images)
    staticHandler := http.FileServer(http.Dir("./static"))
    r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", staticHandler))
    
    // Home route
    r.HandleFunc("/", home).Methods("GET")
    
    // Start the server with error handling
    fmt.Println("Server is running on http://localhost:5000")
    srv := &http.Server{
        Handler: r,
        Addr:    ":5000",
    }
    
    log.Fatal(srv.ListenAndServe())
}
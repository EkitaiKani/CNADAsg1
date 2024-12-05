package handlers

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"CNADASG1/models"
	"CNADASG1/templates"
	"CNADASG1/utils"

	"github.com/gorilla/sessions"
)

var (
	// Create a key for encrypting the session cookie
	// IMPORTANT: Use a secure, random key in production
	sessionKey = []byte("your-secret-key-here-make-it-long-and-random")

	// Create a session store
	store = sessions.NewCookieStore(sessionKey)
)

type login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserHandler struct {
	BaseURL string
}

// NewUserHandler is a constructor function to create a new UserHandler with the API base URL
func NewUserHandler(baseURL string) *UserHandler {
	return &UserHandler{
		BaseURL: baseURL + "/users/", // Store the base URL
	}
}

func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	// If the method is POST, handle form submission
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Retrieve form values by their "name" attribute
	email := r.FormValue("Email")
	username := r.FormValue("Username")
	password := r.FormValue("Password")
	first := r.FormValue("FirstName")
	last := r.FormValue("LastName")
	dobStr := r.FormValue("DateofBirth")
	log.Print(username)

	// Try to parse the date string to time.Time
	var dob sql.NullTime
	if dobStr != "" {
		parsedTime, err := time.Parse("2006-01-02", dobStr) // Expecting date in "YYYY-MM-DD" format
		if err != nil {
			log.Println("Error parsing Date of Birth:", err)
			http.Error(w, "Invalid date format", http.StatusBadRequest)
			return
		}
		dob = sql.NullTime{Time: parsedTime, Valid: true} // Mark as valid with parsed time
	} else {
		dob = sql.NullTime{Valid: false} // If empty, mark as invalid (representing NULL in SQL)
	}

	user := &models.User{
		UserEmail:    email,
		UserName:     username,
		HashPassword: password,
		FirstName:    first,
		LastName:     last,
		DateofBirth:  dob,
	}

	
}

func (h *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// If the method is POST, handle form submission
	// Parse the form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid username")
		return
	}

	// Retrieve form values by their "name" attribute
	var l login
	l.Username = r.FormValue("Username")
	l.Password = r.FormValue("Password")

	var response map[string]interface{}
	url := h.BaseURL
	client := &http.Client{}
	postBody, _ := json.Marshal(l)
	resBody := bytes.NewBuffer(postBody)

	if req, err := http.NewRequest("POST", url, resBody); err == nil {
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

	if response["error"] == "true" {
		renderErr := templates.Templates.ExecuteTemplate(w, "login.html", response)
		if renderErr != nil {
			http.Error(w, renderErr.Error(), http.StatusInternalServerError)
		}
		return
	}

	// log.Print(response)

	// Create a new session
	session, err := store.Get(r, "user-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Store the user ID in the session
	session.Values["user_id"] = response["user"]
	// Retrieve nested "zipcode" map
	if user, ok := response["user"].(map[string]interface{}); ok {
		// Retrieve "code" from the "zipcode" map
		if id, ok := user["id"].(int); ok {
			session.Values["user_id"] = id
		}
	}

	// Save the session
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// After successful login
	// log.Print(user.UserName)
	http.Redirect(w, r, "/", http.StatusSeeOther)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *UserHandler) LogOutUser(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "user-session")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Clear the session
	session.Values = make(map[interface{}]interface{})
	session.Options.MaxAge = -1 // Expire the session immediately

	// Save the session before redirecting
	if err := session.Save(r, w); err != nil {
		http.Error(w, "Session save error", http.StatusInternalServerError)
		return
	}

	// Redirect to login or home page
	log.Print("logged out successfully")
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// func (h *UserHandler) UserDetails(w http.ResponseWriter, r *http.Request) {
// 	session, err := store.Get(r, "user-session")
// 	if err != nil {
// 		http.Redirect(w, r, "/login", http.StatusSeeOther)
// 		return
// 	}

// 	// Convert userid from session state to int
// 	id, ok := session.Values["user_id"].(int)
// 	// log.Print(ok)
// 	if !ok {
// 		// Clear the session
// 		session.Values = make(map[interface{}]interface{})
// 		session.Options.MaxAge = -1 // Expire the session immediately

// 		// Save the session before redirecting
// 		if err := session.Save(r, w); err != nil {
// 			http.Error(w, "Session save error", http.StatusInternalServerError)
// 			return
// 		}

// 		// Redirect to login or home page
// 		http.Redirect(w, r, "/login", http.StatusSeeOther)
// 		return
// 	}

// 	// log.Print(id)

// 	// get user details
// 	user, err := h.Service.GetUserDetails(id)
// 	// log.Print(user.UserEmail)

// 	if err != nil {
// 		// Log the actual error
// 		renderErr := templates.Templates.ExecuteTemplate(w, "profile.html", map[string]interface{}{
// 			"message": "Error getting user details, please try again",
// 			"error":   true,
// 		})
// 		if renderErr != nil {
// 			http.Error(w, "Template render error", http.StatusInternalServerError)
// 		}
// 		return
// 	}

// 	// Render user details
// 	if err := templates.Templates.ExecuteTemplate(w, "profile.html", map[string]interface{}{
// 		"user":  user,
// 		"error": false,
// 	}); err != nil {
// 		http.Error(w, "Template render error", http.StatusInternalServerError)
// 	}
// }

// Middleware to check if user is logged in
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the session
		session, err := store.Get(r, "user-session")
		if err != nil {
			// If there's an error getting the session, redirect to login
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Check if user_id exists in session
		if session.Values["user_id"] == nil {
			// No user ID in session, redirect to login
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Convert userid from session state to int
		id, ok := session.Values["user_id"].(int)
		if !ok {
			// If user_id is not an int, redirect to login
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Attach the user ID to the request context
		ctx := context.WithValue(r.Context(), "user_id", id)

		// Call the next handler with the modified request
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

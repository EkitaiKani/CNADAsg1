package handlers

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"

	"CNADASG1/models"
	"CNADASG1/services"
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

type UserHandler struct {
	Service *services.UserService
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

	// Try to create the user using the service
	createdUser, err := h.Service.CreateUser(user)
	if err != nil {
		// Use a single template execution with an error message
		renderErr := templates.Templates.ExecuteTemplate(w, "register.html", map[string]interface{}{
			"message": "Account was not created. Check your input fields and try again.",
			"error":   true,
		})
		if renderErr != nil {
			http.Error(w, renderErr.Error(), http.StatusInternalServerError)
		}
		return
	}

	// After successful creation, render the user in the template
	err = templates.Templates.ExecuteTemplate(w, "register.html", map[string]interface{}{
		"user":    createdUser,
		"message": "Account created successfully",
		"error":   false,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
	name := r.FormValue("Username")
	password := r.FormValue("Password")

	user, err := h.Service.LogInUser(name, password)
	if err != nil {
		// Use a single template execution with an error message
		renderErr := templates.Templates.ExecuteTemplate(w, "login.html", map[string]interface{}{
			"message": "Username or password is incorrect.",
			"error":   true,
		})
		if renderErr != nil {
			http.Error(w, renderErr.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Create a new session
	session, err := store.Get(r, "user-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Store the user ID in the session
	session.Values["user_id"] = user.UserId

	// Set session options
	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, // 7 days
		HttpOnly: true,
		Secure:   true, // Use only over HTTPS
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

// Middleware to check if user is logged in
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := store.Get(r, "user-session")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Check if user is logged in
		UserId, ok := session.Values["user_id"].(int) // Adjust type as needed
		if !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// You can optionally attach the user ID to the request context
		ctx := context.WithValue(r.Context(), "user_id", UserId)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}



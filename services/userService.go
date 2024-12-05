package services

import (
	"database/sql"
	"log"
	"time"

	"CNADASG1/models"

	"golang.org/x/crypto/bcrypt" // to encrypt/hash passwords
)

type UserService struct {
	DB *sql.DB
}

// register new user
func (s *UserService) CreateUser(user *models.User) (*models.User, error) {
	// encrypt password
	hash, err := bcrypt.GenerateFromPassword([]byte(user.HashPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error hashing password:", err)
		return nil, err
	}
	user.HashPassword = string(hash)

	// Prepare the SQL INSERT statement
	query := "INSERT INTO users (email, password_hash, firstname, lastname, username, date_of_birth) VALUES (?, ?, ?, ?, ?, ?)"
	result, err := s.DB.Exec(query, user.UserEmail, user.HashPassword, user.FirstName, user.LastName, user.UserName, user.DateofBirth)
	if err != nil {
		log.Println("Database insert error:", err)
		return nil, err
	}

	// Get the last inserted ID
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		log.Println("Error getting last insert ID:", err)
		return nil, err
	}

	// Set the UserId in the user model
	user.UserId = int(lastInsertID)

	return user, nil
}

// log in user
func (s *UserService) LogInUser(username string, password string) (*models.User, error) {

	// define temp password and id
	var temppw string
	var id int

	// get user details
	query := "SELECT user_id, password_hash FROM users WHERE username = ?"
	err := s.DB.QueryRow(query, username).Scan(&id, &temppw)
	if err != nil {
		log.Println("Row scan error:", err)
		return nil, err
	}

	// check password
	err = bcrypt.CompareHashAndPassword([]byte(temppw), []byte(password))
	if err != nil {
		// log.Fatal("Password does not match:", err)

	} else {
		// if password is correct, store user
		user := &models.User{
			UserName: username,
			UserId:   id,
		}

		return user, nil
	}

	// return nil if passwords do not match
	return nil, err

}

// fetches user details for the user details page.
func (s *UserService) GetUserDetails(id int) (*models.User, error) {
    // Create a new instance of models.User
    u := &models.User{}

    // Set id
    u.UserId = id

    // Declare a variable for the date_of_birth as a sql.NullString
    var dob sql.NullString

    // Get user details
    query := "SELECT email, username, membership_tier, firstname, lastname, date_of_birth, is_verified FROM users WHERE user_id = ?"
    err := s.DB.QueryRow(query, id).Scan(&u.UserEmail, &u.UserName, &u.MemberTier, &u.FirstName, &u.LastName, &dob, &u.Verified)
    if err != nil {
        log.Println("Row scan error:", err)
        return nil, err
    }

    // Parse the date into a time.Time if it's not NULL
    if dob.Valid {
        parsedTime, err := time.Parse("2006-01-02", dob.String)
        if err != nil {
            log.Fatal(err)
        }
        u.DateofBirth = sql.NullTime{Time: parsedTime, Valid: true}
    } else {
        u.DateofBirth = sql.NullTime{Valid: false} // Handle NULL case
    }

    return u, nil
}

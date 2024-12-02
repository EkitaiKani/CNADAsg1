package services

import (
	"database/sql"
	"log"

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

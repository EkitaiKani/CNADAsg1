package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
)

// define global variables
var db *sql.DB
var cfg mysql.Config

func ConnectDatabase() *sql.DB {
	// Get database connection details from environment variables
	dbUser, dbPassword, dbHost, dbName := os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME")

	// Check if all necessary environment variables are set
	if dbUser == "" || dbPassword == "" || dbHost == "" || dbName == "" {
		log.Fatal("Database credentials not fully set in environment variables")
	}

	// Capture connection properties.
	cfg = mysql.Config{
		User:   dbUser,
		Passwd: dbPassword,
		Net:    "tcp",
		Addr:   dbHost,
		DBName: dbName,
	}

	//fmt.Print(cfg)

	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	// Get a database handle.
	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")
	return db
}

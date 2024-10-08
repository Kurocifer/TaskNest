package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var DB *sql.DB

// SetUpDatabase sets up and Establishes the connection to the database
func SetUpDatabase() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var dbUser = os.Getenv("DB_USER")
	var dbPassword = os.Getenv("DB_PASSWORD")
	var dbName = os.Getenv("DB_NAME")

	fmt.Println("Establishing database connection...")

	dsn := dbUser + ":" + dbPassword + "@tcp(localhost:3306)/" + dbName

	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}
	if err = DB.Ping(); err != nil {
		log.Fatal("Database is unreachable:", err)
	}

	fmt.Println("Database connection successful!")
}

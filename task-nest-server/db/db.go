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

// promptForUsername prompts the user for their mysql database username
//
// parameters:
//
//	None
//
// Returns:
//
//	the username
/*
func promptForUsername() string {
	var username string
	prompt := &survey.Input{
		Message: "Enter your MySQL username: ",
	}

	survey.AskOne(prompt, &username)
	return username
}

func promptForPassword() string {
	fmt.Print("Enter you MySQL password: ")
	bytePassword, err := terminal.ReadPassword(0)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println() // Add new line after password
	return string(bytePassword)
}

func promptForDatabaseName() string {
	var databaseName string
	prompt := &survey.Input{
		Message: "Enter the MySQL database name: ",
	}

	survey.AskOne(prompt, &databaseName)
	return databaseName
}
*/

func SetupDatabase() {
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

	fmt.Println("Database connection successful")
}

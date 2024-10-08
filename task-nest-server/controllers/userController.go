package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/kurocifer/TaskNest/task-nest-server/db"
	"github.com/kurocifer/TaskNest/task-nest-server/models"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("wenderlichPrime_")

// hashPassword hashes the password string using bcrypt.
// Returns the hashed password.
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// verifyPassword compares a bcrypt hashed password with it's passible plain text.
// Returns nil on success, or error on failure.
func verifyPassword(hashedPassword, plainPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
}

// RegisterUser is the handler function for Registering a user in the db database.
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	//check if it's a POST reqeust
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var user models.UserAuthRequestBody
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid input."})
		return
	}

	// Trim whitespace and check for empty values
	checkUsername := strings.TrimSpace(user.Username)
	checkPassword := strings.TrimSpace(user.Password)

	if checkPassword == "" || checkUsername == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Username or password cannot be empty or contain only white space"})
		return
	}

	//check if user already exists
	query := "SELECT id FROM users WHERE username = ?"
	err = db.DB.QueryRow(query, user.Username).Scan(new(int))

	if err == nil {
		// user already exists
		w.WriteHeader(http.StatusConflict) // 409 conflict
		json.NewEncoder(w).Encode(map[string]string{"error": "User aleady exists"})
		return
	} else if err != sql.ErrNoRows {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// generate this token before saving user instead of generating after saving user
	// because a user might be created successfully but some error prevented token creation
	// which can lead to an unknown pending user in the database
	token, err := genereteJWT(user.Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Printf("generated token %s", token)

	//save user to database
	hashedPassword, err := hashPassword(user.Password)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	query = "INSERT INTO users (username, password) VALUES (?, ?)"
	result, err := db.DB.Exec(query, user.Username, hashedPassword)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userID, err := result.LastInsertId()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	createdUser := models.UserAuthResponseBody{
		ID:       int(userID),
		Username: user.Username,
		Token:    token,
	}

	createdUser.Username = user.Username
	createdUser.ID = int(userID)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdUser)
}

// LoginUser is the handler function that is called to log in a user.
// Returns a generated JWT for the user's session.
func LoginUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var user models.UserAuthRequestBody
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid input."})
		return
	}

	// Trim whitespace and check for empty values
	user.Username = strings.TrimSpace(user.Username)
	user.Password = strings.TrimSpace(user.Password)

	if user.Username == "" || user.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Username or password cannot be empty or contain only white space"})
		return
	}

	// compare passwords
	var retrievedUser models.User
	query := "SELECT id, password FROM users WHERE username = ?"
	err = db.DB.QueryRow(query, user.Username).Scan(&retrievedUser.ID, &retrievedUser.Password)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if err := verifyPassword(retrievedUser.Password, user.Password); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Wrong username or password."})
		return
	}

	token, err := genereteJWT(user.Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userResponseDetails := models.UserAuthResponseBody{
		ID:       retrievedUser.ID,
		Username: user.Username,
		Token:    token,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(userResponseDetails)
}

// generateJWT generates the JWT token that expires in 1 hour for a user's session,
// using username for building the claim and signs it using a secrete key.
// Returns the signed generated JWT on Success, or error on failure.
func genereteJWT(username string) (string, error) {
	expirationTime := time.Now().Add(1 * time.Hour) // expires in 1 hour

	claims := &models.Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

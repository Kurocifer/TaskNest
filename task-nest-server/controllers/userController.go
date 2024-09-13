package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/kurocifer/TaskNest/task-nest-server/db"
	"github.com/kurocifer/TaskNest/task-nest-server/models"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("wenderlichPrime_")

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func verifyPassword(hashedPassword, plainPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	//check if it's a POST reqeust
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var user models.UserAuthRequestBody
	json.NewDecoder(r.Body).Decode(&user)
	fmt.Println(user)

	//check if user already exists
	query := "SELECT id FROM users WHERE username = ?"
	err := db.DB.QueryRow(query, user.Username).Scan(new(int))

	if err == nil {
		// user already exists
		w.WriteHeader(http.StatusConflict) // 409 conflict
		json.NewEncoder(w).Encode(map[string]string{"error": "User aleady exists"})
		return
	} else if err != sql.ErrNoRows {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

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

	var createdUser models.UserRegistrationResponse

	createdUser.Username = user.Username
	createdUser.ID = int(userID)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdUser)
	fmt.Println(w)
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var user models.UserAuthRequestBody
	json.NewDecoder(r.Body).Decode(&user)

	// compare passwords
	var retrievedUser models.User
	query := "SELECT id, password FROM users WHERE username = ?"
	err := db.DB.QueryRow(query, user.Username).Scan(&retrievedUser.ID, &retrievedUser.Password)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if err := verifyPassword(retrievedUser.Password, user.Password); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	token, err := genereteJWT(user.Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

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

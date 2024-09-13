package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/kurocifer/TaskNest/task-nest-server/db"
	"github.com/kurocifer/TaskNest/task-nest-server/models"
)

func CreateTodo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	tokenBearer := r.Header.Get("Authorization")
	tokenString := strings.TrimPrefix(tokenBearer, "Bearer ")
	claims := &models.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var todo models.Todo
	json.NewDecoder(r.Body).Decode(&todo)
	todo.Username = claims.Username

	query := "INSERT INTO todos (user_id, task, done) VALUES (?, ?, ?)"
	result, err := db.DB.Exec(query, todo.Username, todo.Task, todo.Done)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	todoID, err := result.LastInsertId()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	todo.ID = int(todoID)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(todo)
}

func GetTodos(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	tokenBearer := r.Header.Get("Authorization")
	tokenString := strings.TrimPrefix(tokenBearer, "Bearer ")
	claims := &models.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	query := "SELECT id, task, done FROM todos WHERE user_id = ?"
	rows, err := db.DB.Query(query, claims.Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var userTodos models.UserTodos
	//var todos []models.Todo

	userTodos.UserId = claims.Username
	for rows.Next() {
		var todo models.Todo

		if err := rows.Scan(&todo.ID, &todo.Task, &todo.Done); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		userTodos.Todos = append(userTodos.Todos, todo)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(userTodos)

}

func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	tokenBearer := r.Header.Get("Authorization")
	tokenString := strings.TrimPrefix(tokenBearer, "Bearer ")
	claims := &models.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Extract the TODO ID from the URL
	todoID := r.URL.Query().Get("id")
	if todoID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	query := "DELETE FROM todos WHERE id = ? AND user_id = ?"
	result, err := db.DB.Exec(query, todoID, claims.Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound) // Todo not found or not owned by this user
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func UpdateTodoStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	//get token
	tokenBearer := r.Header.Get("Authorization")
	tokenString := strings.TrimPrefix(tokenBearer, "Bearer ")
	claims := &models.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Extract the Todo ID and new status from the URL
	todoId := r.URL.Query().Get("id")
	todoStatus := r.URL.Query().Get("done")
	if todoId == "" || todoStatus == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// convert todoStatus to boolean
	done := todoStatus == "true"
	query := "UPDATE todos SET done = ? WHERE id = ? AND user_id = ?"
	_, err = db.DB.Exec(query, done, todoId, claims.Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Todo status changed"})
}

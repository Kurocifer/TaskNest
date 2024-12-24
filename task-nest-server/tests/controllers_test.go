package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kurocifer/TaskNest/task-nest-server/controllers"
	"github.com/kurocifer/TaskNest/task-nest-server/db"
	"github.com/kurocifer/TaskNest/task-nest-server/models"
)

func TestRegisterUser(t *testing.T) {
	db.SetUpDatabase()
	defer db.DB.Close()

	user := models.User{Username: "tesser", Password: "password123"}
	userJson, _ := json.Marshal(user)

	req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(userJson))
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controllers.RegisterUser)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}
}

func TestRegisterUserAlreadyExists(t *testing.T) {
	db.SetUpDatabase()
	defer db.DB.Close()

	// First register the user
	user := models.UserAuthRequestBody{
		Username: "existinguser",
		Password: "password123",
	}
	userJson, _ := json.Marshal(user)

	req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(userJson))
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controllers.RegisterUser)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	// Attempt to register the same user again
	req, _ = http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(userJson))
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusConflict {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusConflict)
	}

	var errorResponse map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&errorResponse); err != nil {
		t.Fatalf("Could not decode error response: %v", err)
	}

	if errorResponse["error"] != "User aleady exists" {
		t.Errorf("Expected error message 'User aleady exists', got '%s'", errorResponse["error"])
	}
}

func TestRegisterUserInvalidInput(t *testing.T) {
	db.SetUpDatabase()
	defer db.DB.Close()

	// Attempt to register with invalid input (empty username and password)
	user := models.UserAuthRequestBody{
		Username: "",
		Password: "",
	}
	userJson, _ := json.Marshal(user)

	req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(userJson))
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controllers.RegisterUser)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	var errorResponse map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&errorResponse); err != nil {
		t.Fatalf("Could not decode error response: %v", err)
	}

	if errorResponse["error"] != "Username or password cannot be empty or contain only white space" {
		t.Errorf("Expected error message 'Username or password cannot be empty or contain only white space', got '%s'", errorResponse["error"])
	}
}

func TestLoginUser(t *testing.T) {
	db.SetUpDatabase()
	defer db.DB.Close()

	user := models.User{Username: "tesser", Password: "password123"}
	userJson, _ := json.Marshal(user)

	// Now try to log in
	req, err := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(userJson))
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controllers.LoginUser)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestLoginUserInvalidInput(t *testing.T) {
	db.SetUpDatabase()
	defer db.DB.Close()

	// Attempt to log in with invalid input (empty username and password)
	user := models.UserAuthRequestBody{
		Username: "",
		Password: "",
	}
	userJson, _ := json.Marshal(user)

	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(userJson))
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controllers.LoginUser)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	var errorResponse map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&errorResponse); err != nil {
		t.Fatalf("Could not decode error response: %v", err)
	}

	if errorResponse["error"] != "Username or password cannot be empty or contain only white space" {
		t.Errorf("Expected error message 'Username or password cannot be empty or contain only white space', got '%s'", errorResponse["error"])
	}
}

func TestLoginUserUnauthorized(t *testing.T) {
	db.SetUpDatabase()
	defer db.DB.Close()

	// Attempt to log in with incorrect password
	user := models.UserAuthRequestBody{
		Username: "tesser",
		Password: "wrongpassword",
	}
	userJson, _ := json.Marshal(user)

	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(userJson))
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controllers.LoginUser)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
	}

	var errorResponse map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&errorResponse); err != nil {
		t.Fatalf("Could not decode error response: %v", err)
	}

	if errorResponse["error"] != "Wrong username or password." {
		t.Errorf("Expected error message 'Wrong username or password.', got '%s'", errorResponse["error"])
	}
}

func TestCreateTodo(t *testing.T) {
	db.SetUpDatabase()
	defer db.DB.Close()

	user := models.User{Username: "tesser", Password: "password123"}
	userJson, _ := json.Marshal(user)

	// Log in to get the token
	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(userJson))
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controllers.LoginUser)
	handler.ServeHTTP(rr, req)

	var response map[string]string
	json.NewDecoder(rr.Body).Decode(&response)
	token := response["token"]

	// Now create a todo
	todo := models.Todo{Task: "Test Todo", Done: false}
	todoJson, _ := json.Marshal(todo)

	req, _ = http.NewRequest(http.MethodPost, "/todos/create", bytes.NewBuffer(todoJson))
	req.Header.Set("Authorization", token)
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(controllers.CreateTodo)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}
}

func TestCreateTodoUnauthorized(t *testing.T) {
	db.SetUpDatabase()
	defer db.DB.Close()

	todo := models.Todo{Task: "Unauthorized Todo", Done: false}
	todoJson, _ := json.Marshal(todo)

	req, _ := http.NewRequest(http.MethodPost, "/todos", bytes.NewBuffer(todoJson))
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controllers.CreateTodo)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
	}
}

func TestGetTodos(t *testing.T) {
	db.SetUpDatabase()
	defer db.DB.Close()

	// First register and log in the user
	user := models.User{Username: "tesser", Password: "password123"}
	userJson, _ := json.Marshal(user)

	// Log in to get the token
	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(userJson))
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controllers.LoginUser)
	handler.ServeHTTP(rr, req)

	var response map[string]string
	json.NewDecoder(rr.Body).Decode(&response)
	token := response["token"]

	// Now get todos
	req, _ = http.NewRequest(http.MethodGet, "/todos/list", nil)
	req.Header.Set("Authorization", token)
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(controllers.GetTodos)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestGetTodosUnauthorized(t *testing.T) {
	db.SetUpDatabase()
	defer db.DB.Close()

	req, _ := http.NewRequest(http.MethodGet, "/todos/list", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controllers.GetTodos)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
	}
}

func TestDeleteTodo(t *testing.T) {
	db.SetUpDatabase()
	defer db.DB.Close()

	user := models.User{Username: "tesser", Password: "password123"}
	userJson, _ := json.Marshal(user)

	// Log in to get the token
	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(userJson))
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controllers.LoginUser)
	handler.ServeHTTP(rr, req)

	var loginResponse map[string]string
	json.NewDecoder(rr.Body).Decode(&loginResponse)
	token := loginResponse["token"]

	// first create the todo to be sure it's present and get it's ID
	todo := models.Todo{Task: "Test Todo deletion", Done: false}
	todoJson, _ := json.Marshal(todo)

	req, _ = http.NewRequest(http.MethodPost, "/todos/create", bytes.NewBuffer(todoJson))
	req.Header.Set("Authorization", token)
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(controllers.CreateTodo)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	var createTodoResponse map[string]int
	json.NewDecoder(rr.Body).Decode(&createTodoResponse)
	todoID := createTodoResponse["id"]

	// Now delete the todo
	req, _ = http.NewRequest(http.MethodDelete, fmt.Sprintf("/todos?id=%d", todoID), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(controllers.DeleteTodo)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
	}

}

func TestDeleteTodoUnauthorized(t *testing.T) {
	db.SetUpDatabase()
	defer db.DB.Close()

	req, _ := http.NewRequest(http.MethodDelete, "/todos?id=1", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controllers.DeleteTodo)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
	}
}

func TestDeleteTodoNotFound(t *testing.T) {
	db.SetUpDatabase()
	defer db.DB.Close()

	user := models.User{Username: "tesser", Password: "password123"}
	userJson, _ := json.Marshal(user)

	// Log in to get the token
	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(userJson))
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controllers.LoginUser)
	handler.ServeHTTP(rr, req)

	var loginResponse map[string]string
	json.NewDecoder(rr.Body).Decode(&loginResponse)
	token := loginResponse["token"]

	// Attempt to delete a non-existent todo
	req, _ = http.NewRequest(http.MethodDelete, "/todos?id=1024", nil) // Assuming todo with ID as 1024 does not exist
	req.Header.Set("Authorization", "Bearer "+token)
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(controllers.DeleteTodo)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}

func TestUpdateTodoStatus(t *testing.T) {
	db.SetUpDatabase()
	defer db.DB.Close()

	user := models.User{Username: "tesser", Password: "password123"}
	userJson, _ := json.Marshal(user)

	// Log in to get the token
	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(userJson))
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controllers.LoginUser)
	handler.ServeHTTP(rr, req)

	var loginResponse map[string]string
	json.NewDecoder(rr.Body).Decode(&loginResponse)
	token := loginResponse["token"]

	// first create the todo to be sure it's present and get it's ID
	todo := models.Todo{Task: "Test Todo deletion", Done: false}
	todoJson, _ := json.Marshal(todo)

	req, _ = http.NewRequest(http.MethodPost, "/todos/create", bytes.NewBuffer(todoJson))
	req.Header.Set("Authorization", token)
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(controllers.CreateTodo)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	var createTodoResponse map[string]int
	json.NewDecoder(rr.Body).Decode(&createTodoResponse)
	todoID := createTodoResponse["id"]

	// update the todo's status.
	newStatus := "false"
	req, _ = http.NewRequest(http.MethodDelete, fmt.Sprintf("/todos?id=%d&done%s", todoID, newStatus), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(controllers.DeleteTodo)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
	}
}

func TestUpdateTodoStatusUnauthorized(t *testing.T) {
	db.SetUpDatabase()
	defer db.DB.Close()

	req, _ := http.NewRequest(http.MethodPut, "/todos?id=1&done=true", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controllers.UpdateTodoStatus)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
	}
}

func TestUpdateTodoStatusBadRequest(t *testing.T) {
	db.SetUpDatabase()
	defer db.DB.Close()

	user := models.User{Username: "tesser", Password: "password123"}
	userJson, _ := json.Marshal(user)

	// Log in to get the token
	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(userJson))
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controllers.LoginUser)
	handler.ServeHTTP(rr, req)

	var loginResponse map[string]string
	json.NewDecoder(rr.Body).Decode(&loginResponse)
	token := loginResponse["token"]

	// Attempt to update with an invalid status
	req, _ = http.NewRequest(http.MethodPut, "/todos?id=1&done=invalid", nil) // Invalid status
	req.Header.Set("Authorization", "Bearer "+token)
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(controllers.UpdateTodoStatus)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestUpdateTodoStatusTodoNotFound(t *testing.T) {
	db.SetUpDatabase()
	defer db.DB.Close()

	user := models.User{Username: "tesser", Password: "password123"}
	userJson, _ := json.Marshal(user)

	// Log in to get the token
	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(userJson))
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controllers.LoginUser)
	handler.ServeHTTP(rr, req)

	var loginResponse map[string]string
	json.NewDecoder(rr.Body).Decode(&loginResponse)
	token := loginResponse["token"]

	// Attempt to update a non-existent todo
	req, _ = http.NewRequest(http.MethodPut, "/todos?id=999&done=false", nil) // Assuming 999 does not exist
	req.Header.Set("Authorization", "Bearer "+token)
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(controllers.UpdateTodoStatus)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

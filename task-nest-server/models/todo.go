package models

type Todo struct {
	ID       int    `json:"id"`
	Username string `json:"user_name"`
	Task     string `json:"task"`
	Done     bool   `json:"done"`
}

type FoundTodo struct {
	ID   int    `json:"id"`
	Task string `json:"task"`
	Done bool   `json:"done"`
}

type UserTodos struct {
	UserId string      `json:"user_name"`
	Todos  []FoundTodo `json:"todos"`
}

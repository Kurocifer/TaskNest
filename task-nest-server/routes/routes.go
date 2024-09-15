package routes

import (
	"net/http"

	"github.com/kurocifer/TaskNest/task-nest-server/controllers"
)

// SetupRoutes sets up the routes, linking the different routes to their respective handle functions.
func SetUpRoutes() {

	http.HandleFunc("/tasknest/register", controllers.RegisterUser)
	http.HandleFunc("/tasknest/login", controllers.LoginUser)
	http.HandleFunc("/tasknest/todos/create", controllers.CreateTodo)
	http.HandleFunc("/tasknest/todos/list", controllers.GetTodos)
	http.HandleFunc("/tasknest/todos/delete", controllers.DeleteTodo)
	http.HandleFunc("/tasknest/todos/update-status", controllers.UpdateTodoStatus)
}

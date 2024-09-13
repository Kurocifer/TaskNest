package routes

import (
	"net/http"

	"github.com/kurocifer/TaskNest/task-nest-server/controllers"
	"github.com/kurocifer/TaskNest/task-nest-server/db"
)

func SetUpRoutes() {
	db.SetupDatabase()

	http.HandleFunc("/register", controllers.RegisterUser)
	http.HandleFunc("/lgin", controllers.LoginUser)
}

package main

import (
	"log"
	"net/http"

	"github.com/kurocifer/TaskNest/task-nest-server/db"
)

func main() {
	db.SetupDatabase()

	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

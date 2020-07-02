package main

import (
	todo "github.com/urbaniakm/todo-list-api/todo"
	"net/http"
)

const apiBasePath = "/api"

func main() {
	todo.SetupRoutes(apiBasePath)
	http.ListenAndServe(":5000", nil)
}

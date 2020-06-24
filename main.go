package main

import (
	"encoding/json"
	"net/http"
	"log"
)

type fooHandler struct {
	Message string
}

type Todo struct {
	TodoId	int			`json:"todoId"`
	Name		string	`json:"name"`
}

var todoList []Todo
func init() {
	todosJSON := `[
		{
			"todoId": 1,
			"name": "Write simple To-Do app in Vue"
		},
		{
			"todoId": 2,
			"name": "Write simple To-Do app in SwiftUI for iOS"
		},
		{
			"todoId": 3,
			"name": "Write simple To-Do app in Angular"
		},
		{
			"todoId": 4,
			"name": "Write simple To-Do app in Kotlin for Android"
		}
	]`

	err := json.Unmarshal([]byte(todosJSON), &todoList)
	if err != nil {
		log.Fatal(err)
	}
}

func todosHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		productsJSON, err := json.Marshal(todoList)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(productsJSON)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func main() {
	http.HandleFunc("/todos", todosHandler)
	http.ListenAndServe(":5000", nil)
}

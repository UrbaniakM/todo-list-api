package main

import (
	"fmt"
	"encoding/json"
	"net/http"
	"log"
	"io/ioutil"
	"time"
)

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

func getNextId() int {
	max := 0
	for index, todo := range todoList {
    if index == 0 || todo.TodoId > max {
        max = todo.TodoId
    }
	}
	return max + 1
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
		return
	case http.MethodPost:
		// add new todo
		var newTodo Todo
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = json.Unmarshal(bodyBytes, &newTodo)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if newTodo.TodoId != 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		newTodo.TodoId = getNextId()
		todoList = append(todoList, newTodo)
		w.WriteHeader(http.StatusCreated)
		return 
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func middlewareHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		fmt.Println("before handler; middleware  start")
		start := time.Now()
		handler.ServeHTTP(w, r)
		fmt.Printf("middleware  finished, %s", time.Since(start))
	})
}

func main() {
	todostHandlerFunc := http.HandlerFunc(todosHandler)

	http.Handle("/todos", middlewareHandler(todostHandlerFunc));
	http.ListenAndServe(":5000", nil)
}

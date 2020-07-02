package todo

import (
	"fmt"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"strings"
	"strconv"
	"time"
	cors "github.com/urbaniakm/todo-list-api/cors"
)

const todosBasePath = "todos"
func SetupRoutes(apiBasePath string) {
	todosHandlerFunc := http.HandlerFunc(todosHandler)
	todoHandlerFunc := http.HandlerFunc(todoHandler)

	http.Handle(fmt.Sprintf("%s/%s", apiBasePath, todosBasePath), 
		middlewareHandler(cors.Middleware(todosHandlerFunc)))
	http.Handle(fmt.Sprintf("%s/%s/", apiBasePath, todosBasePath), 
		middlewareHandler(cors.Middleware(todoHandlerFunc)))
}

func middlewareHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		fmt.Println("before handler; middleware  start")
		start := time.Now()
		handler.ServeHTTP(w, r)
		fmt.Printf("middleware  finished, %s", time.Since(start))
	})
}

func todosHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// get all todos
		todoList := getTodoList()
		todosJSON, err := json.Marshal(todoList)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(todosJSON)
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

		_, err = addOrUpdateTodo(newTodo)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusCreated)
		return 
	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func todoHandler(w http.ResponseWriter, r *http.Request) {
	urlPathSegments := strings.Split(r.URL.Path, "todos/")

	todoId, err := strconv.Atoi(urlPathSegments[len(urlPathSegments) - 1])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return 
	}

	todo := getTodo(todoId)
	if todo == nil {
		w.WriteHeader(http.StatusNotFound)
		return 
	}

	switch r.Method {
	case http.MethodGet:
		// get single todo
		todoJSON, err := json.Marshal(todo)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(todoJSON)
		return
	case http.MethodPut:
		// update todo in the list
		var updatedTodo Todo
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = json.Unmarshal(bodyBytes, &updatedTodo)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if updatedTodo.TodoId != todoId {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		addOrUpdateTodo(updatedTodo)
		w.WriteHeader(http.StatusOK)
	case http.MethodDelete:
		removeTodo(todoId)
		w.WriteHeader(http.StatusOK)
		return
	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
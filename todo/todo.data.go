package todo

import (
	"sync"
	"log"
	"fmt"
	"os"
	"io/ioutil"
	"encoding/json"
	"sort"
)

var todoMap = struct {
	sync.RWMutex
	m map[int]Todo
}{m: make(map[int]Todo)}

func init() {
	fmt.Println("Loading todos")
	
	loadedMap, err := loadTodoMap()
	todoMap.m = loadedMap

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%d todos loaded\n", len(todoMap.m))
}

func loadTodoMap() (map[int]Todo, error) {
	fileName := "todos.json"
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("file [%s] does not exist", fileName)
	}

	todoList := make([]Todo, 0)
	file, _ := ioutil.ReadFile(fileName)
	err = json.Unmarshal(file, &todoList)
	if err != nil {
		log.Fatal(err)

		return nil, err
	}

	todoMap := make(map[int]Todo)
	for _, todo := range(todoList) {
		todoMap[todo.TodoId] = todo
	}

	return todoMap, nil
}

func getTodo(todoId int) *Todo {
	todoMap.RLock()
	defer todoMap.RUnlock()

	if todo, ok := todoMap.m[todoId]; ok {
		return &todo
	}
	return nil
}

func getTodoList() []Todo {
	todoMap.RLock()
	defer todoMap.RUnlock()

	todos := make([]Todo, 0, len(todoMap.m))
	for _, todo := range todoMap.m {
		todos = append(todos, todo)
	}

	return todos
}

func getTodoIds() []int {
	todoMap.RLock()
	defer todoMap.RUnlock()

	todoIds := []int{}
	for key := range todoMap.m {
		todoIds = append(todoIds, key)
	}
	sort.Ints(todoIds)

	return todoIds
}

func getNextTodoId() int {
	todoIds := getTodoIds()
	return todoIds[len(todoIds) - 1] + 1
}

func addOrUpdateTodo(todo Todo) (int, error) {
	// if the todo id is set then update, otherwise add
	addOrUpdateId := -1

	if todo.TodoId > 0 {
		oldTodo := getTodo(todo.TodoId)
		// replace todo if it exists, otherwise return error
		if oldTodo == nil {
			return 0, fmt.Errorf("todo id [%d] does not exist", todo.TodoId)
		}
		addOrUpdateId = todo.TodoId
	} else {
		addOrUpdateId = getNextTodoId()
		todo.TodoId = addOrUpdateId
	}

	todoMap.Lock()
	todoMap.m[addOrUpdateId] = todo
	todoMap.Unlock()

	return addOrUpdateId, nil
}
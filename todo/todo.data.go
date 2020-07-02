package todo

import (
	"sync"
	"log"
	"fmt"
	"os"
	"io/ioutil"
	"encoding/json"
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


package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/danilomarques1/todoexample/model"
)

var app App

func TestMain(m *testing.M) {
	app.Initialize(ConnectionDb{
		Host: "0.0.0.0",
		Port: "5433",
		User: "fitz",
		Password: "123456",
		Dbname: "todoexample",
	})
	code := m.Run()
	clearTables()
	os.Exit(code)
}

func clearTables() {
	stmt, err := app.Db.Prepare("truncate table todo")
	if err != nil {
		log.Fatalf("error preparing statement %v", err)
	}
	if _, err = stmt.Exec(); err != nil {
		log.Fatalf("error cleaning table %v", err)
	}
}

func executeRequest(request *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	app.mux.ServeHTTP(rr, request)

	return rr
}

func addTodo(size int) {
	todoName := "Todo "
	todoDesc := "Desc "
	for i:= 0; i < size; i++ {
		app.Db.Exec("insert into todo(title, description) values($1, $2)", todoName + strconv.Itoa(i), todoDesc + strconv.Itoa(i))
	}
}

func TestAddTodo(t *testing.T) {
	todoRequest := model.Todo{
		Title: "Todo 1",
		Desc: "Desc 1",
	}
	b, err := json.Marshal(todoRequest)
	req, err := http.NewRequest(http.MethodPost, "/todo",  bytes.NewBuffer(b))
	if err != nil {
		t.Errorf("Error creating the request %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	response := executeRequest(req)
	if (response.Code != http.StatusCreated) {
		t.Errorf("Wrong status code returned, expect 201 got %v", response.Code)
	}
	var todo model.Todo
	err = json.Unmarshal(response.Body.Bytes(), &todo)
	if err != nil {
		t.Errorf("Error parsing the response %v", err)
	}
	if todo.Title != "Todo 1" {
		t.Errorf("Wrong title returned. Expect Todo 1 got %v", todo.Title)
	}
	clearTables()
}

func TestGetTodos(t *testing.T) {
	size := 3
	addTodo(size)
	request, err := http.NewRequest(http.MethodGet, "/todo", nil)
	if err != nil {
		t.Errorf("Error creating the request %v", err)
	}
	response := executeRequest(request)
	if response.Code != http.StatusOK {
		t.Errorf("Wrong status code, expect 200 got %v", response.Code)
	}
	todos := make([]model.Todo, size)
	err = json.Unmarshal(response.Body.Bytes(), &todos)
	if err != nil {
		t.Errorf("Error parsing body %v", err)
	}
	if len(todos) != size {
		t.Errorf("Invalid return, expect size %v, got %v", size, len(todos))
	}

	clearTables()
}

func TestGetTodo(t *testing.T) {
	addTodo(3)
	var id int
	app.Db.QueryRow("select (id) from todo where title = 'Todo 1'").Scan(&id)

	request, err := http.NewRequest(http.MethodGet, "/todo/"+strconv.Itoa(id), nil)
	if err != nil {
		t.Errorf("Error creating request %v", err)
	}
	response := executeRequest(request)
	if response.Code != http.StatusOK {
		t.Errorf("[todos request] Wrong status code expect 200 got %v", response.Code)
	}
	var todo model.Todo
	json.Unmarshal(response.Body.Bytes(), &todo)
	if todo.Title != "Todo 1" {
		t.Errorf("Wrong todo returned expect title Todo 1 got %v", todo.Title)
	}

	clearTables()
}

package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/danilomarques1/todoexample/model"
	"github.com/danilomarques1/todoexample/util"
	"github.com/gorilla/mux"
)

type TodoHandler struct {
	todoRepository model.ITodo
}

func NewTodoHandler(todoRepository model.ITodo) *TodoHandler {
	return &TodoHandler{
		todoRepository: todoRepository,
	}
}

func (todoHandler *TodoHandler) AddTodo(w http.ResponseWriter, r *http.Request) {
	var todo model.Todo
	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		util.RespondError(w, "Invalid body", http.StatusBadRequest)
		return
	}

	err = todoHandler.todoRepository.AddTodo(&todo)
	if err != nil {
		util.RespondError(w, "Error adding the new todo", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(todo)
}

func (todoHandler *TodoHandler) GetTodos(w http.ResponseWriter, r *http.Request) {
	todos, err := todoHandler.todoRepository.GetTodos()
	if err != nil {
		util.RespondError(w, "Error getting todos", http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(todos)
}

func (todoHandler *TodoHandler) GetTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		util.RespondError(w, "Wrong parameters", 400)
		return
	}
	todo, err := todoHandler.todoRepository.GetTodo(id)
	if err != nil {
		util.RespondError(w, "Todo not found", 404)
		return
	}

	json.NewEncoder(w).Encode(todo)
}

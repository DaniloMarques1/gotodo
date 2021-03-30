package repository

import (
	"database/sql"
	"fmt"

	"github.com/danilomarques1/todoexample/model"
)

type TodoRepository struct {
	db *sql.DB
}

func NewTodoRepository(db *sql.DB) *TodoRepository {
	return &TodoRepository{
		db: db,
	}
}

// implementing ITodo
func (todoRepository *TodoRepository) AddTodo(todo *model.Todo) error {
	stmt, err := todoRepository.db.Prepare(`insert into todo(title, description) values($1, $2) RETURNING id`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	err = stmt.QueryRow(todo.Title, todo.Desc).Scan(&todo.Id)
	if err != nil {
		return err
	}

	return nil
}


func (todoRepository *TodoRepository) GetTodos() ([]model.Todo, error) {
	todos := make([]model.Todo, 0)
	stmt, err := todoRepository.db.Prepare("select * from todo")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var todo model.Todo
		rows.Scan(&todo.Id, &todo.Title, &todo.Desc)
		todos = append(todos, todo)
	}

	return todos, nil
}

func (todoRepository *TodoRepository) GetTodo(id int) (*model.Todo, error) {
	stmt, err := todoRepository.db.Prepare("select * from todo where id = $1")
	if err != nil {
		fmt.Printf("Error = %v\n", err)
		return nil, err
	}
	rows := stmt.QueryRow(id)
	var todo model.Todo
	err = rows.Scan(&todo.Id, &todo.Title, &todo.Desc)
	if err != nil {
		fmt.Printf("Error = %v\n", err)
		return nil, err
	}

	return &todo, nil
}

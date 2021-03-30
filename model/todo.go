package model

type Todo struct {
	Id    int     `json:"id,omitempty"`
	Title string  `json:"title"`
	Desc  string  `json:"desc"`
}

type ITodo interface {
	AddTodo(t *Todo) error
	GetTodos() ([]Todo, error)
	GetTodo(id int) (*Todo, error)
}

package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/danilomarques1/todoexample/handler"
	"github.com/danilomarques1/todoexample/repository"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type App struct {
	Db  *sql.DB
	mux *mux.Router
}

type ConnectionDb struct {
	Host 	 string
	Port 	 string
	User 	 string
	Password string
	Dbname   string
}

const table_creation = `
	CREATE TABLE IF NOT EXISTS todo(
		id SERIAL PRIMARY KEY,
		title VARCHAR(100) NOT NULL,
		description  VARCHAR(200)
	)
`

func (app *App) Initialize(connectionDb ConnectionDb) {
	connectionStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			connectionDb.Host, connectionDb.Port, connectionDb.User,
			connectionDb.Password, connectionDb.Dbname)

	var err error
	app.Db, err = sql.Open("postgres", connectionStr)
	if err != nil {
		log.Fatalf("Error loading database %v", err)
	}

	if _, err = app.Db.Exec(table_creation); err != nil {
		log.Fatalf("Error creating table %v", err)
	}
	app.mux = mux.NewRouter()

	todoRepository := repository.NewTodoRepository(app.Db)
	todoHandler := handler.NewTodoHandler(todoRepository)

	app.mux.HandleFunc("/todo", todoHandler.AddTodo).Methods(http.MethodPost)
	app.mux.HandleFunc("/todo", todoHandler.GetTodos).Methods(http.MethodGet)
	app.mux.HandleFunc("/todo/{id}", todoHandler.GetTodo).Methods(http.MethodGet)
}

func (app *App) Run(port string) {
	log.Fatal(http.ListenAndServe(port, app.mux))
}

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"todo/internal/database"
	"todo/internal/handlers"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	databaseUrl := os.Getenv("DATABASE_URL")
	serverPort := os.Getenv("SERVER_PORT")

	if databaseUrl == "" || serverPort == "" {
		log.Fatal("DATABASE_URL or SERVER_PORT not set")
	}

	log.Printf("Starting server on port %s", serverPort)

	db, err := database.Connect(databaseUrl)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	fmt.Println("Connected to database")

	taskStore := database.NewTaskStore(db)
	handlersTask := handlers.NewTodoHandler(taskStore)

	mux := http.NewServeMux()

	mux.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlersTask.GetAllTodos(w, r)
		case http.MethodPost:
			handlersTask.CreateTodo(w, r)
		default:
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/tasks/", taskIdHandler(handlersTask))

	loggedMux := loggingMiddleware(mux)
	serverAdr := ":" + serverPort

	err = http.ListenAndServe(serverAdr, loggedMux)
	if err != nil {
		log.Fatal(err)
	}
}

func taskIdHandler(handler *handlers.TodoHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetTodoById(w, r)
		case http.MethodDelete:
			handler.DeleteTodoById(w, r)
		case http.MethodPut:
			handler.UpdateTodoById(w, r)
		default:
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.Method, r.URL.Path, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

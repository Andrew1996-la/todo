package handlers

import (
	"log"
	"net/http"
)

type HTTPServer struct {
	todoHandler *TodoHandler
	port        string
}

func NewHTTPServer(todoHandler *TodoHandler, port string) *HTTPServer {
	return &HTTPServer{
		todoHandler: todoHandler,
		port:        port,
	}
}

func (s *HTTPServer) Start() {
	mux := http.NewServeMux()

	mux.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			s.todoHandler.GetAllTodos(w, r)
		case http.MethodPost:
			s.todoHandler.CreateTodo(w, r)
		default:
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/tasks/", taskIdHandler(s.todoHandler))

	loggedMux := loggingMiddleware(mux)
	serverAdr := ":" + s.port

	err := http.ListenAndServe(serverAdr, loggedMux)
	if err != nil {
		log.Fatal(err)
	}
}

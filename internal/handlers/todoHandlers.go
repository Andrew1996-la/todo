package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"todo/internal/database"
	"todo/internal/models"
)

type TodoHandler struct {
	todoStore *database.TaskStore
}

func NewTodoHandler(todoStore *database.TaskStore) *TodoHandler {
	return &TodoHandler{todoStore: todoStore}
}

func responseWithJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func responseWithError(w http.ResponseWriter, statusCode int, err error) {
	responseWithJSON(w, statusCode, map[string]string{"error": err.Error()})
}

func findId(w http.ResponseWriter, r *http.Request) (int, error) {
	pathParts := strings.Split(r.URL.Path, "/")
	stringUrl := pathParts[len(pathParts)-1]

	id, err := strconv.Atoi(stringUrl)
	if err != nil {
		responseWithError(w, http.StatusBadRequest, fmt.Errorf("error converting id to int: %w", err))
		return 0, err
	}

	return id, nil
}

func (h *TodoHandler) GetAllTodos(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.todoStore.GetAll()
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, fmt.Errorf("get all todos: %w", err))
		return
	}

	responseWithJSON(w, http.StatusOK, tasks)
}

func (h *TodoHandler) GetTodoById(w http.ResponseWriter, r *http.Request) {
	id, err := findId(w, r)
	if err != nil {
		responseWithError(w, http.StatusBadRequest, err)
	}

	task, err := h.todoStore.GetById(id)
	if err != nil {
		responseWithError(w, http.StatusNotFound, fmt.Errorf("error getting task by id: %w", err))
		return
	}

	responseWithJSON(w, http.StatusOK, task)
}

func (h *TodoHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	var todo models.CreateTaskInput

	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		responseWithError(w, http.StatusBadRequest, fmt.Errorf("error decoding body: %w", err))
		return
	}

	if strings.TrimSpace(todo.Title) == "" {
		responseWithError(w, http.StatusBadRequest, fmt.Errorf("error empty title"))
		return
	}

	task, err := h.todoStore.Create(todo)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, fmt.Errorf("error creating task: %w", err))
		return
	}

	responseWithJSON(w, http.StatusCreated, task)
}

func (h *TodoHandler) UpdateTodoById(w http.ResponseWriter, r *http.Request) {
	id, err := findId(w, r)
	if err != nil {
		responseWithError(w, http.StatusBadRequest, err)
	}

	var input models.UpdateTaskInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		responseWithError(w, http.StatusBadRequest, fmt.Errorf("error decoding body: %w", err))
		return
	}

	if input.Title != nil && strings.TrimSpace(*input.Title) == "" {
		responseWithError(w, http.StatusBadRequest, fmt.Errorf("error empty title"))
		return
	}

	todo, err := h.todoStore.Update(id, input)

	if err != nil {
		responseWithError(w, http.StatusInternalServerError, fmt.Errorf("error updating task by id: %w", err))
		return
	}

	responseWithJSON(w, http.StatusOK, todo)
}

func (h *TodoHandler) DeleteTodoById(w http.ResponseWriter, r *http.Request) {
	id, err := findId(w, r)
	if err != nil {
		responseWithError(w, http.StatusBadRequest, err)
	}

	err = h.todoStore.Delete(id)

	if err != nil {
		responseWithError(w, http.StatusInternalServerError, fmt.Errorf("error deleting task by id: %w", err))
	}

	responseWithJSON(w, http.StatusNoContent, nil)
}

package database

import (
	"database/sql"
	"fmt"
	"time"
	"todo/internal/models"

	"github.com/jmoiron/sqlx"
)

type TaskStore struct {
	db *sqlx.DB
}

func NewTaskStore(db *sqlx.DB) *TaskStore {
	return &TaskStore{db: db}
}

func (t *TaskStore) GetAll() ([]models.Task, error) {
	var tasks []models.Task

	query := `
		SELECT id, title, description, completed, created_at, updated_at 
		FROM tasks 
		order by created_at desc
	`

	err := t.db.Select(&tasks, query)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch tasks: %w", err)
	}

	return tasks, nil
}

func (t *TaskStore) GetById(id int) (*models.Task, error) {
	var task models.Task

	query := `
		SELECT id, title, description, completed, created_at, updated_at 
		FROM tasks 
		where id = $1
	`

	err := t.db.Get(&task, query, id)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("task with id %d not found", id)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to fetch task: %w", err)
	}

	return &task, nil
}

func (t *TaskStore) Create(input models.CreateTaskInput) (*models.Task, error) {
	var task models.Task

	query := `
		INSERT INTO tasks (title, description, completed, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5) 
		RETURNING id, title, description, completed, created_at, updated_at
	`

	currentTime := time.Now()

	err := t.db.QueryRowx(query, input.Title, input.Description, input.Description, currentTime, currentTime).StructScan(&task)

	if err != nil {
		return nil, fmt.Errorf("failed to insert task: %w", err)
	}

	return &task, nil
}

func (t *TaskStore) Update(id int, input models.UpdateTaskInput) (*models.Task, error) {
	task, err := t.GetById(id)

	if err != nil {
		return nil, fmt.Errorf("failed task with id %d not found: %w", id, err)
	}

	if input.Title != nil {
		task.Title = *input.Title
	}

	if input.Description != nil {
		task.Description = *input.Description
	}
	if input.Completed != nil {
		task.Completed = *input.Completed
	}

	task.CreatedAt = time.Now()

	query := `
		UPDATE tasks
		SET title = $1, description = $2, completed = $3, updated_at = $4
		WHERE id = $5
		RETURNING id, title, description, completed, created_at, updated_at
	`

	var updatedTask models.Task

	errUpdate := t.db.QueryRowx(query, task.Title, task.Description, task.Completed, task.UpdatedAt).StructScan(&updatedTask)
	if errUpdate != nil {
		return nil, fmt.Errorf("failed to update task: %w", err)
	}

	return &updatedTask, nil
}

func (t *TaskStore) Delete(id int) error {
	query := `DELETE FROM tasks WHERE id = $1`

	res, err := t.db.Exec(query, id)

	if err != nil {
		return fmt.Errorf("failed to delete task with id %d: %w", id, err)
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		return fmt.Errorf("failed to delete task with id %d: %w", id, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("task with id %d not found", id)
	}

	return nil
}

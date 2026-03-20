package repository

import (
	"database/sql"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/kubegogo/genvideo/internal/model"
)

type Repository struct {
	db  *sql.DB
	rdb *redis.Client
}

func NewRepository(db *sql.DB, rdb *redis.Client) *Repository {
	return &Repository{db: db, rdb: rdb}
}

func (r *Repository) CreateTask(task *model.Task) error {
	query := `INSERT INTO tasks (type, status, input, progress, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`
	now := time.Now()
	result, err := r.db.Exec(query, task.Type, task.Status, task.Input, task.Progress, now, now)
	if err != nil {
		return err
	}
	id, _ := result.LastInsertId()
	task.ID = id
	task.CreatedAt = now
	task.UpdatedAt = now
	return nil
}

func (r *Repository) GetTask(id int64) (*model.Task, error) {
	task := &model.Task{}
	query := `SELECT id, type, status, input, output, error, progress, created_at, updated_at FROM tasks WHERE id = ?`
	err := r.db.QueryRow(query, id).Scan(&task.ID, &task.Type, &task.Status, &task.Input, &task.Output, &task.Error, &task.Progress, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (r *Repository) UpdateTask(task *model.Task) error {
	query := `UPDATE tasks SET status = ?, output = ?, error = ?, progress = ?, updated_at = ? WHERE id = ?`
	task.UpdatedAt = time.Now()
	_, err := r.db.Exec(query, task.Status, task.Output, task.Error, task.Progress, task.UpdatedAt, task.ID)
	return err
}

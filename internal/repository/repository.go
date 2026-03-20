package repository

import (
	"database/sql"
	"encoding/json"

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

// Task operations
func (r *Repository) CreateTask(task *model.Task) error {
	query := `INSERT INTO tasks (type, status, input, progress, created_at, updated_at) VALUES (?, ?, ?, ?, NOW(), NOW())`
	result, err := r.db.Exec(query, task.Type, task.Status, task.Input, task.Progress)
	if err != nil {
		return err
	}
	id, _ := result.LastInsertId()
	task.ID = id
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
	query := `UPDATE tasks SET status = ?, output = ?, error = ?, progress = ?, updated_at = NOW() WHERE id = ?`
	_, err := r.db.Exec(query, task.Status, task.Output, task.Error, task.Progress, task.ID)
	return err
}

// VideoProvider operations
func (r *Repository) GetVideoProviders() ([]model.VideoProvider, error) {
	query := `SELECT id, platform, cookie, status FROM video_providers`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var providers []model.VideoProvider
	for rows.Next() {
		var p model.VideoProvider
		if err := rows.Scan(&p.ID, &p.Platform, &p.Cookie, &p.Status); err != nil {
			return nil, err
		}
		providers = append(providers, p)
	}
	return providers, nil
}

func (r *Repository) SaveVideoProvider(p *model.VideoProvider) error {
	query := `INSERT INTO video_providers (platform, cookie, status) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE cookie = ?, status = ?`
	_, err := r.db.Exec(query, p.Platform, p.Cookie, p.Status, p.Cookie, p.Status)
	return err
}

// AIProvider operations
func (r *Repository) GetAIProviders() ([]model.AIProvider, error) {
	query := `SELECT id, type, api_key, base_url, is_active FROM ai_providers`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var providers []model.AIProvider
	for rows.Next() {
		var p model.AIProvider
		if err := rows.Scan(&p.ID, &p.Type, &p.APIKey, &p.BaseURL, &p.IsActive); err != nil {
			return nil, err
		}
		providers = append(providers, p)
	}
	return providers, nil
}

func (r *Repository) SaveAIProvider(p *model.AIProvider) error {
	query := `INSERT INTO ai_providers (type, api_key, base_url, is_active) VALUES (?, ?, ?, ?) ON DUPLICATE KEY UPDATE api_key = ?, base_url = ?, is_active = ?`
	_, err := r.db.Exec(query, p.Type, p.APIKey, p.BaseURL, p.IsActive, p.APIKey, p.BaseURL, p.IsActive)
	return err
}

// OSSConfig operations
func (r *Repository) GetOSSConfig() (*model.OSSConfig, error) {
	query := `SELECT id, endpoint, access_key, secret_key, bucket, is_active FROM oss_config LIMIT 1`
	cfg := &model.OSSConfig{}
	err := r.db.QueryRow(query).Scan(&cfg.ID, &cfg.Endpoint, &cfg.AccessKey, &cfg.SecretKey, &cfg.Bucket, &cfg.IsActive)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func (r *Repository) SaveOSSConfig(cfg *model.OSSConfig) error {
	query := `INSERT INTO oss_config (endpoint, access_key, secret_key, bucket, is_active) VALUES (?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE endpoint = ?, access_key = ?, secret_key = ?, bucket = ?, is_active = ?`
	_, err := r.db.Exec(query, cfg.Endpoint, cfg.AccessKey, cfg.SecretKey, cfg.Bucket, cfg.IsActive, cfg.Endpoint, cfg.AccessKey, cfg.SecretKey, cfg.Bucket, cfg.IsActive)
	return err
}

// Cache operations using Redis
func (r *Repository) CacheTaskResult(ctx context.Context, taskID int64, result interface{}) error {
	key := fmt.Sprintf("task_result:%d", taskID)
	data, err := json.Marshal(result)
	if err != nil {
		return err
	}
	return r.rdb.Set(ctx, key, data, 0).Err()
}

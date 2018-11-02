package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/plopezm/cloud-kaiser/types"
	"strings"
	"time"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgres(url string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	return &PostgresRepository{
		db,
	}, nil
}

func (r *PostgresRepository) Close() {
	r.db.Close()
}

func (r *PostgresRepository) InsertTask(ctx context.Context, task types.JobTask) error {
	var onSuccessName string
	var onSuccessVersion string
	var onFailureName string
	var onFailureVersion string

	if len(task.OnSuccess) > 0 {
		taskIdentifiers := strings.Split(task.OnSuccess, ":")
		if len(taskIdentifiers) > 1 {
			onSuccessName = taskIdentifiers[0]
			onSuccessVersion = taskIdentifiers[1]
		}
	}
	if len(task.OnFailure) > 0 {
		taskIdentifiers := strings.Split(task.OnFailure, ":")
		if len(taskIdentifiers) > 1 {
			onFailureName = taskIdentifiers[0]
			onFailureVersion = taskIdentifiers[1]
		}
	}

	_, err := r.db.Exec(`INSERT INTO tasks(name, version, created_at, script, on_success_name, on_success_version, on_failure_name, on_failure_version)
								VALUES($1, $2, $3, $4, $5, $6, $7, $8)`,
		task.Name, task.Version, time.Now(), task.Script, onSuccessName, onSuccessVersion, onFailureName, onFailureVersion)
	return err
}

func (r *PostgresRepository) ListTasks(ctx context.Context, offset uint64, limit uint64) ([]types.JobTask, error) {
	rows, err := r.db.Query("SELECT * FROM tasks ORDER BY name DESC, version DESC OFFSET $1 LIMIT $2", offset, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []types.JobTask{}
	for rows.Next() {
		task := types.JobTask{}
		var onSuccessName string
		var onSuccessVersion string
		var onFailureName string
		var onFailureVersion string
		if err = rows.Scan(&task.Name, &task.Version, &task.CreatedAt, &task.Script, &onSuccessName, &onSuccessVersion, &onFailureName, &onFailureVersion); err == nil {
			task.OnSuccess = fmt.Sprintf("%s:%s", onSuccessName, onSuccessVersion)
			task.OnFailure = fmt.Sprintf("%s:%s", onFailureName, onFailureVersion)
			tasks = append(tasks, task)
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}

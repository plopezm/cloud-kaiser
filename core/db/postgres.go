package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"github.com/plopezm/cloud-kaiser/core/types"
	"log"
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
	var onSuccessName *string
	var onSuccessVersion *string
	var onFailureName *string
	var onFailureVersion *string

	if len(task.OnSuccess) > 0 {
		taskIdentifiers := strings.Split(task.OnSuccess, ":")
		if len(taskIdentifiers) > 1 {
			onSuccessName = &taskIdentifiers[0]
			onSuccessVersion = &taskIdentifiers[1]
		}
	}
	if len(task.OnFailure) > 0 {
		taskIdentifiers := strings.Split(task.OnFailure, ":")
		if len(taskIdentifiers) > 1 {
			onFailureName = &taskIdentifiers[0]
			onFailureVersion = &taskIdentifiers[1]
		}
	}

	_, err := r.db.Exec(`INSERT INTO tasks(name, version, created_at, script, on_success_name, on_success_version, on_failure_name, on_failure_version)
								VALUES($1, $2, $3, $4, $5, $6, $7, $8)`,
		task.Name, task.Version, time.Now(), task.Script, onSuccessName, onSuccessVersion, onFailureName, onFailureVersion)
	return err
}

func (r *PostgresRepository) ListTasks(ctx context.Context, offset uint64, limit uint64) ([]types.JobTask, error) {
	rows, err := r.db.Query(
		`SELECT name, version, created_at, script, on_success_name, on_success_version, on_failure_name, on_failure_version 
				FROM tasks ORDER BY name DESC, version DESC OFFSET $1 LIMIT $2`, offset, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []types.JobTask{}
	for rows.Next() {
		task := types.JobTask{}
		task.Script = new(string)
		var onSuccessName sql.NullString
		var onSuccessVersion sql.NullString
		var onFailureName sql.NullString
		var onFailureVersion sql.NullString
		if err = rows.Scan(&task.Name, &task.Version, &task.CreatedAt, task.Script, &onSuccessName, &onSuccessVersion, &onFailureName, &onFailureVersion); err == nil {
			if onSuccessName.Valid && onSuccessVersion.Valid {
				task.OnSuccess = fmt.Sprintf("%s:%s", onSuccessName.String, onSuccessVersion.String)
			}
			if onFailureName.Valid && onFailureVersion.Valid {
				task.OnFailure = fmt.Sprintf("%s:%s", onFailureName.String, onFailureVersion.String)
			}
			tasks = append(tasks, task)
		} else {
			return nil, err
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *PostgresRepository) FindTaskByName(ctx context.Context, name string) ([]types.JobTask, error) {
	rows, err := r.db.Query(
		`SELECT name, version, created_at, script, on_success_name, on_success_version, on_failure_name, on_failure_version 
				FROM tasks WHERE name = $1 ORDER BY name DESC, version DESC`, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []types.JobTask{}
	for rows.Next() {
		task := types.JobTask{}
		task.Script = new(string)
		var onSuccessName sql.NullString
		var onSuccessVersion sql.NullString
		var onFailureName sql.NullString
		var onFailureVersion sql.NullString
		if err = rows.Scan(&task.Name, &task.Version, &task.CreatedAt, task.Script, &onSuccessName, &onSuccessVersion, &onFailureName, &onFailureVersion); err == nil {
			if onSuccessName.Valid && onSuccessVersion.Valid {
				task.OnSuccess = fmt.Sprintf("%s:%s", onSuccessName.String, onSuccessVersion.String)
			}
			if onFailureName.Valid && onFailureVersion.Valid {
				task.OnFailure = fmt.Sprintf("%s:%s", onFailureName.String, onFailureVersion.String)
			}
			tasks = append(tasks, task)
		} else {
			return nil, err
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *PostgresRepository) FindTaskByNameAndVersion(ctx context.Context, name string, version string) (*types.JobTask, error) {
	rows, err := r.db.Query(
		`SELECT name, version, created_at, script, on_success_name, on_success_version, on_failure_name, on_failure_version 
				FROM tasks WHERE name = $1 AND version = $2 ORDER BY name DESC, version DESC`, name, version)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var task *types.JobTask
	if rows.Next() {
		task = &types.JobTask{}
		task.Script = new(string)
		var onSuccessName sql.NullString
		var onSuccessVersion sql.NullString
		var onFailureName sql.NullString
		var onFailureVersion sql.NullString
		if err = rows.Scan(&task.Name, &task.Version, &task.CreatedAt, task.Script, &onSuccessName, &onSuccessVersion, &onFailureName, &onFailureVersion); err == nil {
			if onSuccessName.Valid && onSuccessVersion.Valid {
				task.OnSuccess = fmt.Sprintf("%s:%s", onSuccessName.String, onSuccessVersion.String)
			}
			if onFailureName.Valid && onFailureVersion.Valid {
				task.OnFailure = fmt.Sprintf("%s:%s", onFailureName.String, onFailureVersion.String)
			}
		} else {
			return nil, err
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return task, nil
}

func (r *PostgresRepository) InsertJobArgument(ctx context.Context, job *types.Job, argument types.JobArgs) error {
	if job == nil {
		return errors.New("Job instance cannot be null")
	}
	var query = `INSERT INTO arguments(name, value, job_name, job_version )
								VALUES($1, $2, $3, $4)`

	var err error
	tx, ok := ctx.Value("tx").(*sql.Tx)
	if ok {
		_, err = tx.Exec(query, argument.Name, argument.Value, job.Name, job.Version)
	} else {
		_, err = r.db.Exec(query, argument.Name, argument.Value, job.Name, job.Version)
	}

	return err
}

func (r *PostgresRepository) InsertJob(ctx context.Context, job types.Job) error {
	var duration *string
	if job.Activation.Duration != "" {
		duration = &job.Activation.Duration
	}

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(`INSERT INTO jobs(name, version, created_at, activation_type, duration)
								VALUES($1, $2, $3, $4, $5)`,
		job.Name, job.Version, time.Now(), job.Activation.Type, duration)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Insert tasks relations
	for _, jobtask := range job.Tasks {
		task, err := r.FindTaskByNameAndVersion(ctx, jobtask.Name, jobtask.Version)
		if err != nil {
			tx.Rollback()
			return errors.Wrap(err, fmt.Sprintf("Job creation failed beacuse of a wrong task asignment. TASK [%s, %s]", task.Name, task.Version))
		}
		err = r.insertJobTaskRelation(tx, &job, task)
		if err != nil {
			tx.Rollback()
			return errors.Wrap(err,
				fmt.Sprintf("Job creation failed beacuse job-task relation could not be created. JOB [%s, %s] TASK [%s, %s]",
					job.Name, job.Version, task.Name, task.Version))
		}
	}

	contextWithTx := context.WithValue(ctx, "tx", tx)
	//Insert arguments
	for _, jobArg := range job.Parameters {
		// if someone fails the argument is not added but the job insertion does not fail
		err := r.InsertJobArgument(contextWithTx, &job, jobArg)
		if err != nil {
			log.Println("Error creating argument", err)
		}
	}

	tx.Commit()
	return nil
}

func (r *PostgresRepository) insertJobTaskRelation(tx *sql.Tx, job *types.Job, task *types.JobTask) error {
	_, err := tx.Exec(`INSERT INTO jobs_tasks(job_name, job_version, task_name, task_version)
								VALUES($1, $2, $3, $4)`,
		job.Name, job.Version, task.Name, task.Version)
	return err
}

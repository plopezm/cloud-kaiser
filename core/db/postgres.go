package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/plopezm/cloud-kaiser/core/types"
)

type DBManager interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
}

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

func (r *PostgresRepository) Tx(ctx context.Context, opts *sql.TxOptions, txF TransactionFunction) error {
	var tx *sql.Tx
	var err error
	var ok bool

	tx, ok = ctx.Value(ContextTX).(*sql.Tx)
	if tx == nil || !ok {
		tx, err = r.db.BeginTx(ctx, opts)
		if err != nil {
			return err
		}
	}

	err = txF(tx)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}
	return tx.Commit()
}

func (r *PostgresRepository) InsertTask(ctx context.Context, task types.Task) error {
	var dbm DBManager
	tx, ok := ctx.Value(ContextTX).(*sql.Tx)
	if tx != nil && ok {
		dbm = tx
	} else {
		dbm = r.db
	}

	_, err := dbm.Exec(`INSERT INTO tasks(name, version, created_at, script)
								VALUES($1, $2, $3, $4)`,
		task.Name, task.Version, time.Now(), task.Script)
	return err
}

func (r *PostgresRepository) ListTasks(ctx context.Context, offset uint64, limit uint64) ([]types.Task, error) {
	rows, err := r.db.Query(
		`SELECT name, version, created_at, script
				FROM tasks ORDER BY name DESC, version DESC OFFSET $1 LIMIT $2`, offset, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []types.Task{}
	for rows.Next() {
		task := types.Task{}
		if err = rows.Scan(&task.Name, &task.Version, &task.CreatedAt, &task.Script); err == nil {
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

func (r *PostgresRepository) FindTaskByName(ctx context.Context, name string) ([]types.Task, error) {
	rows, err := r.db.Query(
		`SELECT name, version, created_at, script 
				FROM tasks WHERE name = $1 ORDER BY name DESC, version DESC`, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []types.Task{}
	for rows.Next() {
		task := types.Task{}
		if err = rows.Scan(&task.Name, &task.Version, &task.CreatedAt, &task.Script); err == nil {
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

func (r *PostgresRepository) FindTaskByNameAndVersion(ctx context.Context, name string, version string) (*types.Task, error) {
	rows, err := r.db.Query(
		`SELECT name, version, created_at, script 
				FROM tasks WHERE name = $1 AND version = $2 ORDER BY name DESC, version DESC`, name, version)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var task *types.Task
	if rows.Next() {
		task = &types.Task{}
		err = rows.Scan(&task.Name, &task.Version, &task.CreatedAt, &task.Script)
		if err != nil {
			return nil, err
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return task, nil
}

func (r *PostgresRepository) InsertJob(ctx context.Context, job *types.Job) error {
	return r.Tx(ctx, &sql.TxOptions{Isolation: sql.LevelReadUncommitted}, func(tx *sql.Tx) error {
		// Insert job and tasks relations
		if err := r.createJobInTx(ctx, tx, job); err != nil {
			return err
		}

		//Insert arguments
		for _, jobArg := range job.Parameters {
			// if someone fails the argument is not added but the job insertion does not fail
			err := r.AddJobArgument(ctx, job, jobArg)
			if err != nil {
				log.Println("Error creating argument", err)
			}
		}
		return nil
	})
}

func (r *PostgresRepository) createJobInTx(ctx context.Context, tx *sql.Tx, job *types.Job) error {
	var duration *string
	if job.Activation.Duration != "" {
		duration = &job.Activation.Duration
	}
	_, err := tx.Exec(`INSERT INTO jobs(name, version, created_at, activation_type, duration)
								VALUES($1, $2, $3, $4, $5)`,
		job.Name, job.Version, time.Now(), job.Activation.Type, duration)
	if err != nil {
		//tx.Rollback()
		return err
	}

	// Insert tasks relations
	for taskName, jobtask := range job.Tasks {
		task, err := r.FindTaskByNameAndVersion(ctx, taskName, jobtask.Version)
		if task == nil {
			//tx.Rollback()
			return errors.New(fmt.Sprintf("Job creation failed beacuse the task referenced does not exist. TASK [%s, %s]", taskName, jobtask.Version))
		}
		if err != nil {
			//tx.Rollback()
			return errors.Wrap(err, fmt.Sprintf("Job creation failed beacuse of a wrong task asignment. TASK [%s, %s]", task.Name, task.Version))
		}
		jobtask.Name = task.Name
		jobtask.Script = task.Script
		err = r.insertJobTaskRelation(tx, job, &jobtask)
		if err != nil {
			//tx.Rollback()
			return errors.Wrap(err,
				fmt.Sprintf("Job creation failed beacuse job-task relation could not be created. JOB [%s, %s] TASK [%s, %s]",
					job.Name, job.Version, task.Name, task.Version))
		}
		job.Tasks[taskName] = jobtask
	}
	return nil
}

func (r *PostgresRepository) insertJobTaskRelation(tx *sql.Tx, job *types.Job, task *types.JobTask) error {
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
	var isEntryPoint = 0
	if job.Entrypoint == task.Name {
		isEntryPoint = 1
	}

	_, err := tx.Exec(`INSERT INTO jobs_tasks(job_name, job_version, task_name, task_version, on_success_name, on_success_version, on_failure_name, on_failure_version, entrypoint)
								VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		job.Name, job.Version, task.Name, task.Version, onSuccessName, onSuccessVersion, onFailureName, onFailureVersion, isEntryPoint)
	return err
}

func (r *PostgresRepository) AddJobArgument(ctx context.Context, job *types.Job, argument types.JobArgs) error {
	if job == nil {
		return errors.New("Job instance cannot be null")
	}

	var dbm DBManager
	tx, ok := ctx.Value(ContextTX).(*sql.Tx)
	if tx != nil && ok {
		dbm = tx
	} else {
		dbm = r.db
	}

	var query = `INSERT INTO arguments(name, value, job_name, job_version )
								VALUES($1, $2, $3, $4)`

	_, err := dbm.Exec(query, argument.Name, argument.Value, job.Name, job.Version)

	return err
}

func (r *PostgresRepository) ListJobs(ctx context.Context, offset uint64, limit uint64) ([]types.Job, error) {
	rows, err := r.db.Query(
		`SELECT name, version, created_at, activation_type, duration 
				FROM jobs ORDER BY name DESC, version DESC OFFSET $1 LIMIT $2`, offset, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	jobs := []types.Job{}
	for rows.Next() {
		var job types.Job
		var duration sql.NullString
		if err = rows.Scan(&job.Name, &job.Version, &job.CreatedAt, &job.Activation.Type, &duration); err == nil {
			if duration.Valid {
				job.Activation.Duration = duration.String
			}
			job.Tasks, err = r.findJobTasks(&job)
			if err != nil {
				log.Printf("Error fetching job tasks. JOB [%s, %s]\n", job.Name, job.Version)
			}
			jobs = append(jobs, job)
		} else {
			return nil, err
		}
	}

	return jobs, nil
}

func (r *PostgresRepository) FindJobByName(ctx context.Context, name string) ([]types.Job, error) {
	rows, err := r.db.Query(
		`SELECT name, version, created_at, activation_type, duration 
				FROM jobs WHERE name = $1 ORDER BY name DESC, version DESC`, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	jobs := []types.Job{}
	for rows.Next() {
		var job types.Job
		var duration sql.NullString
		if err = rows.Scan(&job.Name, &job.Version, &job.CreatedAt, &job.Activation.Type, &duration); err == nil {
			if duration.Valid {
				job.Activation.Duration = duration.String
			}
			job.Tasks, err = r.findJobTasks(&job)
			if err != nil {
				log.Printf("Error fetching job tasks. JOB [%s, %s]\n", job.Name, job.Version)
			}
			jobs = append(jobs, job)
		} else {
			return nil, err
		}
	}

	return jobs, nil
}

func (r *PostgresRepository) FindJobByNameAndVersion(ctx context.Context, name string, version string) (*types.Job, error) {
	rows, err := r.db.Query(
		`SELECT name, version, created_at, activation_type, duration 
				FROM jobs WHERE name = $1 AND version = $2`, name, version)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var job = new(types.Job)
	if rows.Next() {
		var duration sql.NullString
		if err = rows.Scan(&job.Name, &job.Version, &job.CreatedAt, &job.Activation.Type, &duration); err == nil {
			if duration.Valid {
				job.Activation.Duration = duration.String
			}
			job.Tasks, err = r.findJobTasks(job)
			if err != nil {
				log.Printf("Error fetching job tasks. JOB [%s, %s]\n", job.Name, job.Version)
			}
		} else {
			return nil, err
		}
	}

	return job, nil
}

func (r *PostgresRepository) findJobTasks(job *types.Job) (map[string]types.JobTask, error) {
	rows, err := r.db.Query(
		`SELECT t.name, t.version, t.created_at, t.script, jt.on_success_name, jt.on_success_version, jt.on_failure_name, jt.on_failure_version, jt.entrypoint 
				FROM tasks t, jobs_tasks jt WHERE t.name = jt.task_name AND t.version = jt.task_version
				AND jt.job_name = $1 AND jt.job_version = $2 ORDER BY name DESC, version DESC`, job.Name, job.Version)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make(map[string]types.JobTask)
	for rows.Next() {
		task := types.JobTask{}
		var onSuccessName sql.NullString
		var onSuccessVersion sql.NullString
		var onFailureName sql.NullString
		var onFailureVersion sql.NullString
		var isEntryPoint int8
		if err = rows.Scan(&task.Name, &task.Version, &task.CreatedAt, &task.Script, &onSuccessName, &onSuccessVersion, &onFailureName, &onFailureVersion, &isEntryPoint); err == nil {
			if onSuccessName.Valid && onSuccessVersion.Valid {
				task.OnSuccess = fmt.Sprintf("%s:%s", onSuccessName.String, onSuccessVersion.String)
			}
			if onFailureName.Valid && onFailureVersion.Valid {
				task.OnFailure = fmt.Sprintf("%s:%s", onFailureName.String, onFailureVersion.String)
			}
			tasks[task.Name] = task
			if isEntryPoint == 1 {
				job.Entrypoint = task.Name
			}
		} else {
			return nil, err
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}

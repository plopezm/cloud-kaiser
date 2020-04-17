package db

import (
	"context"
	"database/sql"

	"github.com/plopezm/cloud-kaiser/core/types"
)

const (
	//ContextTX The transaction key for the context
	ContextTX = "tx"
)

//TransactionFunction The action to complete in order to commit a transaction successfuly
type TransactionFunction func(*sql.Tx) error

//Repository The repository interface
type Repository interface {
	Close()
	//Tx Creates a transaction scope
	Tx(context.Context, *sql.TxOptions, TransactionFunction) error
	// InsertTask Creates a new task
	InsertTask(ctx context.Context, job types.Task) error
	// ListTasks Returns a paginated list of tasks
	ListTasks(ctx context.Context, offset uint64, limit uint64) ([]types.Task, error)
	// FindTaskByName Returns all versions of a task by name
	FindTaskByName(ctx context.Context, name string) ([]types.Task, error)
	// FindTaskByNameAndVersion Returns a version of a task
	FindTaskByNameAndVersion(ctx context.Context, name string, version string) (*types.Task, error)

	// InsertJob Creates a new job
	InsertJob(ctx context.Context, job *types.Job) error
	// ListJobs Returns a paginated list of jobs
	ListJobs(ctx context.Context, offset uint64, limit uint64) ([]types.Job, error)
	// FindJobByName Returns all versions of a job by name
	FindJobByName(ctx context.Context, name string) ([]types.Job, error)
	// FindJobByNameAndVersion Returns a version of a job
	FindJobByNameAndVersion(ctx context.Context, name string, version string) (*types.Job, error)

	// AddJobArgument Adds a new argument to an existing job
	AddJobArgument(ctx context.Context, job *types.Job, argument types.JobArgs) error
}

var impl Repository

//SetRepository Sets the a repository implementation
func SetRepository(repository Repository) {
	impl = repository
}

//Close Closes DB connection
func Close() {
	impl.Close()
}

//Tx Creates a transaction scope
func Tx(ctx context.Context, opts *sql.TxOptions, txF TransactionFunction) error {
	return impl.Tx(ctx, opts, txF)
}

//InsertTask Creates a new task
func InsertTask(ctx context.Context, job types.Task) error {
	return impl.InsertTask(ctx, job)
}

//ListTasks Returns a paginated list of tasks
func ListTasks(ctx context.Context, offset uint64, limit uint64) ([]types.Task, error) {
	return impl.ListTasks(ctx, offset, limit)
}

//FindTaskByName Returns all versions of a task by name
func FindTaskByName(ctx context.Context, name string) ([]types.Task, error) {
	return impl.FindTaskByName(ctx, name)
}

//FindTaskByNameAndVersion Returns a version of a task
func FindTaskByNameAndVersion(ctx context.Context, name string, version string) (*types.Task, error) {
	return impl.FindTaskByNameAndVersion(ctx, name, version)
}

func InsertJobArgument(ctx context.Context, job *types.Job, argument types.JobArgs) error {
	return impl.AddJobArgument(ctx, job, argument)
}

func InsertJob(ctx context.Context, job *types.Job) error {
	return impl.InsertJob(ctx, job)
}

func ListJobs(ctx context.Context, page uint64, limit uint64) ([]types.Job, error) {
	return impl.ListJobs(ctx, page, limit)
}

func FindJobByName(ctx context.Context, name string) ([]types.Job, error) {
	return impl.FindJobByName(ctx, name)
}

func FindJobByNameAndVersion(ctx context.Context, name string, version string) (*types.Job, error) {
	return impl.FindJobByNameAndVersion(ctx, name, version)
}

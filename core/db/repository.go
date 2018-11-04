package db

import (
	"context"
	"github.com/plopezm/cloud-kaiser/core/types"
)

type Repository interface {
	Close()
	// InsertTask Creates a new task
	InsertTask(ctx context.Context, job types.JobTask) error
	// ListTasks Returns a paginated list of tasks
	ListTasks(ctx context.Context, offset uint64, limit uint64) ([]types.JobTask, error)
	// FindTaskByName Returns all versions of a task by name
	FindTaskByName(ctx context.Context, name string) ([]types.JobTask, error)
	// FindTaskByNameAndVersion Returns a version of a task
	FindTaskByNameAndVersion(ctx context.Context, name string, version string) (*types.JobTask, error)

	InsertJobArgument(ctx context.Context, job *types.Job, argument types.JobArgs) error

	InsertJob(ctx context.Context, job types.Job) error
	ListJobs(ctx context.Context, offset uint64, limit uint64) ([]types.Job, error)
}

var impl Repository

func SetRepository(repository Repository) {
	impl = repository
}

func Close() {
	impl.Close()
}

func InsertTask(ctx context.Context, job types.JobTask) error {
	return impl.InsertTask(ctx, job)
}

func ListTasks(ctx context.Context, offset uint64, limit uint64) ([]types.JobTask, error) {
	return impl.ListTasks(ctx, offset, limit)
}

func FindTaskByName(ctx context.Context, name string) ([]types.JobTask, error) {
	return impl.FindTaskByName(ctx, name)
}

func FindTaskByNameAndVersion(ctx context.Context, name string, version string) (*types.JobTask, error) {
	return impl.FindTaskByNameAndVersion(ctx, name, version)
}

func InsertJobArgument(ctx context.Context, job *types.Job, argument types.JobArgs) error {
	return impl.InsertJobArgument(ctx, job, argument)
}

func InsertJob(ctx context.Context, job types.Job) error {
	return impl.InsertJob(ctx, job)
}

func ListJobs(ctx context.Context, page uint64, limit uint64) ([]types.Job, error) {
	return impl.ListJobs(ctx, page, limit)
}

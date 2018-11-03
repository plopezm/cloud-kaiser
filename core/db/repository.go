package db

import (
	"context"
	"github.com/plopezm/cloud-kaiser/core/types"
)

type Repository interface {
	Close()
	InsertTask(ctx context.Context, job types.JobTask) error
	ListTasks(ctx context.Context, page uint64, limit uint64) ([]types.JobTask, error)
	FindTaskByName(ctx context.Context, name string) ([]types.JobTask, error)
	FindTaskByNameAndVersion(ctx context.Context, name string, version string) (*types.JobTask, error)
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

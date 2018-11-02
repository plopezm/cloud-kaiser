package repository

import (
	"context"
	"github.com/plopezm/cloud-kaiser/types"
)

type Repository interface {
	Close()
	InsertTask(ctx context.Context, job types.JobTask) error
	ListTasks(ctx context.Context, page uint64, limit uint64) ([]types.JobTask, error)
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

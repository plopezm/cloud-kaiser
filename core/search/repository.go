package search

import (
	"context"

	"github.com/plopezm/cloud-kaiser/core/types"
)

type Repository interface {
	Close()
	InsertTask(ctx context.Context, job types.Task) error
	FindTasks(ctx context.Context, query string, page uint64, limit uint64) ([]types.Task, error)
	InsertLog(ctx context.Context, taskExecutionLog types.TaskExecutionLog) error
	FindLogs(ctx context.Context, query string, field string, offset uint64, limit uint64) ([]types.TaskExecutionLog, error)
}

var impl Repository

func SetRepository(repository Repository) {
	impl = repository
}

func Close() {
	impl.Close()
}

func IsConfigured() bool {
	return impl != nil
}

func InsertTask(ctx context.Context, job types.Task) error {
	return impl.InsertTask(ctx, job)
}

func FindTasks(ctx context.Context, query string, offset uint64, limit uint64) ([]types.Task, error) {
	return impl.FindTasks(ctx, query, offset, limit)
}

func InsertLog(ctx context.Context, taskExecutionLog types.TaskExecutionLog) error {
	return impl.InsertLog(ctx, taskExecutionLog)
}

func FindLogs(ctx context.Context, query string, field string, offset uint64, limit uint64) ([]types.TaskExecutionLog, error) {
	return impl.FindLogs(ctx, query, field, offset, limit)
}

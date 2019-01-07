package search

import (
	"context"
	"github.com/plopezm/cloud-kaiser/core/types"
)

type Repository interface {
	Close()
	InsertTask(ctx context.Context, job types.Task) error
	FindTasks(ctx context.Context, query string, page uint64, limit uint64) ([]types.Task, error)
	InsertLog(ctx context.Context, jobname string, jobversion string, taskname string, taskversion string, logline string) error
	FindLogs(ctx context.Context, query string, fields []string, offset uint64, limit uint64) ([]interface{}, error)
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

func InsertLog(ctx context.Context, jobname string, jobversion string, taskname string, taskversion string, logline string) error {
	return impl.InsertLog(ctx, jobname, jobversion, taskname, taskversion, logline)
}

func FindLogs(ctx context.Context, query string, fields []string, offset uint64, limit uint64) ([]interface{}, error) {
	return impl.FindLogs(ctx, query, fields, offset, limit)
}

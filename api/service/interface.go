package service

import (
	"context"
	"todotree/entity"
	"todotree/store"
)

//go:generate go run github.com/matryer/moq -out moq_test.go . TaskAdder TaskListener
type TaskAdder interface {
	AddTask(ctx context.Context, db store.Execer, t *entity.Task) error
}

type TaskLister interface {
	ListTasks(ctx context.Context, db store.Queryer) (entity.Tasks, error)
}

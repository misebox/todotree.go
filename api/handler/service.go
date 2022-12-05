package handler

import (
	"context"
	"todotree/entity"
)

//go:generate go run github.com/matryer/moq -out moq_test.go . ListTasksService AddTaskService RegisterUserService LoginService
type ListTasksService interface {
	ListTasks(ctx context.Context) (entity.Tasks, error)
}

type AddTaskService interface {
	AddTask(ctx context.Context, title string, parentID *entity.TaskID) (*entity.Task, error)
}

type RegisterUserService interface {
	RegisterUser(ctx context.Context, name, email, password, role string) (*entity.User, error)
}
type LoginService interface {
	Login(ctx context.Context, name, pw string) (string, error)
}

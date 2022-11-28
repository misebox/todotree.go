package service

import (
	"context"
	"fmt"
	"todotree/entity"
	"todotree/store"
)

type AddTask struct {
	DB   store.Execer
	Repo TaskAdder
}

func (at *AddTask) AddTask(ctx context.Context, title string) (*entity.Task, error) {
	t := &entity.Task{
		Title:  title,
		Status: entity.TaskStatusTodo,
	}
	err := at.Repo.AddTask(ctx, at.DB, t)
	if err != nil {
		return nil, fmt.Errorf("failed to register: %w", err)
	}
	return t, nil
}

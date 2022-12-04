package service

import (
	"context"
	"fmt"
	"todotree/auth"
	"todotree/entity"
	"todotree/store"
)

type ListTask struct {
	DB   store.Queryer
	Repo TaskLister
}

func (lt *ListTask) ListTasks(ctx context.Context) (entity.Tasks, error) {
	userID, ok := auth.GetUserID(ctx)
	if !ok {
		return nil, fmt.Errorf("user_id not found")
	}
	ts, err := lt.Repo.ListTasks(ctx, lt.DB, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list tasks: %w", err)
	}
	return ts, nil
}

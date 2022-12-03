package store

import (
	"context"
	"fmt"
	"todotree/auth"
	"todotree/entity"
)

func (r *Repository) ListTasks(
	ctx context.Context, db Queryer,
) (entity.Tasks, error) {
	user_id, ok := auth.GetUserID(ctx)
	if !ok {
		return nil, fmt.Errorf("user_id not found")
	}
	tasks := entity.Tasks{}
	sql := `
		SELECT
			id, title, status,
			created, modified
		FROM task
		WHERE user_id = ?;`
	err := db.SelectContext(
		ctx, &tasks, sql,
		user_id,
	)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *Repository) AddTask(
	ctx context.Context, db Execer, t *entity.Task,
) error {
	user_id, ok := auth.GetUserID(ctx)
	if !ok {
		return fmt.Errorf("user_id not found")
	}
	t.Created = r.Clocker.Now()
	t.Modified = r.Clocker.Now()
	sql := `INSERT INTO task
	(user_id, title, status, created, modified)
	VALUES (?, ?, ?, ?, ?);`
	result, err := db.ExecContext(
		ctx, sql,
		user_id, t.Title, t.Status, t.Created, t.Modified,
	)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	t.ID = entity.TaskID(id)
	return nil
}

package store

import (
	"context"
	"todotree/entity"
)

func (r *Repository) ListTasks(
	ctx context.Context, db Queryer, userID entity.UserID,
) (entity.Tasks, error) {
	tasks := entity.Tasks{}
	sql := `
		SELECT
			id, user_id, title, status,
			created, modified
		FROM task
		WHERE user_id = ?;`
	err := db.SelectContext(
		ctx, &tasks, sql, userID,
	)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *Repository) AddTask(
	ctx context.Context, db Execer, t *entity.Task,
) error {
	t.Created = r.Clocker.Now()
	t.Modified = r.Clocker.Now()
	sql := `INSERT INTO task
	(user_id, title, status, created, modified)
	VALUES (?, ?, ?, ?, ?);`
	result, err := db.ExecContext(
		ctx, sql,
		t.UserID, t.Title, t.Status, t.Created, t.Modified,
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

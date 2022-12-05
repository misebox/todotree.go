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
	(user_id, root_id, parent_id, title, status, created, modified)
	VALUES (?, ?, ?, ?, ?, ?, ?);`
	result, err := db.ExecContext(
		ctx, sql,
		t.UserID, t.RootID, t.ParentID, t.Title, t.Status, t.Created, t.Modified,
	)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	t.ID = entity.TaskID(id)
	if t.ParentID == nil {
		t.RootID = &t.ID
	} else {
		// 親タスクのルートIDを継承する
		sql := `SELECT root_id FROM task WHERE id = ?;`
		err := db.GetContext(ctx, &t.RootID, sql, t.ParentID)
		if err != nil {
			return err
		}
	}

	upd_sql := `UPDAWTE task set root_id = ? WHERE id = ?;`
	db.ExecContext(ctx, upd_sql, t.RootID, t.ID)

	return nil
}

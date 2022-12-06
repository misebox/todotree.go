package store

import (
	"context"
	"testing"

	"todotree/clock"
	"todotree/entity"
	"todotree/testutil"
	"todotree/testutil/fixture"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/go-cmp/cmp"
	"github.com/jmoiron/sqlx"
)

func prepareUser(ctx context.Context, t *testing.T, db Execer) entity.UserID {
	t.Helper()
	u := fixture.NewUserForTest()
	result, err := db.ExecContext(
		ctx,
		`INSERT INTO user (name, email, password, role, created, modified)
		VALUES (?, ?, ?, ?, ?, ?);`,
		u.Name, u.Email, u.Password, u.Role, u.Created, u.Modified,
	)
	if err != nil {
		t.Fatalf("insert user: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("got user_id: %v", err)
	}
	return entity.UserID(id)
}

func prepareTasks(ctx context.Context, t *testing.T, con Execer) (entity.UserID, entity.Tasks) {
	t.Helper()

	userID := prepareUser(ctx, t, con)
	otherUserID := prepareUser(ctx, t, con)
	c := clock.FixedClocker{}
	wants := entity.Tasks{
		{
			UserID: userID,
			Title:  "want task 1", Status: "todo",
			Created: c.Now(), Modified: c.Now(),
		},
		{
			UserID: userID,
			Title:  "want task 2", Status: "done",
			Created: c.Now(), Modified: c.Now(),
		},
	}
	tasks := entity.Tasks{
		wants[0],
		{
			UserID: otherUserID,
			Title:  "not want task", Status: "todo",
			Created: c.Now(), Modified: c.Now(),
		},
		wants[1],
	}
	result, err := con.ExecContext(ctx,
		`INSERT INTO task (title, user_id, status, created, modified)
		VALUES
		(?, ?, ?, ?, ?),
		(?, ?, ?, ?, ?),
		(?, ?, ?, ?, ?);`,
		tasks[0].Title, tasks[0].UserID, tasks[0].Status, tasks[0].Created, tasks[0].Modified,
		tasks[1].Title, tasks[1].UserID, tasks[1].Status, tasks[1].Created, tasks[1].Modified,
		tasks[2].Title, tasks[2].UserID, tasks[2].Status, tasks[2].Created, tasks[2].Modified,
	)
	if err != nil {
		t.Fatal(err)
	}
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}
	// MySQLでは発行したIDのうち最小のIDがlastInsertIDとなる
	tasks[0].ID = entity.TaskID(lastInsertID)
	tasks[1].ID = entity.TaskID(lastInsertID + 1)
	tasks[2].ID = entity.TaskID(lastInsertID + 2)
	return userID, wants
}

func TestRepository_AddTask(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	c := clock.FixedClocker{}
	// 親タスク(ルート)
	pWantID := entity.TaskID(20)
	pTask := &entity.Task{
		UserID:   999,
		RootID:   nil,
		ParentID: nil,
		Title:    "parent task",
		Status:   "todo",
		Created:  c.Now(),
		Modified: c.Now(),
	}
	// サブタスク
	cWantID := entity.TaskID(pWantID + 1)
	cTask := &entity.Task{
		UserID:   999,
		RootID:   nil,
		ParentID: &pWantID,
		Title:    "child task",
		Status:   "todo",
		Created:  c.Now(),
		Modified: c.Now(),
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	// bhavior for parent task
	mock.ExpectExec(
		// エスケープが必要
		`INSERT INTO task \(user_id, root_id, parent_id, title, status, created, modified\)
		VALUES \(\?, \?, \?, \?, \?, \?, \?\);`,
	).WithArgs(pTask.UserID, pTask.RootID, pTask.ParentID, pTask.Title, pTask.Status, c.Now(), c.Now()).
		WillReturnResult(sqlmock.NewResult(int64(pWantID), 1))
	mock.ExpectExec(
		`UPDATE task set root_id = \? WHERE id = \?;`,
	).WithArgs(pWantID, pWantID).
		WillReturnResult(sqlmock.NewResult(int64(cWantID), 1))
	// behavior for child task
	mock.ExpectExec(
		// エスケープが必要
		`INSERT INTO task \(user_id, root_id, parent_id, title, status, created, modified\)
		VALUES \(\?, \?, \?, \?, \?, \?, \?\);`,
	).WithArgs(cTask.UserID, cTask.RootID, cTask.ParentID, cTask.Title, cTask.Status, c.Now(), c.Now()).
		WillReturnResult(sqlmock.NewResult(int64(cWantID), 1))

	mock.ExpectQuery(
		`SELECT root_id FROM task WHERE id = \?;`,
	).WithArgs(pWantID).
		WillReturnRows(sqlmock.NewRows([]string{"root_id"}).AddRow(pWantID))

	mock.ExpectExec(
		`UPDATE task set root_id = \? WHERE id = \?;`,
	).WithArgs(pWantID, cWantID).
		WillReturnResult(sqlmock.NewResult(int64(cWantID), 1))

	xdb := sqlx.NewDb(db, "mysql")
	r := &Repository{Clocker: c}
	if err := r.AddTask(ctx, xdb, pTask); err != nil {
		t.Errorf("want no error, but got %v", err)
	}

	if err := r.AddTask(ctx, xdb, cTask); err != nil {
		t.Errorf("want no error, but got %v", err)
	}
	if cTask.RootID == nil || *cTask.RootID != pWantID {
		t.Fatalf("want %v, got %v", pWantID, cTask.RootID)
	}
	if cTask.ParentID == nil || *cTask.ParentID != pWantID {
		t.Fatalf("want %v, got %v", pWantID, cTask.ParentID)
	}

}

func TestRepository_ListTasks(t *testing.T) {
	ctx := context.Background()
	// entity.Taskを作成するため、トランザクションを張る
	tx, err := testutil.OpenDBForTest(t).BeginTxx(ctx, nil)
	t.Cleanup(func() { _ = tx.Rollback() })
	if err != nil {
		t.Fatal(err)
	}
	userID, wants := prepareTasks(ctx, t, tx)

	sut := &Repository{}
	gots, err := sut.ListTasks(ctx, tx, userID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d := cmp.Diff(gots, wants); len(d) != 0 {
		t.Errorf("differs: (got +want)\n%s", d)
	}
}

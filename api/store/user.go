package store

import (
	"context"
	"errors"
	"fmt"

	"todotree/entity"

	"github.com/go-sql-driver/mysql"
)

func (r *Repository) RegisterUser(ctx context.Context, db Execer, u *entity.User) error {
	u.Created = r.Clocker.Now()
	u.Modified = r.Clocker.Now()
	sql := `INSERT INTO user (
		name, email, password, role, created, modified
	) VALUES (?, ?, ?, ?, ?, ?)`
	result, err := db.ExecContext(ctx, sql, u.Name, u.Email, u.Password, u.Role, u.Created, u.Modified)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == ErrCodeMySQLDuplicateEntry {
			return fmt.Errorf("cannot create user having same name or email address: %w", ErrAlreadyEntry)
		}
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	u.ID = entity.UserID(id)
	return err
}

func (r *Repository) GetUserByName(
	ctx context.Context, db Queryer, name string,
) (*entity.User, error) {
	u := &entity.User{}
	sql := `SELECT
		id, name, email, password, role, created, modified
		FROM user WHERE name = ?`
	if err := db.GetContext(ctx, u, sql, name); err != nil {
		return nil, err
	}
	return u, nil
}

func (r *Repository) GetUserByEmail(
	ctx context.Context, db Queryer, email string,
) (*entity.User, error) {
	u := &entity.User{}
	sql := `SELECT
		id, name, email, password, role, created, modified
		FROM user WHERE email = ?`
	if err := db.GetContext(ctx, u, sql, email); err != nil {
		return nil, err
	}
	return u, nil
}

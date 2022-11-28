package service

import (
	"context"
	"fmt"
	"todotree/entity"
	"todotree/store"
)

type RegisterUser struct {
	DB   store.Execer
	Repo UserRegister
}

func (r *RegisterUser) RegisterUser(ctx context.Context, name, email, password, role string,
) (*entity.User, error) {
	u := &entity.User{
		Name:     name,
		Email:    email,
		Password: password,
		Role:     role,
	}
	err := r.Repo.RegisterUser(ctx, r.DB, u)
	if err != nil {
		return nil, fmt.Errorf("failed to register: %w", err)
	}
	return u, nil
}

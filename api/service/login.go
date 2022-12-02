package service

import (
	"context"
	"fmt"

	"todotree/store"
)

type Login struct {
	DB             store.Queryer
	Repo           UserGetter
	TokenGenerator TokenGenerator
}

func (l *Login) Login(ctx context.Context, name, pw string) (string, error) {
	user, err := l.Repo.GetUserByName(ctx, l.DB, name)
	if err != nil {
		return "", fmt.Errorf("failed to get user: %w", err)
	}

	if err := user.ComparePassword(pw); err != nil {
		return "", fmt.Errorf("wrong password: %w", err)
	}

	token, err := l.TokenGenerator.GenerateToken(ctx, *user)
	if err != nil {
		return "", fmt.Errorf("failed to generate JWT: %w", err)
	}

	return string(token), nil
}

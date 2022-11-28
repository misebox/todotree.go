package main

import (
	"context"
	"net/http"
	"todotree/clock"
	"todotree/config"
	"todotree/handler"
	"todotree/service"
	"todotree/store"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

func NewMux(ctx context.Context, cfg *config.Config) (http.Handler, func(), error) {
	mux := chi.NewRouter()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, _ = w.Write([]byte(`{"status": "ok"}`))
	})

	v := validator.New()
	db, cleanup, err := store.New(ctx, cfg)
	if err != nil {
		return nil, cleanup, err
	}
	repo := store.Repository{Clocker: clock.RealClocker{}}

	// handle /tasks
	// - add task
	add_task := &handler.AddTask{Service: &service.AddTask{DB: db, Repo: &repo}, Validator: v}
	mux.Post("/tasks", add_task.ServeHTTP)
	// - list tasks
	list_task := &handler.ListTask{Service: &service.ListTask{DB: db, Repo: &repo}}
	mux.Get("/tasks", list_task.ServeHTTP)
	ru := &handler.RegisterUser{
		Service:   &service.RegisterUser{DB: db, Repo: &repo},
		Validator: v,
	}
	mux.Post("/users", ru.ServeHTTP)
	return mux, cleanup, nil
}

package main

import (
	"context"
	"net/http"
	"todotree/auth"
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
	clocker := clock.RealClocker{}
	repo := store.Repository{Clocker: clocker}
	redis_cli, err := store.NewKVS(ctx, cfg)
	if err != nil {
		return nil, cleanup, err
	}
	jwter, err := auth.NewJWTer(redis_cli, clocker)
	if err != nil {
		return nil, cleanup, err
	}
	// handle
	ru := &handler.RegisterUser{
		Service:   &service.RegisterUser{DB: db, Repo: &repo},
		Validator: v,
	}
	mux.Post("/register", ru.ServeHTTP)

	login := &handler.Login{
		Service: &service.Login{
			DB:             db,
			Repo:           &repo,
			TokenGenerator: jwter,
		},
		Validator: v,
	}
	mux.Post("/login", login.ServeHTTP)

	add_task := &handler.AddTask{Service: &service.AddTask{DB: db, Repo: &repo}, Validator: v}
	list_task := &handler.ListTask{Service: &service.ListTask{DB: db, Repo: &repo}}

	mux.Route("/tasks", func(r chi.Router) {
		r.Use(handler.AuthMiddleware(jwter))
		r.Post("/", add_task.ServeHTTP)
		r.Get("/", list_task.ServeHTTP)
	})

	return mux, cleanup, nil
}

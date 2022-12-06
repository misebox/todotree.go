package main

import (
	"context"
	"fmt"
	"os"
	"testing"
	"todotree/config"
	"todotree/store"
)

func TestMain(m *testing.M) {
	// set up
	ctx := context.Background()
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}
	if cfg.Env != "testing" {
		panic(fmt.Sprintf("want env testing, but got %q", cfg.Env))
	}
	db, cleanup, err := store.NewDAO(ctx, cfg)
	if err != nil {
		panic(err)
	}
	tables := []string{"user", "task"}
	for _, t := range tables {
		db.ExecContext(ctx, fmt.Sprintf("TRUNCATE TABLE %s;", t))
	}
	cleanup()

	// run test
	code := m.Run()
	// tear down

	// end test
	os.Exit(code)
}

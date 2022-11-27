package main

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"todotree/config"
)

func TestNewMux(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/health", nil)
	ctx := context.Background()
	cfg, err := config.NewConfig()
	if err != nil {
		t.Fatalf("failed to create config: %v", err)
	}
	sut_mux, cleanup, err := NewMux(ctx, cfg)
	if err != nil {
		t.Errorf("failed to create mux: %v", err)
	}
	t.Cleanup(cleanup)
	sut_mux.ServeHTTP(w, r)
	resp := w.Result()
	t.Cleanup(func() { _ = resp.Body.Close() })

	if resp.StatusCode != http.StatusOK {
		t.Error("want status code 200, but", resp.StatusCode)
	}
	got, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}
	want := `{"status": "ok"}`
	if string(got) != want {
		t.Errorf("want %q, but got %q", want, got)
	}
}

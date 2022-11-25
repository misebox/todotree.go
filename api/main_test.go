package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"

	"golang.org/x/sync/errgroup"
)

func TestRun(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return run(ctx)
	})
	in := "message"
	rsp, err := http.Get(fmt.Sprintf("http://localhost:8000/%s", in))
	if err != nil {
		t.Errorf("failed to get: %+v", err)
	}
	defer rsp.Body.Close()
	got, err := io.ReadAll(rsp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}
	// HTTPサーバーの戻り値を検証
	want := fmt.Sprintf("Path: /%s", in)
	if string(got) != want {
		t.Errorf("want %q, but got %q", want, got)
	}

	// run()のコンテキストをキャンセル
	cancel()
	if err := eg.Wait(); err != nil {
		t.Fatal(err)
	}
}

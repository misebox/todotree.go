package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"testing"

	"golang.org/x/sync/errgroup"
)

func TestServer_Run(t *testing.T) {
	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("failed to listen port %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	eg, ctx := errgroup.WithContext(ctx)
	mux := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Path: %s", r.URL.Path)
	})
	eg.Go(func() error {
		s := NewServer(l, mux)
		return s.Run(ctx)
	})
	t.Logf("debug: %v", l.Addr().String())
	addr := l.Addr().String()
	in := "message"
	rsp, err := http.Get(fmt.Sprintf("http://%s/%s", addr, in))
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

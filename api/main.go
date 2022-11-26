package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"todotree/config"

	"golang.org/x/sync/errgroup"
)

func run(ctx context.Context) error {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.NewConfig()
	if err != nil {
		return err
	}
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		log.Fatalf("failed to listen port %d: %v", cfg.Port, err)
	}
	url := fmt.Sprintf("http://%s", listener.Addr().String())
	log.Printf("start with: %v", url)
	server := &http.Server{
		Addr: ":8000",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(5 * time.Second)
			fmt.Fprintf(w, "Path: %s", r.URL.Path)
		}),
	}
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		// ErrServerClosed(正常終了)以外のエラーが発生したらClose失敗
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			log.Printf("failed to close: %+v", err)
			return err
		}
		return nil
	})
	// シグナル受信時にサーバーを終了
	<-ctx.Done()
	if err := server.Shutdown(context.Background()); err != nil {
		log.Printf("failed to shutdown: %+v", err)
	}
	log.Printf("nothing to do")
	return eg.Wait()
}

func main() {
	ctx := context.Background()
	err := run(ctx)
	if err != nil {
		fmt.Printf("failed to terminate server: %v", err)
		os.Exit(1)
	}
}

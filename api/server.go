package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

type Server struct {
	server   *http.Server
	listener net.Listener
}

func NewServer(lsn net.Listener, mux http.Handler) *Server {
	return &Server{
		server: &http.Server{
			Handler: mux,
		},
		listener: lsn,
	}
}
func (s *Server) Run(ctx context.Context) error {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		// ErrServerClosed(正常終了)以外のエラーが発生したらClose失敗
		if err := s.server.Serve(s.listener); err != nil && err != http.ErrServerClosed {
			log.Printf("failed to close: %+v", err)
			return err
		}
		return nil
	})
	// シグナル受信時にサーバーを終了
	<-ctx.Done()
	if err := s.server.Shutdown(context.Background()); err != nil {
		log.Printf("failed to shutdown: %+v", err)
	}
	return eg.Wait()
}

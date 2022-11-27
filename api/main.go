package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"todotree/config"
)

func run(ctx context.Context) error {
	cfg, err := config.NewConfig()
	if err != nil {
		return err
	}
	lsn, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", cfg.Port))
	if err != nil {
		log.Fatalf("failed to listen port %d: %v", cfg.Port, err)
	}
	url := fmt.Sprintf("http://%s", lsn.Addr().String())
	log.Printf("start with: %v", url)
	mux, cleanup, err := NewMux(ctx, cfg)
	if err != nil {
		return err
	}
	defer cleanup()
	server := NewServer(lsn, mux)

	return server.Run(ctx)
}

func main() {
	ctx := context.Background()
	err := run(ctx)
	if err != nil {
		fmt.Printf("failed to terminate server: %v", err)
		os.Exit(1)
	}
}

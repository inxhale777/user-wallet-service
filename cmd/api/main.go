package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
	"user-wallet-service/config"
	v1 "user-wallet-service/internal/http/v1"
	"user-wallet-service/internal/postgres"
	"user-wallet-service/internal/repo/pg_transactions"
	"user-wallet-service/internal/service/wallet"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Unable to create config: %s", err)
	}

	p, err := postgres.New(ctx, cfg.PgURL)
	if err != nil {
		log.Fatalf("postgres initialization failed: %s", err)
	}
	defer p.Close()

	err = p.Pool.Ping(ctx)
	if err != nil {
		log.Fatalf("can't ping postgres : %s", err)
	}

	// wallet service that is not wrapped around database TX
	// used in handlers that logic does not require TX, e.g: GET /balance request
	w := wallet.New(pg_transactions.New(p.Pool), nil)

	r := v1.Run(&v1.SetupRequest{
		DB:     p,
		Wallet: w,
	})

	srv := &http.Server{
		Addr:    cfg.Addr,
		Handler: r,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-ctx.Done()

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("bye bye")
}

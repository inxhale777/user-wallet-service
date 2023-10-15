package integration_test

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"user-wallet-service/config"
	v1 "user-wallet-service/internal/http/v1"
	"user-wallet-service/internal/postgres"
	"user-wallet-service/internal/repo/pgtransactions"
	"user-wallet-service/internal/service/wallet"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/stretchr/testify/require"
)

func runMigrations(t *testing.T, cfg *config.C) {
	mg, err := migrate.New(
		"file://../migrations",
		cfg.PgDSN)
	if err != nil {
		require.Nil(t, err)
	}

	err = mg.Up()
	require.Nil(t, err)
}

func TestDeposit(t *testing.T) {
	cfg, err := config.New()
	require.Nil(t, err)

	runMigrations(t, cfg)

	ctx := context.Background()
	p, err := postgres.New(ctx, cfg.PgDSN)
	require.Nil(t, err)
	defer p.Close()

	w := wallet.New(pgtransactions.New(p), nil)
	r := v1.Run(&v1.SetupRequest{
		DB:     p,
		Wallet: w,
	})

	ts := httptest.NewServer(r)

	//TODO: WIP, need to think about test cases
	userID := 111
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			resp, err := http.Post(fmt.Sprintf("%s/wallet/deposit/%d", ts.URL, userID), //nolint
				"application/json; charset=utf-8",
				bytes.NewBufferString("{\"amount\": 100}"))
			require.Nil(t, err)
			defer func() {
				_ = resp.Body.Close()
			}()
		}()
	}

	wg.Wait()

	b, err := w.Balance(ctx, userID)
	require.Nil(t, err)
	require.Equal(t, 100*100, b)
}

package main

import (
	"context"
	"crypto/tls"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/labstack/echo/v5"
	"rsp.random/services"

	"rsp.random/config"
	"rsp.random/db"
	"rsp.random/server"
)

var Version string // Injected by ldflags at build time

func main() {

	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	// Open the Badger database located in the /tmp/badger directory.
	// It is created if it doesn't exist.
	badgerDb, err := badger.Open(badger.DefaultOptions("/tmp/badger"))
	if err != nil {
		log.Fatal(err)
	}

	defer badgerDb.Close()

	rspConfig, err := config.NewConfigWithVersion(Version)
	if err != nil {
		panic(err)
	}
	pgPool, err := db.NewDatabase(rspConfig)
	mgr := db.NewSqlManager(pgPool)

	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS12,
			},
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 10,
			IdleConnTimeout:     30 * time.Second,
		},
		Timeout: 10 * time.Second,
	}

	backgroundCh := make(chan services.UpdateCounterProcess, 1)

	go func() {
		for f := range backgroundCh {
			err := f(ctx)
			if err != nil {
				slog.Warn("background task failed", "error", err)
			}
		}
	}()

	echoServer := server.NewEchoServer(rspConfig, mgr, httpClient, badgerDb, backgroundCh)
	if err != nil {
		panic(err)
	}

	sc := echo.StartConfig{
		Address:         ":1333",
		GracefulTimeout: 5 * time.Second,
	}

	if err := sc.Start(ctx, echoServer); err != nil {
		echoServer.Logger.Error("failed to start server", "error", err)
	}

}

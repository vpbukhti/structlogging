package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"structlog/clog"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {
	err := run()
	if err != nil {
		log.Fatalf("error running my_service", err.Error())
	}
}

func run() error {
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer cancel()

	clientLogger := clog.NewLogger("my_client", slog.LevelInfo)
	serverLogger := clog.NewLogger("my_server", slog.LevelInfo)

	server := SetupServer(clog.WithLogger(ctx, serverLogger))

	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		ctx = clog.WithLogger(ctx, clientLogger)

		err := SendImportantRequests(ctx)
		if err != nil && !errors.Is(err, context.Canceled) {
			return fmt.Errorf("unable to make important requests: %w", err)
		}

		return nil
	})

	eg.Go(func() error {
		<-ctx.Done()

		cctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel()

		err := server.Shutdown(cctx)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("unable to shutdown server gracefully: %w", err)
		}

		return nil
	})

	eg.Go(func() error {
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("unable to listen and serve: %w", err)
		}

		return nil
	})

	err := eg.Wait()
	if err != nil {
		return fmt.Errorf("error during service execution: %w", err)
	}

	return nil
}

package main

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"Beanefits/internal/app"
	"Beanefits/internal/config"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(log)

	cfg := config.MustLoad()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := runMigrations(ctx, log, cfg.PostgresURL); err != nil {
		log.Error("migrations failed", "err", err)
		os.Exit(1)
	}

	a, err := app.New(ctx, cfg, log)
	if err != nil {
		log.Error("app init failed", "err", err)
		os.Exit(1)
	}
	defer func() {
		if err := a.Close(); err != nil {
			log.Error("app close failed", "err", err)
		}
	}()

	errCh := make(chan error, 1)
	go func() {
		errCh <- a.HTTP.ListenAndServe(ctx)
	}()

	select {
	case <-ctx.Done():
		log.Info("shutdown signal received", "signal", ctx.Err().Error())

		shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
		defer cancel()

		type shutdowner interface {
			Shutdown(ctx context.Context) error
		}
		if s, ok := any(a.HTTP).(shutdowner); ok {
			if err := s.Shutdown(shutdownCtx); err != nil {
				log.Error("http shutdown failed", "err", err)
			} else {
				log.Info("http server stopped gracefully")
			}
		} else {
			if err := a.HTTP.Close(); err != nil {
				log.Error("http close failed", "err", err)
			} else {
				log.Info("http server closed")
			}
		}

	case err := <-errCh:
		if err == nil {
			log.Info("http server stopped")
			return
		}
		if errors.Is(err, http.ErrServerClosed) {
			log.Info("http server closed")
			return
		}
		log.Error("http server stopped with error", "err", err)
		os.Exit(1)
	}

	time.Sleep(10 * time.Millisecond)
}

func runMigrations(ctx context.Context, log *slog.Logger, dsn string) error {
	dir := os.Getenv("GOOSE_MIGRATIONS_DIR")
	if dir == "" {
		dir = "migrations"
	}

	migCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return err
	}
	defer func() { _ = db.Close() }()

	db.SetMaxOpenConns(2)
	db.SetMaxIdleConns(2)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.PingContext(migCtx); err != nil {
		return err
	}

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	log.Info("running migrations", "dir", dir)
	if err := goose.UpContext(migCtx, db, dir); err != nil {
		return err
	}
	log.Info("migrations applied")

	return nil
}

package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/vadimbarashkov/workmate-test-task/internal/api"
	"github.com/vadimbarashkov/workmate-test-task/internal/config"

	executor "github.com/vadimbarashkov/workmate-test-task/internal/executor/memory"
	manager "github.com/vadimbarashkov/workmate-test-task/internal/manager/memory"
)

var configPath string

func main() {
	flag.StringVar(&configPath, "configPath", ".config.yml", "Path to the config file")
	flag.Parse()

	cfg, err := config.Load(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Loading config: %v\n", err)
		os.Exit(1)
	}

	setupLogger(cfg.Env)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	executor := executor.New(cfg.QueueSize, cfg.MaxWorkers)
	taskManager := manager.New(ctx, executor, cfg.TaskCleanup)

	router := api.NewRouter(slog.Default(), taskManager)

	server := &http.Server{
		Addr:           cfg.Server.Addr(),
		Handler:        router,
		ReadTimeout:    cfg.Server.ReadTimeout,
		WriteTimeout:   cfg.Server.WriteTimeout,
		IdleTimeout:    cfg.Server.IdleTimeout,
		MaxHeaderBytes: cfg.Server.MaxHeaderBytes,
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		slog.Info("starting server", slog.String("addr", cfg.Server.Addr()))
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("unexpected server error", slog.Any("err", err))
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	exitCode := 0

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := server.Shutdown(ctx); err != nil {
			slog.Error("failed to shutdown server", slog.Any("err", err))
			exitCode = 1
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := executor.Shutdown(ctx); err != nil {
			slog.Error("failed to shutdown executor", slog.Any("err", err))
			exitCode = 1
		}
	}()

	wg.Wait()
	if exitCode == 1 {
		slog.Info("server shutdown completed")
	}
	os.Exit(exitCode)
}

func setupLogger(env string) {
	level := slog.LevelDebug
	if env == config.EnvProd {
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{Level: level}

	var handler slog.Handler
	if env == config.EnvProd || env == config.EnvTest {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	logger := slog.New(handler).With(slog.String("env", env))
	slog.SetDefault(logger)
}

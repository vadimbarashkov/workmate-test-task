package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

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
	taskManager := manager.New(executor)

	router := api.NewRouter(slog.Default(), taskManager)

	server := &http.Server{
		Addr:           cfg.Server.Addr(),
		Handler:        router,
		ReadTimeout:    cfg.Server.ReadTimeout,
		WriteTimeout:   cfg.Server.WriteTimeout,
		IdleTimeout:    cfg.Server.IdleTimeout,
		MaxHeaderBytes: cfg.Server.MaxHeaderBytes,
		BaseContext: func(_ net.Listener) context.Context {
			return ctx
		},
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

	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			slog.Error("failed to shutdown server", slog.Any("err", err))
			os.Exit(1)
		}
	}()

	wg.Wait()
	slog.Info("server stopped gracefully")
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

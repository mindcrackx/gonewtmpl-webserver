package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	"log/slog"

	"github.com/mindcrackx/gonewtmpl-webserver/internal/metrics"
	"github.com/mindcrackx/gonewtmpl-webserver/internal/server"
	"github.com/mindcrackx/gonewtmpl-webserver/internal/tmpl"
	"github.com/mindcrackx/gonewtmpl-webserver/ui"
	"go.uber.org/automaxprocs/maxprocs"
)

func main() {
	if err := run(); err != nil {
		slog.Error("during run", "err", err)
		os.Exit(1)
	}
}

func setupLogger(level, handler string) {
	handlerOpts := &slog.HandlerOptions{AddSource: true, Level: slog.LevelInfo}
	switch strings.ToLower(level) {
	case "debug":
		handlerOpts.Level = slog.LevelDebug
	case "info":
		handlerOpts.Level = slog.LevelInfo
	case "warn":
		handlerOpts.Level = slog.LevelWarn
	case "error":
		handlerOpts.Level = slog.LevelError
	}

	var slogHandler slog.Handler
	switch strings.ToLower(handler) {
	case "text":
		slogHandler = slog.NewTextHandler(os.Stderr, handlerOpts)
	case "json":
		slogHandler = slog.NewJSONHandler(os.Stderr, handlerOpts)
	default:
		slogHandler = slog.NewTextHandler(os.Stderr, handlerOpts)
	}

	slog.SetDefault(slog.New(slogHandler))
}

func envOrDefault(key, def string) string {
	val := os.Getenv(key)
	if val != "" {
		return val
	}
	return def
}

func run() error {
	slog.Debug("configuring...")
	serverListenAddr := envOrDefault("MYAPP_SERVER_LISTEN_ADDR", ":8080")
	metricsListenAddr := envOrDefault("MYAPP_METRICS_LISTEN_ADDR", ":8081")
	logLevel := envOrDefault("MYAPP_LOG_LEVEL", "info")
	logHandler := envOrDefault("MYAPP_LOG_HANDLER", "text")

	setupLogger(logLevel, logHandler)

	if _, err := maxprocs.Set(); err != nil {
		return fmt.Errorf("setting maxprocs: %w", err)
	}
	slog.Info("starting up", "maxprocs", runtime.GOMAXPROCS(0))

	slog.Debug("parsing templates...")
	templates, err := tmpl.NewTemplateRenderer(ui.EmbeddedContentHTML, "html/*.html", "html/**/*.html")
	if err != nil {
		return fmt.Errorf("could not create templates: %w", err)
	}

	slog.Debug("setting up metrics server...")
	metricsServer, err := metrics.NewServer(metricsListenAddr)
	if err != nil {
		return fmt.Errorf("could not create metrics server: %w", err)
	}
	go func() {
		slog.Info("metrics server startup", "status", "server starting", "addr", metricsListenAddr)
		err := metricsServer.ListenAndServe()
		defer slog.Error("metrics server shutdown", "status", "stopped", "err", err)
	}()

	slog.Debug("setting up server...")
	appServer, err := server.New(serverListenAddr, templates, ui.EmbeddedContentStatic)
	if err != nil {
		return fmt.Errorf("could not create server: %w", err)
	}

	// graceful shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	serverErrors := make(chan error, 1)

	go func() {
		slog.Info("server startup", "status", "server starting", "addr", serverListenAddr)
		serverErrors <- appServer.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		slog.Info("server shutdown", "status", "shutdown started", "signal", sig.String())
		defer slog.Info("server shutdown", "status", "shutdown complete", "signal", sig.String())

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := appServer.Shutdown(ctx); err != nil {
			appServer.Close()
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}

	return nil
}

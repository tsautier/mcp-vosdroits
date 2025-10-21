// Package main provides the entry point for the VosDroits MCP server.
package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/guigui42/mcp-vosdroits/internal/config"
	"github.com/guigui42/mcp-vosdroits/internal/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	// Load configuration
	cfg := config.Load()

	// Set up logging
	setupLogging(cfg.LogLevel)

	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh
		slog.Info("Shutting down gracefully...")
		cancel()
	}()

	// Create MCP server
	server := mcp.NewServer(
		&mcp.Implementation{
			Name:    cfg.ServerName,
			Version: cfg.ServerVersion,
		},
		nil,
	)

	// Register tools
	if err := tools.RegisterTools(server, cfg); err != nil {
		return fmt.Errorf("failed to register tools: %w", err)
	}

	slog.Info("Starting MCP server",
		"name", cfg.ServerName,
		"version", cfg.ServerVersion,
	)

	// Use stdio transport
	transport := &mcp.StdioTransport{}
	slog.Info("Using stdio transport")

	// Run server
	if err := server.Run(ctx, transport); err != nil {
		return fmt.Errorf("server error: %w", err)
	}

	return nil
}

func setupLogging(level string) {
	var logLevel slog.Level
	switch level {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level: logLevel,
	}))
	slog.SetDefault(logger)
}

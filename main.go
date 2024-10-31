package main

import (
	"database/sql"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"mercury/internal/api"
	"mercury/internal/config"
	"mercury/internal/core"
	"mercury/internal/smtpserver"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Parse command line flags
	configFile := flag.String("config", "config/default.yaml", "Path to configuration file")
	flag.Parse()

	// Load configuration
	cfg, err := config.Load(*configFile)
	if err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize database
	db, err := sql.Open(cfg.Database.Driver, cfg.Database.URL)
	if err != nil {
		fmt.Printf("Failed to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	// Create core with new config
	core := core.NewCore(&core.Config{
		SMTPPort:    cfg.Server.SMTP.Port,
		HTTPPort:    cfg.Server.HTTP.Port,
		DatabaseURL: cfg.Database.URL,
		LogLevel:    cfg.Logging.Level,
	}, db)

	core.Logger.Info("Starting application with configuration from %s", *configFile)

	if err := core.Repository.InitializeTables(); err != nil {
		core.Logger.Fatal("Failed to initialize database tables: %v", err)
	}

	apiServer := api.NewServer(core)
	smtpSrv := smtpserver.NewServer(core)

	// Create a channel to listen for interrupt signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Create error channel for servers
	errChan := make(chan error, 2)

	// Start HTTP server with goroutine
	go func() {
		core.Logger.Info("Starting HTTP server at %s", cfg.Server.HTTP.Port)
		if err := apiServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- fmt.Errorf("HTTP server error: %v", err)
		}
	}()

	// Start SMTP server with goroutine
	go func() {
		core.Logger.Info("Starting SMTP server at %s", cfg.Server.SMTP.Port)
		if err := smtpSrv.ListenAndServe(); err != nil {
			errChan <- fmt.Errorf("SMTP server error: %v", err)
		}
	}()

	// Wait for shutdown signal or error
	select {
	case <-stop:
		core.Logger.Info("Shutting down servers...")
		// Implement graceful shutdown here
	case err := <-errChan:
		core.Logger.Fatal("Server error: %v", err)
	}
}

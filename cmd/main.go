package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"mercury/internal/api"
	"mercury/internal/config"
	"mercury/internal/core"
	"mercury/internal/imap"
	"mercury/internal/smtp"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func initDB(cfg *config.Config) (*sqlx.DB, error) {
	var db *sqlx.DB
	var err error

	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		db, err = sqlx.Connect("postgres", cfg.Database.URL)
		if err == nil {
			break
		}
		fmt.Printf("Failed to connect to database, retrying in 2 seconds... (%d/%d)\n", i+1, maxRetries)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database after %d retries: %v", maxRetries, err)
	}

	db.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	db.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime)

	return db, nil
}

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
	db, err := initDB(cfg)
	if err != nil {
		fmt.Printf("Failed to initialize database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	// Create core
	core, err := core.NewCore(&core.Config{
		SMTPPort:    cfg.Server.SMTP.Port,
		HTTPPort:    cfg.Server.HTTP.Port,
		IMAPPort:    cfg.Server.IMAP.Port,
		DatabaseURL: cfg.Database.URL,
		LogLevel:    cfg.Logging.Level,
	}, db)
	if err != nil {
		fmt.Printf("Failed to create core: %v\n", err)
		os.Exit(1)
	}

	core.Logger.Info("Starting application with configuration from %s", *configFile)

	if err := core.Repository.InitializeTables(); err != nil {
		core.Logger.Fatal("Failed to initialize database tables: %v", err)
	}

	apiServer := api.NewServer(core)
	smtpServer := smtp.NewServer(core)
	imapServer := imap.NewServer(core)

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
		if err := smtpServer.ListenAndServe(); err != nil {
			errChan <- fmt.Errorf("SMTP server error: %v", err)
		}
	}()

	// Start IMAP server with goroutine
	go func() {
		core.Logger.Info("Starting IMAP server at %s", cfg.Server.IMAP.Port)
		if err := imapServer.ListenAndServe(); err != nil {
			errChan <- fmt.Errorf("IMAP server error: %v", err)
		}
	}()

	// Wait for shutdown signal or error
	select {
	case <-stop:
		core.Logger.Info("Shutting down servers...")
		// TODO: Implement graceful shutdown here
	case err := <-errChan:
		core.Logger.Fatal("Server error: %v", err)
	}
}

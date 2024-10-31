package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"mercury/internal/api"
	"mercury/internal/core"
	"mercury/internal/smtpserver"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./email.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	config := core.LoadConfig()
	logger := log.Default()

	core := core.NewCore(config, db, logger)

	if err := core.Repository.InitializeTables(); err != nil {
		logger.Fatal("Failed to initialize database tables:", err)
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
		logger.Println("Starting HTTP server at", config.HTTPPort)
		if err := apiServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- fmt.Errorf("HTTP server error: %v", err)
		}
	}()

	// Start SMTP server with goroutine
	go func() {
		logger.Println("Starting SMTP server at", config.SMTPPort)
		if err := smtpSrv.ListenAndServe(); err != nil {
			errChan <- fmt.Errorf("SMTP server error: %v", err)
		}
	}()

	// Wait for shutdown signal or error
	select {
	case <-stop:
		logger.Println("Shutting down servers...")
		// Implement graceful shutdown here
	case err := <-errChan:
		logger.Fatal(err)
	}
}

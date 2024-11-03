package storage

import (
	"database/sql"
	"errors"
	"fmt"
)

var (
	// ErrNotFound is returned when a resource is not found in the database
	ErrNotFound = errors.New("resource not found")
	// ErrNoRowsAffected is returned when an update/delete operation affects no rows
	ErrNoRowsAffected = errors.New("no rows affected")
)

// handleDBError standardizes database error handling
func handleDBError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	return fmt.Errorf("database error: %w", err)
}

// handleRowsAffected checks if any rows were affected by an operation
func handleRowsAffected(result sql.Result) error {
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}
	if rows == 0 {
		return ErrNoRowsAffected
	}
	return nil
}

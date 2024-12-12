package storage

import (
	"database/sql"
	"errors"
	"testing"
)

func TestHandleDBError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected error
	}{
		{
			name:     "nil error returns nil",
			err:      nil,
			expected: nil,
		},
		{
			name:     "sql.ErrNoRows returns ErrNotFound",
			err:      sql.ErrNoRows,
			expected: ErrNotFound,
		},
		{
			name:     "other errors are wrapped",
			err:      errors.New("some error"),
			expected: errors.New("database error: some error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := handleDBError(tt.err)
			if tt.expected == nil {
				if got != nil {
					t.Errorf("handleDBError() = %v, want nil", got)
				}
				return
			}
			if got == nil || got.Error() != tt.expected.Error() {
				t.Errorf("handleDBError() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// mockResult implements sql.Result for testing
type mockResult struct {
	rowsAffected int64
	err          error
}

func (m mockResult) LastInsertId() (int64, error) { return 0, nil }
func (m mockResult) RowsAffected() (int64, error) { return m.rowsAffected, m.err }

func TestHandleRowsAffected(t *testing.T) {
	tests := []struct {
		name     string
		result   sql.Result
		expected error
	}{
		{
			name:     "one row affected returns nil",
			result:   mockResult{rowsAffected: 1},
			expected: nil,
		},
		{
			name:     "zero rows affected returns ErrNoRowsAffected",
			result:   mockResult{rowsAffected: 0},
			expected: ErrNoRowsAffected,
		},
		{
			name:     "error checking rows affected is wrapped",
			result:   mockResult{err: errors.New("rows error")},
			expected: errors.New("error checking rows affected: rows error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := handleRowsAffected(tt.result)
			if tt.expected == nil {
				if got != nil {
					t.Errorf("handleRowsAffected() = %v, want nil", got)
				}
				return
			}
			if got == nil || got.Error() != tt.expected.Error() {
				t.Errorf("handleRowsAffected() = %v, want %v", got, tt.expected)
			}
		})
	}
}

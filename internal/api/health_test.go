package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"inbox451/internal/config"
	"inbox451/internal/core"
	"inbox451/internal/logger"

	"github.com/labstack/echo/v4"
)

func TestHealthCheck(t *testing.T) {
	tests := []struct {
		name           string
		version        string
		commitSHA      string
		expectedStatus int
		expectedBody   HealthResponse
	}{
		{
			name:           "successful health check",
			version:        "1.0.0",
			commitSHA:      "abc123",
			expectedStatus: http.StatusOK,
			expectedBody: HealthResponse{
				Status:    "ok",
				Version:   "1.0.0",
				CommitSHA: "abc123",
			},
		},
		{
			name:           "empty version and commit",
			version:        "",
			commitSHA:      "",
			expectedStatus: http.StatusOK,
			expectedBody: HealthResponse{
				Status:    "ok",
				Version:   "",
				CommitSHA: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/health", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Create minimal core for testing
			testCore := &core.Core{
				Config:  &config.Config{},
				Logger:  logger.New(nil, logger.ERROR),
				Version: tt.version,
				Commit:  tt.commitSHA,
			}

			s := &Server{
				core: testCore,
			}

			err := s.healthCheck(c)
			if err != nil {
				t.Errorf("healthCheck() error = %v", err)
				return
			}

			if rec.Code != tt.expectedStatus {
				t.Errorf("healthCheck() status = %v, want %v", rec.Code, tt.expectedStatus)
			}

			var got HealthResponse
			if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
				t.Errorf("failed to unmarshal response: %v", err)
			}

			if got != tt.expectedBody {
				t.Errorf("healthCheck() body = %v, want %v", got, tt.expectedBody)
			}
		})
	}
}
